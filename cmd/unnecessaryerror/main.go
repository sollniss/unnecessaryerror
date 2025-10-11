package main

import (
	"github.com/sollniss/unnecessaryerror/unnecessaryerror"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(unnecessaryerror.Analyzer)
}
