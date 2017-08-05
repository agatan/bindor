## bindor

`binder` provides binary vendoring, using Go's vendoring mechanism.

### Install

`go get github.com/agatan/bindor`

### Usage

```sh
$ # vendor github.com/golang/mock with package management system like glide or dep
$ bindor build github.com/golang/mock/mockgen
$ bindor exec which mockgen # => ./.bindor/mockgen
```

`bindor` doesn't provide package vendoring itself, it just provide the way to build executables of vendored packages and execute them.
