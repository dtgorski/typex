[![Build Status](https://travis-ci.org/dtgorski/typex.svg?branch=master)](https://travis-ci.org/dtgorski/typex)
[![Coverage Status](https://coveralls.io/repos/github/dtgorski/typex/badge.svg?branch=master)](https://coveralls.io/github/dtgorski/typex?branch=master)

## typex

Examine [Go](https://golang.org/) types and their transitive dependencies.

### Installation
```
go get -u github.com/dtgorski/typex
```

### CLI example output
```
$ typex -f Rune io/...
├── error interface {
│       Error() string
│   }
└── io
    ├── RuneReader interface {
    │       ReadRune() (r rune, size int, err error)
    │   }
    └── RuneScanner interface {
            io.RuneReader
            ReadRune() (r rune, size int, err error)
            UnreadRune() error
        }
```

```
$ typex -u -f URL net/url
└── net
    └── url
        ├── URL struct {
        │       Scheme string
        │       Opaque string
        │       User *url.Userinfo
        │       Host string
        │       Path string
        │       RawPath string
        │       ForceQuery bool
        │       RawQuery string
        │       Fragment string
        │   }
        └── Userinfo struct {
                username string
                password string
                passwordSet bool
            }
```

```
$ typex -f Render github.com/dtgorski/typex/...
├── error interface {
│       Error() string
│   }
├── github.com
│   └── dtgorski
│       └── typex
│           └── internal
│               ├── Renderer interface {
│               │       Render(w io.Writer, m internal.ViewMap) error
│               │   }
│               └── ViewMap map[string]string
└── io
    └── Writer interface {
            Write(p []byte) (n int, err error)
        }
```

```
$ typex -h

Usage: typex [options] package...
Examine Go types and their transitive dependencies.

Options:
    -f <string>
          Type name filter expression. Repeating the -f option
          is allowed, all expressions aggregate to an OR query.

          The <string> filter can be a type name, a path part
          or a regular expression. Especially in the latter
          case, <string> should be quoted or escaped correctly
          to avoid errors during shell interpolation. Filter
          expressions are case sensitive, see examples below.

          The result set will contain additional references to
          transitive dependencies vital for the filtered types.

    -t    Go tests (files suffixed _test.go) will be included
          in the result set available for a filter expression.

    -u    Unexported types (lowercase names) will be included
          in the result set available for a filter expression.

More options:
    -h    Display this usage help and exit.
    -v    Print program version and exit.

The 'package' argument denotes one or more package import path
patterns to be inspected. Patterns must be separated by space.
A pattern containing '...' specifies the active modules whose
modules paths match the pattern.

Examples:
     $ typex -u go/...
     $ typex -u -f URL net/url
     $ typex github.com/your/repository/...

This tool relies heavily on Go's package managing subsystem
is bound to its features and environmental execution context.
```

### Disclaimer
The implementation and features of ```typex``` follow the [YAGNI](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it) principle.
There is no claim for completeness or reliability.

### @dev
Try ```make```:
```
$ make

 make help       Displays this list
 make clean      Removes build/test artifacts
 make build      Builds a static binary to ./bin/typex
 make debug      Starts debugger [:2345] with ./bin/typex
 make install    Compiles and installs typex in Go environment
 make test       Runs tests, reports coverage
 make tidy       Formats source files, cleans go.mod
 make sniff      Checks format and runs linter (void on success)
```

### License
[MIT](https://opensource.org/licenses/MIT) - © dtg [at] lengo [dot] org
