package compiler

import (
	"clojurev/parser"
	"fmt"
)

type Target string

const (
	TargetVerilog    Target = "verilog"
	TargetJavaScript Target = "javascript"
	TargetPython     Target = "python"
	TargetWasm       Target = "wasm"
	TargetGo         Target = "go"
	TargetDart       Target = "dart"
)

// Compile takes an AST and generates code for the specified target.
func Compile(ast *parser.Program, target Target, pkgName string) (string, error) {
	switch target {
	case TargetVerilog:
		return EmitVerilog(ast)
	case TargetGo:
		return EmitGo(ast, pkgName)
	case TargetJavaScript:
		return EmitJS(ast)
	case TargetPython:
		return EmitPython(ast)
	case TargetWasm:
		return EmitWasm(ast)
	default:
		return "", fmt.Errorf("target %s not yet fully implemented in new AST compiler", target)
	}
}
