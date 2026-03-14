package compiler

import (
	"clojurev/parser"
	"strings"
	"testing"
)

func TestScriptEmitters(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai quantum_mock [in]
			(let [x in]
				(fractal-zip x)
				(sin x)
				(qurq/phi-scale x)))
	`

	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Test JS
	js, err := EmitJS(ast)
	if err != nil {
		t.Fatalf("EmitJS failed: %v", err)
	}
	if !strings.Contains(js, "export function quantum_mock(in)") {
		t.Errorf("JS output missing function definition")
	}

	// Test Python
	py, err := EmitPython(ast)
	if err != nil {
		t.Fatalf("EmitPython failed: %v", err)
	}
	if !strings.Contains(py, "def quantum_mock(in):") {
		t.Errorf("Python output missing function definition")
	}

	// Test Wasm (C base)
	wasm, err := EmitWasm(ast)
	if err != nil {
		t.Fatalf("EmitWasm failed: %v", err)
	}
	if !strings.Contains(wasm, "EMSCRIPTEN_KEEPALIVE") {
		t.Errorf("WASM output missing Emscripten keepalive macro")
	}
}
