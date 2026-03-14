package compiler

import (
	"clojurev/parser"
	"strings"
	"testing"
)

// Exhaustive coverage for all remaining Verilog operator cases

func TestExhaustiveVerilogOperators(t *testing.T) {
	cljvInput := `
		(ns ClojureV.qurq)
		(defn-ai verilog_kitchen_sink [clk rst_n in]
			(let [
				x in
				p (qurq/measure-intent-pressure x)
				d (qurq/read-topological-dimension x)
			]
				(qurq/greater-than out x)
				(qurq/less-than out x)
				(qurq/equal out x)
				(qurq/spawn-macro-cell C1 x)
				(qurq/collapse-macro-cell C1)
				(qurq/bit-not out x)
				(qurq/sum-pair out x)
				(qurq/matrix-dot out ". . . . . .")
				(qurq/matrix-dot out ". . . .")
				(qurq/quat-map out x)
				(qurq/torsional-pair out x)
				(qurq/fractal-zip x)
				(qurq/photonic-stream out 3)
				(qurq/swave-interaction out x)
				(qurq/read-qudot out x)
				(qurq/read-sound-pixel out x)
				(qurq/read-thought-pixel out x)
				(qurq/transmit-qu-dot out x)
				(qurq/transmit-sound-pixel out x)
				(qurq/when (= x 1)
					(qurq/assign out x)
				)
				(qurq/assert-invariant out)
			))
	`

	p := parser.NewParser(cljvInput)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	code, err := EmitVerilog(ast)
	if err != nil {
		t.Fatalf("Failed to compile: %v", err)
	}

	// Just checking that compilation didn't panic and produced code.
	// The primary goal here is line execution coverage.
	if !strings.Contains(code, "module verilog_kitchen_sink") {
		t.Errorf("Verilog output missing module definition")
	}
}
