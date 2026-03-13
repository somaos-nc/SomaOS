package compiler

import (
	"clojurev/parser"
	"strings"
	"testing"
)

func TestExhaustiveGoOperators(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai go_kitchen_sink [in]
			(let [
				x in
				a (qurq/measure-intent-pressure x)
				b (qurq/read-topological-dimension x)
			]
				(qurq/greater-than out x)
				(qurq/less-than out x)
				(qurq/equal out x)
				(qurq/bit-not out x)
				(qurq/sum-pair out x)
				(qurq/matrix-dot out ". . . . . .")
				(qurq/matrix-dot out ". . . .")
				(qurq/quat-map out x)
				(qurq/torsional-pair out x)
				(qurq/fractal-zip x)
				(qurq/swave-interaction out x)
				(qurq/read-qudot out x)
				(qurq/read-sound-pixel out x)
				(qurq/read-thought-pixel out x)
				(qurq/transmit-qu-dot out x)
				(qurq/transmit-sound-pixel out x)
				(qurq/if (= x 1)
					(qurq/assign out x)
					(qurq/assign out 0)
				)
				(qurq/when (= x 1)
					(qurq/assign out x)
				)
			))
	`

	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	code, err := EmitGo(ast, "main")
	if err != nil {
		t.Fatalf("Failed to compile: %v", err)
	}

	if !strings.Contains(code, "func go_kitchen_sink") {
		t.Errorf("Go output missing function definition")
	}
}