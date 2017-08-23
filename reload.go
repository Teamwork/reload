// Package reload offers lightweight automatic reloading of running processes.
//
// After initialisation with `reload.Do()` any changes to the binary will
// restart the process.
package reload // import "github.com/teamwork/reload"

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/teamwork/log"
)

// Do reload the current process when its binary changes.
func Do() error {
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
				log.Error(err)
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write && event.Name == bin {
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

	log.Printf("restarting %#v when it changes", bin)
	<-done
	return nil
}

func exec(bin string) {
	err := syscall.Exec(bin, []string{bin}, os.Environ())
	if err != nil {
		panic(fmt.Sprintf("cannot restart: %v", err))
	}
}
