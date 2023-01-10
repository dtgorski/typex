// MIT license · Daniel T. Gorski · dtg [at] lengo [dot] org · 06/2020

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	typex "github.com/dtgorski/typex/internal"
	g0 "github.com/dtgorski/typex/internal/go"
	"github.com/dtgorski/typex/internal/ts"
)

type (
	options struct {
		includeParts flagArray
		excludeParts flagArray
		replaceParts flagArray
		pathPatterns []string
		outputLayout *string
		includeTests *bool
		includeUnexp *bool
		printVersion *bool
	}
	flagArray []string
)

func main() {
	write := func(msg string) {
		out := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(out, "typex: %s\n", msg)
	}

	opts, err := getOpts()
	if err != nil {
		write(err.Error())
		os.Exit(1)
	}
	if *opts.printVersion {
		write(Version + " " + runtime.GOOS + " " + runtime.GOARCH)
		os.Exit(0)
	}

	pac := typex.Packagist{
		PathFilterFunc:    typex.CreatePathFilterFunc(opts.includeParts, opts.excludeParts),
		IncludeUnexported: *opts.includeUnexp,
		IncludeTestFiles:  *opts.includeTests,
	}
	types, err := pac.Inspect(opts.pathPatterns...)
	if err != nil {
		write(err.Error())
		os.Exit(1)
	}

	switch *opts.outputLayout {
	case "go":
		err = exportGo(opts, types)
	case "ts-type":
		err = exportTs(opts, types, false)
	case "ts-class":
		err = exportTs(opts, types, true)
	}
	if err != nil {
		write(err.Error())
		os.Exit(1)
	}
}

func exportGo(opts options, types typex.TypeMap) error {
	tr := g0.TypeRender{
		PathReplaceFunc:   typex.CreatePathReplaceFunc(opts.replaceParts),
		IncludeUnexported: *opts.includeUnexp,
	}
	tw := typex.TreeWalk{
		Layout: g0.NewTreeLayout(os.Stdout),
	}
	return tw.Walk(tr.Render(types))
}

func exportTs(opts options, types typex.TypeMap, exportObjs bool) error {
	tr := ts.TypeRender{
		PathReplaceFunc:   typex.CreatePathReplaceFunc(opts.replaceParts),
		IncludeUnexported: *opts.includeUnexp,
	}
	tw := typex.TreeWalk{
		Layout: ts.NewModuleLayout(os.Stdout),
	}
	return tw.Walk(tr.Render(types, exportObjs))
}

func getOpts() (options, error) {
	opts := options{
		includeParts: flagArray{},
		excludeParts: flagArray{},
		replaceParts: flagArray{},
		outputLayout: flag.String("l", "", ""),
		includeTests: flag.Bool("t", false, ""),
		includeUnexp: flag.Bool("u", false, ""),
		printVersion: flag.Bool("v", false, ""),
	}
	flag.Var(&opts.includeParts, "f", "")
	flag.Var(&opts.excludeParts, "x", "")
	flag.Var(&opts.replaceParts, "r", "")
	flag.Usage = usage
	flag.Parse()

	switch *opts.outputLayout {
	case "go", "ts-type", "ts-class":
	default:
		*opts.outputLayout = "go"
	}

	opts.pathPatterns = flag.Args()
	if len(opts.includeParts) == 0 {
		opts.includeParts = []string{".*"}
	}
	return opts, nil
}

func (a *flagArray) String() string {
	return strings.Join(*a, " ")
}

func (a *flagArray) Set(value string) error {
	*a = append(*a, strings.TrimSpace(value))
	return nil
}

func write(w io.Writer, f string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, f, a...)
}

func usage() {
	write(flag.CommandLine.Output(), `
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

Sources: <https://github.com/dtgorski/typex>
`)
}
