[![GoDoc](https://godoc.org/github.com/teamwork/reload?status.svg)](https://godoc.org/github.com/teamwork/reload)
[![Build Status](https://travis-ci.org/Teamwork/reload.svg?branch=master)](https://travis-ci.org/Teamwork/reload)
[![codecov](https://codecov.io/gh/Teamwork/reload/branch/master/graph/badge.svg?token=n0k8YjbQOL)](https://codecov.io/gh/Teamwork/reload)

Lightweight automatic reloading of Go processes.

After initialisation with `reload.Do()` any changes to the binary (and *only*
the binary) will restart the process. For example:

```go
func main() {
    go func() {
        err := reload.Do(log.Printf)
        if err != nil {
            panic(err) // Only returns initialisation errors.
        }
    }()

    fmt.Println(os.Args)
    fmt.Println(os.Environ())
    ch := make(chan bool)
    <-ch
}
```

Now use `go install` or `go build` to restart the process.

Additional directories can be watched using `reload.Dir()`; this is useful for
reloading templates:

```go
func main() {
    go func() {
        err := reload.Do(log.Printf, reload.Dir("tpl", reloadTpl))
        if err != nil {
            panic(err)
        }
    }()
}
```

This is an alternative to the "restart binary after any `*.go` file
changed"-strategy that some other projects – such as
[gin](https://github.com/codegangsta/gin) or
[go-watcher](https://github.com/canthefason/go-watcher) – take.
The advantage of `reload`'s approach is that you have a more control over when
the process restarts, and it only watches a single directory for changes which
has some performance benefits, especially when used over NFS or Docker with a
large number of files.

It also means you won't start a whole bunch of builds if you update 20 files in
a quick succession. On a desktop this probably isn't a huge deal, but on a
laptop it'll save some battery power.

Because it's in-process you can also do things like reloading just templates
instead of recompiling/restarting everything.

Caveat: the old process will continue running happily if `go install` has a
compile error, so if you missed any compile errors due to switching the window
too soon you may get confused.

You can also use `reload.Exec()` to manually restart your process, instead of 
automatically restarting when a file changes. 
