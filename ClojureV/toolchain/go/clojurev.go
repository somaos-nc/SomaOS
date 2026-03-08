package clojurev

import (
	"clojurev/compiler"
	"clojurev/parser"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Target = compiler.Target

const (
	TargetVerilog    = compiler.TargetVerilog
	TargetJavaScript = compiler.TargetJavaScript
	TargetPython     = compiler.TargetPython
	TargetWasm       = compiler.TargetWasm
	TargetGo         = compiler.TargetGo
	TargetDart       = compiler.TargetDart
)

func Transpile(cljvCode string, target Target, pkgName string) (string, error) {
	if errs := Lint(cljvCode); len(errs) > 0 {
		return "", fmt.Errorf("lint fracture: %v", errs[0])
	}
	p := parser.NewParser(cljvCode)
	ast, err := p.Parse()
	if err != nil {
		return "", err
	}
	return compiler.Compile(ast, target, pkgName)
}

func Lint(cljvCode string) []error {
	var errs []error
	if !strings.Contains(cljvCode, "(ns ClojureV.qurq)") &&
		!strings.Contains(cljvCode, "(ns soma.co-math)") &&
		!strings.Contains(cljvCode, "(ns ClojureV.gui)") &&
		!strings.Contains(cljvCode, "(ns ClojureV.quantum)") {
		errs = append(errs, fmt.Errorf("missing required namespace"))
	}
	open := strings.Count(cljvCode, "(")
	close := strings.Count(cljvCode, ")")
	if open != close {
		errs = append(errs, fmt.Errorf("unbalanced parentheses: (=%d, )=%d", open, close))
	}
	return errs
}

func CompileBinary(cljvCode string, outputPath string) error {
	goCode, err := Transpile(cljvCode, TargetGo, "main")
	if err != nil { return err }
	
	funcName := "main"
	if idx := strings.Index(goCode, "func "); idx != -1 {
		if endIdx := strings.Index(goCode[idx+5:], "("); endIdx != -1 {
			funcName = goCode[idx+5 : idx+5+endIdx]
		}
	}
	
	if !strings.Contains(goCode, "func main()") {
		goCode = strings.Replace(goCode, "package main", "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n)\n", 1)
		goCode += fmt.Sprintf(`
func main() {
	if len(os.Args) > 1 {
		v, _ := strconv.ParseFloat(os.Args[1], 64)
		fmt.Printf("%%.0f", %s(0, 0, v))
	}
}
`, funcName)
	}
	
	tmpDir, _ := os.MkdirTemp("", "cljv_build_*")
	defer os.RemoveAll(tmpDir)
	goFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(goFile, []byte(goCode), 0644)
	cmd := exec.Command("go", "build", "-o", outputPath, goFile)
	output, err := cmd.CombinedOutput()
	if err != nil { return fmt.Errorf("go build failed: %v\n%s", err, string(output)) }
	return nil
}
