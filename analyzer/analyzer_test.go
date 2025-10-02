package analyzer_test

import (
	"testing"

	"github.com/sollniss/unnecessaryerror/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a/...")
}

func TestAnalyzer2(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a/func_arg_to_external/...")
}
