package compiler

import (
	"clojurev/parser"
	"testing"
)

func TestExhaustiveScriptEmitHelpers(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai quantum_mock [in]
			(let [x in]
				(qurq/fractal-zip x)
				(qurq.math/sin x)))
	`
	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	for _, node := range ast.Body {
		if defn, ok := node.(*parser.Defn); ok {
			if !containsZip(defn) {
				t.Errorf("containsZip failed to detect qurq/fractal-zip")
			}
			if !containsSin(defn) {
				t.Errorf("containsSin failed to detect qurq.math/sin")
			}
		}
	}
}
