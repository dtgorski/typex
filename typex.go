// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 03/2020

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	typex "github.com/dtgorski/typex/internal"
)

type (
	options struct {
		filterParts  flagArray
		pathPatterns []string
		includeTests *bool
		includeUnexp *bool
		printVersion *bool
	}
	flagArray []string
)

func main() {
	die := func(err error) {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "typex: %s\n", err)
		os.Exit(1)
	}

	opts := getOpts()
	if *opts.printVersion {
		printf("typex: %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
		return
	}

	find := typex.TypeGrep{
		FilterFunc:   typex.CreateFilterFunc(opts.filterParts),
		IncludeUnexp: *opts.includeUnexp,
		IncludeTests: *opts.includeTests,
	}

	types, err := find.Grep(opts.pathPatterns...)
	if err != nil {
		die(err)
	}

	expo := &typex.ExGoType{
		Renderer:     &typex.TreeView{},
		IncludeUnexp: *opts.includeUnexp,
	}
	err = expo.Export(os.Stdout, types)
	if err != nil {
		die(err)
	}
}

func getOpts() options {
	opts := options{
		filterParts:  flagArray{},
		includeTests: flag.Bool("t", false, ""),
		includeUnexp: flag.Bool("u", false, ""),
		printVersion: flag.Bool("v", false, ""),
	}
	flag.Var(&opts.filterParts, "f", "")
	flag.Usage = usage
	flag.Parse()

	opts.pathPatterns = flag.Args()
	if len(opts.filterParts) == 0 {
		opts.filterParts = []string{".*"}
	}
	return opts
}

func (a *flagArray) String() string {
	return strings.Join(*a, " ")
}

func (a *flagArray) Set(value string) error {
	*a = append(*a, strings.TrimSpace(value))
	return nil
}

func printf(f string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, f, a...)
}

func usage() {
	printf(`
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
and is bound to its features and environmental execution
context.

Sources: <https://github.com/dtgorski/typex>
`)
}
