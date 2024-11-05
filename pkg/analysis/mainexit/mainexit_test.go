package mainexit

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestMainExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
