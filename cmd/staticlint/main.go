package main

import (
	"honnef.co/go/tools/analysis/facts/deprecated"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/staticcheck"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"github.com/VladBag2022/goshort/internal/exitcheck"
)

func main() {
	// some analyzers from golang.org/x/tools/go/analysis/passes.
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	// all SA analyzers from staticcheck.io.
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	// one another analyzer from staticcheck.io.
	checks = append(checks, deprecated.Analyzer)

	// two another public analyzers.
	checks = append(checks,
		directives.Analyzer,
		generated.Analyzer)

	// exit analyzer.
	checks = append(checks, exitcheck.Analyzer)

	multichecker.Main(checks...)
}
