package compiler

import (
	"clojurev/parser"
	"strings"
	"testing"
)

func TestVerilogCompiler(t *testing.T) {
	input := `
	(ns ClojureV.qurq)
	(defn-ai fractal_seed [clk rst_n in]
		"Manifesting Linear & Fractal Coherence"
		(qurq/quat-map mid in)
		(qurq/torsional-pair out mid)
		(qurq/phi-scale out)
		(qurq/matrix-dot ". . . . . ."))`

	p := parser.NewParser(input)
	ast, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	code, err := Compile(ast, TargetVerilog, "")
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}

	expectedParts := []string{
		"// SomaOS v3.0.0: Verilog Manifestation of fractal_seed",
		"// AI INTENT: Manifesting Linear & Fractal Coherence",
		"module fractal_seed",
		"reg [23:0] mid;",
		"mid = in;",
		"out = ~mid;",
		"out = (out * 1657) >> 10;",
		"out = in & 24'hAAA000;",
	}

	for _, part := range expectedParts {
		if !strings.Contains(code, part) {
			t.Errorf("Expected output to contain:\n%s\n\nGot:\n%s", part, code)
		}
	}
}
