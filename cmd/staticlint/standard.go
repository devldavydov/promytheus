package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
)

// Standard go analyzers.
var standardAnalyzers = []*analysis.Analyzer{
	// report mismatches between assembly files and Go declarations
	asmdecl.Analyzer,
	// check for useless assignments
	assign.Analyzer,
	// check for common mistakes using the sync/atomic package
	atomic.Analyzer,
	// check for non-64-bits-aligned arguments to sync/atomic functions
	atomicalign.Analyzer,
	// check for common mistakes involving boolean operators
	bools.Analyzer,
	// defines an Analyzer that checks build tags
	buildtag.Analyzer,
	// defines an Analyzer that detects some violations of the cgo pointer passing rules
	cgocall.Analyzer,
	// defines an Analyzer that checks for unkeyed composite literals
	composite.Analyzer,
	// defines an Analyzer that checks for locks erroneously passed by value
	copylock.Analyzer,
	// defines an Analyzer that checks for the use of reflect.DeepEqual with error values
	deepequalerrors.Analyzer,
	// defines an Analyzer that checks known Go toolchain directives
	directive.Analyzer,
	// defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error
	errorsas.Analyzer,
	// defines an Analyzer that detects structs that would use less memory if their fields were sorted
	fieldalignment.Analyzer,
	// defines an Analyzer that reports assembly code that clobbers the frame pointer before saving it
	framepointer.Analyzer,
	// defines an Analyzer that checks for mistakes using HTTP responses
	httpresponse.Analyzer,
	// defines an Analyzer that flags impossible interface-interface type assertions
	ifaceassert.Analyzer,
	// defines an Analyzer that checks for references to enclosing loop variables from within nested functions
	loopclosure.Analyzer,
	// defines an Analyzer that checks for failure to call a context cancellation function
	lostcancel.Analyzer,
	// defines an Analyzer that checks for useless comparisons against nil
	nilfunc.Analyzer,
	// check for redundant or impossible nil comparisons
	nilness.Analyzer,
	// printf defines an Analyzer that checks consistency of Printf format strings and arguments
	printf.Analyzer,
	// check for comparing reflect.Value values with == or reflect.DeepEqual
	reflectvaluecompare.Analyzer,
	// defines an Analyzer that checks for shadowed variables
	shadow.Analyzer,
	// defines an Analyzer that checks for shifts that exceed the width of an integer
	shift.Analyzer,
	// defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify
	sigchanyzer.Analyzer,
	// defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument
	sortslice.Analyzer,
	// defines an Analyzer that checks for misspellings in the signatures of methods similar to well-known interfaces
	stdmethods.Analyzer,
	// defines an Analyzer that flags type conversions from integers to strings
	stringintconv.Analyzer,
	// defines an Analyzer that checks struct field tags are well formed
	structtag.Analyzer,
	// defines an Analyzerfor detecting calls to Fatal from a test goroutine
	testinggoroutine.Analyzer,
	// defines an Analyzer that checks for common mistaken usages of tests and examples
	tests.Analyzer,
	// defines an Analyzer that checks for the use of time.Format or time.Parse calls with a bad format
	timeformat.Analyzer,
	// defines an Analyzer that checks for passing non-pointer or non-interface types to unmarshal and decode
	unmarshal.Analyzer,
	// defines an Analyzer that checks for unreachable code
	unreachable.Analyzer,
	// defines an Analyzer that checks for invalid conversions of uintptr to unsafe.Pointer
	unsafeptr.Analyzer,
	// defines an analyzer that checks for unused results of calls to certain pure functions
	unusedresult.Analyzer,
	// checks for unused writes to the elements of a struct or array object
	unusedwrite.Analyzer,
	// defines an Analyzer that checks for usage of generic features added in Go 1.18
	usesgenerics.Analyzer,
}
