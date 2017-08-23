The reload package offers lightweight automatic reloading of running processes.
After initialisation with `reload.Do()` any changes to the binary will restart
the process.

This works well with the standard `go install` and `go build` commands.

This is an alternative to the "restart binary after any `*.go` file
changed"-strategy that some other projects take (such as
[go-watcher](https://github.com/canthefason/go-watcher)). The advantage of
`reload`'s approach is that you have a bit more control over when the process
restarts, and it only watches a single directory for changes, which has some
performance benefits when used over NFS or Docker.

Caveat: the old process will continue running happily if `go install` has a
compile error, so if you missed any compile errors due to switching the window
too soon you may get confused.