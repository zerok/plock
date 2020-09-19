# plock: flock but for POSIX locks (experimental)

This is a little experiment where I wanted to try creating a POSIX-style lock
as opposed to a FLOCK lock when using the `flock` command-line utility.

```
$ plock lockfile ./test.sh
```

For detailed information about the supported flags, run `plock --help`.
