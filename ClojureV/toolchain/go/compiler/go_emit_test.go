package compiler

import (
	"clojurev/parser"
	"strings"
	"testing"
)

func TestGoEmitter(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai go_mock [in]
			(let [x in]
				(qurq/phi-scale x)))
	`

	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	goCode, err := EmitGo(ast, "main")
	if err != nil {
		t.Fatalf("EmitGo failed: %v", err)
	}
	if !strings.Contains(goCode, "package main") {
		t.Errorf("Go output missing package declaration")
	}
	if !strings.Contains(goCode, "func go_mock(in") {
		t.Errorf("Go output missing function definition")
	}
}

func TestCompileMainTarget(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai quantum_main [in]
			(qurq/assign out in))
	`

	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	_, err = Compile(ast, TargetGo, "main")
	if err != nil {
		t.Fatalf("Compile(TargetGo) failed: %v", err)
	}
	_, err = Compile(ast, TargetJavaScript, "")
	if err != nil {
		t.Fatalf("Compile(TargetJavaScript) failed: %v", err)
	}
	_, err = Compile(ast, TargetPython, "")
	if err != nil {
		t.Fatalf("Compile(TargetPython) failed: %v", err)
	}
	_, err = Compile(ast, TargetWasm, "")
	if err != nil {
		t.Fatalf("Compile(TargetWasm) failed: %v", err)
	}
}
