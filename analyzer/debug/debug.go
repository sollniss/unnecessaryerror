package debug

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Enabled = false

// Reportf calls [analysis.Pass.Reportf] iff [Enabled] is true.
func Reportf(pass *analysis.Pass, pos token.Pos, format string, args ...any) {
	if !Enabled {
		return
	}
	pass.Reportf(pos, format, args...)
}

// PrintInstr reports the given instruction in human readable form iff [Enabled] is true.
func PrintInstr(pass *analysis.Pass, i ssa.Instruction, comment ...any) {
	if !Enabled {
		return
	}
	var comstr string
	if len(comment) > 0 {
		comstr = " // " + fmt.Sprint(comment...)
	}
	// Pos can be NoPos, but more context is better than nothing.
	if v, ok := i.(ssa.Value); ok {
		Reportf(pass, i.Pos(), "%T: %s() %s = %s%s", i, i.Parent(), v.Name(), i, comstr)
	} else {
		Reportf(pass, i.Pos(), "%T: %s() %s%s", i, i.Parent(), i, comstr)
	}
}

// PrintSSA prints the SSA data including addresses iff [Enabled] is true.
func PrintSSA(ssaInput *buildssa.SSA) {
	if !Enabled {
		return
	}
	for _, fn := range ssaInput.SrcFuncs {
		if fn.Synthetic != "" {
			continue
		}
		fmt.Printf("FUNC %p | %+#v\n", fn, fn)
		for _, block := range fn.Blocks {
			for _, instr := range block.Instrs {
				fmt.Printf("%p | %+#v\n", instr, instr)

			}
		}
	}
	fmt.Println("------")
}
