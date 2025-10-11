package unnecessaryerror_test

import (
	"testing"

	"github.com/sollniss/unnecessaryerror/unnecessaryerror"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, unnecessaryerror.Analyzer, "a/...")
}
