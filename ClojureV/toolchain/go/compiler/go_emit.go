package compiler

import (
	"clojurev/parser"
	"fmt"
	"strings"
)

func EmitGo(ast *parser.Program, pkgName string) (string, error) {
	if pkgName == "" {
		pkgName = "main"
	}

	var sb strings.Builder

	sb.WriteString("package " + pkgName + "\n\n")
	for _, node := range ast.Body {
		if defn, ok := node.(*parser.Defn); ok {
			code, err := emitGoDefn(defn)
			if err != nil {
				return "", err
			}
			sb.WriteString(code)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

func emitGoDefn(d *parser.Defn) (string, error) {
	var g strings.Builder

	g.WriteString(fmt.Sprintf("func %s(", strings.ReplaceAll(d.Name, "-", "_")))
	for i, p := range d.Params {
		g.WriteString(fmt.Sprintf("%s float64", strings.ReplaceAll(p, "-", "_")))
		if i < len(d.Params)-1 {
			g.WriteString(", ")
		}
	}
	g.WriteString(") float64 {\n")

	for i, stmt := range d.Body {
		val := translateToGo(stmt)
		if i == len(d.Body)-1 {
			g.WriteString(fmt.Sprintf("\treturn %s\n", val))
		} else {
			g.WriteString(fmt.Sprintf("\t_ = %s\n", val))
		}
	}
	if len(d.Body) == 0 {
		g.WriteString("\treturn 0.0\n")
	}
	g.WriteString("}\n")
	return g.String(), nil
}

func translateToGo(node parser.Node) string {
	switch n := node.(type) {
	case *parser.Identifier:
		name := strings.ReplaceAll(n.Name, "-", "_")
		if name == "out" {
			return "0.0"
		}
		return name
	case *parser.Number:
		return n.Value
	case *parser.Call:
		op := strings.TrimPrefix(n.Callee, "qurq/")
		args := n.Args
		
		var argVals []string
		for _, arg := range args {
			v := translateToGo(arg)
			if v != "0.0" { // Hack to ignore out
				argVals = append(argVals, v)
			}
		}

		switch op {
		case "bit-xor":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) ^ uint64(%s))", argVals[0], argVals[1])
			}
		case "bit-and":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) & uint64(%s))", argVals[0], argVals[1])
			}
		case "bit-or":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) | uint64(%s))", argVals[0], argVals[1])
			}
		case "bit-not":
			if len(argVals) >= 1 {
				return fmt.Sprintf("float64(^uint64(%s))", argVals[0])
			}
		case "bit-shift-left":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) << %s)", argVals[0], argVals[1])
			}
		case "bit-shift-right":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) >> %s)", argVals[0], argVals[1])
			}
		case "bit-clear":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) &^ uint64(%s))", argVals[0], argVals[1])
			}
		case "bit-set":
			if len(argVals) >= 2 {
				return fmt.Sprintf("float64(uint64(%s) | uint64(%s))", argVals[0], argVals[1])
			}
		case "sum-pair":
			if len(argVals) >= 1 {
				return fmt.Sprintf("float64(uint64(int64(%s) >> 12) + uint64(int64(%s) & 0xFFF))", argVals[0], argVals[0])
			}
		}
		
		return "0.0"
	}
	return "0.0"
}
