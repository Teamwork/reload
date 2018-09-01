// Package reload offers lightweight automatic reloading of running processes.
//
// After initialisation with reload.Do() any changes to the binary will
// restart the process.
//
// Example:
//
//    go func() {
//        err := Do(log.Printf)
//        if err != nil {
//            panic(err)
//        }
//    }()
package reload // import "github.com/teamwork/reload"

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

type logger interface {
	Printf(string, ...interface{})
}

// Do reload the current process when its binary changes.
//
// The log function is used to display an informational startup message and
// errors. It works well with e.g. the standard log package or Logrus.
//
// The error return will only return initialisation/startup errors. Once
// initialized it will use the log function to print errors, rather than return.
func Do(log func(string, ...interface{})) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "cannot setup watcher")
	}
	defer watcher.Close() // nolint: errcheck

	bin, err := filepath.Abs(os.Args[0])
	if err != nil {
		return errors.Wrapf(err, "cannot get Abs of %#v", os.Args[0])
	}

	dir := filepath.Dir(bin)

	done := make(chan bool)
	go func() {
		for {
			select {
			case err := <-watcher.Errors:
				// Standard logger doesn't have anything other than Print,
				// Panic, and Fatal :-/ Printf() is probably best.
				log("reload error: %v", err)
			case event := <-watcher.Events:
				// Ensure that we use the correct events, as they are not uniform accross
				// platforms. See https://github.com/fsnotify/fsnotify/issues/74
				var trigger bool
				switch runtime.GOOS {
				case "darwin", "freebsd", "openbsd", "netbsd", "dragonfly":
					trigger = event.Op&fsnotify.Create == fsnotify.Create
				case "linux":
					trigger = event.Op&fsnotify.Write == fsnotify.Write
				default:
					trigger = event.Op&fsnotify.Create == fsnotify.Create
					log("reload: untested GOOS %q; this package may not work correctly", runtime.GOOS)
				}

				if trigger && event.Name == bin {
					// Wait for writes to finish.
					time.Sleep(100 * time.Millisecond)
					exec(bin)
				}

			}
		}
	}()

	// Watch the directory, because a recompile renames the existing file
	// (rather than rewriting it), so we won't get events for that.
	if err := watcher.Add(dir); err != nil {
		return errors.Wrapf(err, "cannot add %#v to watcher", dir)
	}

	log("restarting %#v when it changes", bin)
	<-done
	return nil
}

func exec(bin string) {
	err := syscall.Exec(bin, []string{bin}, os.Environ())
	if err != nil {
		panic(fmt.Sprintf("cannot restart: %v", err))
	}
}
