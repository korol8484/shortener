/*
Code analysis
Usage:

	go run ./cmd/staticlint/main.go  ./...
*/
package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/korol8484/shortener/pkg/analysis/mainexit"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/quickfix/qf1001"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1001"
)

func main() {
	checks := []*analysis.Analyzer{
		// check consistency of Printf format strings and arguments
		printf.Analyzer,
		// check for possible unintended shadowing of variables
		shadow.Analyzer,
		// check that struct field tags
		structtag.Analyzer,
		// checks for unreachable code
		unreachable.Analyzer,
		// report passing non-pointer or non-interface values to unmarshal
		unmarshal.Analyzer,
		// report calls to (*testing.T).Fatal from goroutines started by a test
		testinggoroutine.Analyzer,
		// check cancel func returned by context.WithCancel is called
		lostcancel.Analyzer,
		// check for mistakes using HTTP responses
		httpresponse.Analyzer,
		errorsas.Analyzer,
		// report passing non-pointer or non-error values to errors.As
		appends.Analyzer,
		// dot imports are discouraged
		st1001.Analyzer,
		// Apply De Morgan's law
		qf1001.Analyzer,
		// detect ineffectual assignments in Go code
		ineffassign.Analyzer,
		// checks whether HTTP response body is closed successfully
		bodyclose.Analyzer,
		// Checking the use of a direct call to os.Exit in the main function of the main package
		mainexit.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	multichecker.Main(
		checks...,
	)
}
