[![Build Status](https://travis-ci.org/dtgorski/typex.svg?branch=master)](https://travis-ci.org/dtgorski/typex)
[![Coverage Status](https://coveralls.io/repos/github/dtgorski/typex/badge.svg?branch=master)](https://coveralls.io/github/dtgorski/typex?branch=master)
[![Open Issues](https://img.shields.io/github/issues/dtgorski/typex.svg)](https://github.com/dtgorski/typex/issues)
[![Report Card](https://goreportcard.com/badge/github.com/dtgorski/typex)](https://goreportcard.com/report/github.com/dtgorski/typex)
[![Awesome Go](https://awesome.re/badge.svg)](https://github.com/avelino/awesome-go#user-content-go-tools)

## typex

Examine [Go](https://golang.org/) types and their transitive dependencies. Export results as TypeScript value objects (or types) declaration.

### Installation
```
go get -u github.com/dtgorski/typex
```

### Synopsis
The CLI command ```typex``` filters and displays [Go](https://golang.org/) type structures, interfaces and their relationships across package boundaries.
It generates a type hierarchy tree with additional references to transitive dependencies vital for the filtered types.
As an additional feature, ```typex``` exports the result tree as a [TypeScript](https://www.typescriptlang.org/) projection representing value objects or bare types.

### Examples
**Go type hierarchy layout**
  ```
  $ typex -f=Rune io/...

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
  $ typex -f=Render github.com/dtgorski/typex/...

  └── github.com
      └── dtgorski
          └── typex
              └── internal
                  ├── PathReplaceFunc func(string) string
                  ├── go
                  │   └── TypeRender struct {
                  │           PathReplaceFunc internal.PathReplaceFunc
                  │           IncludeUnexported bool
                  │       }
                  └── ts
                      └── TypeRender struct {
                              PathReplaceFunc internal.PathReplaceFunc
                              IncludeUnexported bool
                          }

  ```

**TypeScript value object layout**
  ```
  $ typex -f=File -l=ts-class mime/multipart

  export module mime {
      export module multipart {
          export class FileHeader {
              constructor(
                  readonly Filename: string,
                  readonly Header: net.textproto.MIMEHeader,
                  readonly Size: number,
              ) {}
          }
      }
  }
  export module net {
      export module textproto {
          export type MIMEHeader = Record<string, string[]>
      }
  }
  ```

**TypeScript bare type layout**
  ```
  $ typex -f=File -l=ts-type mime/multipart

  export module mime {
      export module multipart {
          export type FileHeader = {
              Filename: string,
              Header: net.textproto.MIMEHeader,
              Size: number,
          }
      }
  }
  export module net {
      export module textproto {
          export type MIMEHeader = Record<string, string[]>
      }
  }
  ```

### TypeScript and reserved keywords
Basically, the names of types and fields will be exported from Go without modification.
Collisions with reserved keywords or standard type names in the target language may occur.
To avoid conflicts, you may use the JSON tag annotation for the exported fields of a struct as described in the [json.Marshal(...)](https://golang.org/pkg/encoding/json/#Marshal) documentation.

### TypeScript and exportable types
Due to fundamental language differences, ```typex``` is not capable of exporting all type declarations one-to-one. Refer to the type mapping table below. 
Go channel, interface and function declarations will be omitted, references to these declarations will be typed with ```any```.

### TypeScript type mapping
TypeScript (resp. JavaScript aka ECMAScript) lacks a native integer number type.
The numeric type provided there is inherently a 64 bit float.
You should keep this in mind when working with exported numeric types - this includes `byte` and `rune` type aliases as well.    

|Go native type|TypeScript type
| --- | ---
|```bool```|```boolean```
|```string```|```string```
|```map```|```Record<K, V>```
|```struct``` ```(named)```|```T```
|```struct``` ```(anonymous)```|```{}```
|```array``` ```(slice)```|```T[]```
|```complex```[```64```&vert;```128```]|```any```
|```chan```, ```func```, ```interface```|```any```
|```int```[```8```&vert;```16```&vert;```32```&vert;```64```]|```number```
|```uint```[```8```&vert;```16```&vert;```32```&vert;```64```]|```number```
|```byte```(=```uint8```)|```number```
|```rune```(=```int32```)|```number```
|```float```[```32```&vert;```64```]|```number```
|```uintptr```|```number```

### Usage

```
$ typex -h
```
```
Usage: typex [options] package...
Examine Go types and their transitive dependencies. Export
results as TypeScript value objects (or types) declaration.

Options:
    -f <name>
        Type name filter expression. Repeating the -f option
        is allowed, all expressions aggregate to an OR query.

        The <name> filter can be a type name, a path part or
        a regular expression. Especially in the latter case,
        <name> should be quoted or escaped correctly to avoid
        errors during shell interpolation. Filters are case
        sensitive, see examples below.

        The result tree will contain additional references to
        transitive dependencies vital for the filtered types.

    -l <layout>
        Modify the export layout. Available layouts are:
          * "go":       the default Go type dependency tree
          * "ts-type":  TypeScript type declaration projection
          * "ts-class": TypeScript value object projection

    -r <old-path>:<new-path>
        Replace matching portions of <old-path> in a fully
        qualified type name with <new-path> string. Repeating
        the -r option is allowed, substitutions will perform
        successively. <old-path> can be a regular expression.
        
        The path replacement/relocation can be used to modify
        locations of type hierarchies, e.g. prune off the
        "github.com" reference from qualified type name path
        by omitting the <new-path> part after the colon. 

    -t  Go tests (files suffixed _test.go) will be included
        in the result tree available for a filter expression

    -u  Unexported types (lowercase names) will be included
        in the result tree available for a filter expression.

    -x <name> 
        Exclude type names from export. Repeating this option
        is allowed, all expressions aggregate to an OR query.
        The exclusion filter can be a type name, a path part
        or a regular expression.

More options:
    -h  Display this usage help and exit.
    -v  Print program version and exit.

The 'package' argument denotes one or more package import path
patterns to be inspected. Patterns must be separated by space.
A pattern containing '...' specifies the active modules whose
modules paths match the pattern.

Examples:
    $ typex -u go/...
    $ typex -u -f=URL net/url
    $ typex github.com/your/repository/...
    $ typex -l=ts-type github.com/your/repository/...
    $ typex -r=github.com:a/b/c github.com/your/repository/...

This tool relies heavily on Go's package managing subsystem and
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
