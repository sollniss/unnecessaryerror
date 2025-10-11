package unnecessaryerror

import (
	"maps"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

// reportStable iterates fnErrs in a stable order and reports the findings.
func reportStable(
	pass *analysis.Pass,
	fnErrs map[*ssa.Function]map[int][]ssa.Value,
	isUsed func(instr ssa.Instruction) bool,
) {
	sortedFns := slices.SortedFunc(maps.Keys(fnErrs), func(a, b *ssa.Function) int {
		return int(a.Pos() - b.Pos())
	})
	for _, fn := range sortedFns {
		for _, i := range slices.Sorted(maps.Keys(fnErrs[fn])) {
			errs := fnErrs[fn][i]
			var used bool
		errIdxLoop:
			for _, err := range errs {
				for _, ref := range *err.Referrers() {
					if isUsed(ref) {
						used = true
						break errIdxLoop
					}
				}
			}
			if !used {
				pass.Reportf(retValPos(fn, i), "error is only ever nil-checked; consider returning a bool instead")
			}
		}
	}
}
