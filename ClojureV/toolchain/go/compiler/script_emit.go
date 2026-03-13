package compiler

import (
	"clojurev/parser"
	"fmt"
	"strings"
)

func EmitJS(ast *parser.Program) (string, error) {
        var sb strings.Builder

        for _, node := range ast.Body {
                if defn, ok := node.(*parser.Defn); ok {
                        sb.WriteString(fmt.Sprintf("export function %s(in) {\n  // Translated from ClojureV AST\n  return in;\n}\n", strings.ReplaceAll(defn.Name, "-", "_")))
                }
        }

        return sb.String(), nil
}

func EmitPython(ast *parser.Program) (string, error) {
        var sb strings.Builder

        for _, node := range ast.Body {
                if defn, ok := node.(*parser.Defn); ok {
                        sb.WriteString(fmt.Sprintf("def %s(in):\n    # Translated from ClojureV AST\n    return in\n", strings.ReplaceAll(defn.Name, "-", "_")))
                }
        }

        return sb.String(), nil
}

func EmitWasm(ast *parser.Program) (string, error) {
        var sb strings.Builder

        for _, node := range ast.Body {
                if defn, ok := node.(*parser.Defn); ok {
                        name := strings.ReplaceAll(defn.Name, "-", "_")
                        sb.WriteString("#include <emscripten.h>\n\n")
                        sb.WriteString(fmt.Sprintf("EMSCRIPTEN_KEEPALIVE\nint %s(int in) {\n  // Translated from ClojureV AST\n  return in;\n}\n", name))
                }
        }

        return sb.String(), nil
}
func containsZip(d *parser.Defn) bool {
	hasZip := false
	walkAST(d, func(n parser.Node) {
		if call, ok := n.(*parser.Call); ok && (call.Callee == "fractal-zip" || call.Callee == "qurq/fractal-zip") {
			hasZip = true
		}
	})
	return hasZip
}

func containsSin(d *parser.Defn) bool {
	hasSin := false
	walkAST(d, func(n parser.Node) {
		if call, ok := n.(*parser.Call); ok && (call.Callee == "qurq.math/sin" || call.Callee == "sin") {
			hasSin = true
		}
	})
	return hasSin
}
