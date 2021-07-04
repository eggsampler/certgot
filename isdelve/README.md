# `certgot/isdelve`

---

This package uses a simple build tag trick and redeclaration of a variable to tell if the program is being debugged by
[delve](https://github.com/go-delve/delve).

Mostly this is just used to catch panics in release code, while letting panics fall through when debugging. This is 
easier to read a stack trace printed to console compared to whatever the logging library/package prints.

TODO: Update all build tags whenever this is out https://go.googlesource.com/proposal/+/master/design/draft-gobuild.md
Should be [1.17 onwards](https://tip.golang.org/doc/go1.17#build-lines)