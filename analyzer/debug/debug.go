package debug

import (
	"fmt"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const Enabled = true

// Reportf calls [analysis.Pass.Reportf] iff [Enabled] is true.
func Reportf(pass *analysis.Pass, pos token.Pos, format string, args ...any) {
	if !Enabled {
		return
	}
	pass.Reportf(pos, format, args...)
}

// PrintInstr reports the given instruction in human readable form iff [Enabled] is true.
func PrintInstr(pass *analysis.Pass, i ssa.Instruction) {
	if !Enabled {
		return
	}
	// Pos can be NoPos, but more context is better than nothing.
	if v, ok := i.(ssa.Value); ok {
		Reportf(pass, i.Pos(), "%s: %s = %s", reflect.TypeOf(i), v.Name(), i)
	} else {
		Reportf(pass, i.Pos(), "%s: %s", reflect.TypeOf(i), i)
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
