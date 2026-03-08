package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `(ns ClojureV.qurq) 
	; comment
	(defn test_seed [clk rst_n in] (qurq/bit-xor out in 0xABCDEF))`
	
	lex := NewLexer(input)
	
	expected := []struct{
		tokType TokenType
		val     string
	}{
		{TokenLParen, "("}, {TokenIdent, "ns"}, {TokenIdent, "ClojureV.qurq"}, {TokenRParen, ")"},
		{TokenLParen, "("}, {TokenIdent, "defn"}, {TokenIdent, "test_seed"},
		{TokenLBrack, "["}, {TokenIdent, "clk"}, {TokenIdent, "rst_n"}, {TokenIdent, "in"}, {TokenRBrack, "]"},
		{TokenLParen, "("}, {TokenIdent, "qurq/bit-xor"}, {TokenIdent, "out"}, {TokenIdent, "in"}, {TokenNum, "0xABCDEF"}, {TokenRParen, ")"},
		{TokenRParen, ")"},
	}

	for i, exp := range expected {
		tok := lex.NextToken()
		if tok.Type != exp.tokType || tok.Value != exp.val {
			t.Errorf("Token %d: expected %s '%s', got %s '%s'", i, exp.tokType, exp.val, tok.Type, tok.Value)
		}
	}
}

func TestParser(t *testing.T) {
	input := `
	(ns ClojureV.qurq)
	(defn-ai xor_seed [clk rst_n in] 
		"Manifesting intent"
		(qurq/bit-xor out in 0xABCDEF))`

	p := NewParser(input)
	ast, err := p.Parse()
	
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	if len(ast.Body) != 2 {
		t.Fatalf("Expected 2 body nodes, got %d", len(ast.Body))
	}

	ns, ok := ast.Body[0].(*Namespace)
	if !ok || ns.Name != "ClojureV.qurq" {
		t.Errorf("Expected Namespace ClojureV.qurq")
	}

	defn, ok := ast.Body[1].(*Defn)
	if !ok || defn.Name != "xor_seed" {
		t.Fatalf("Expected Defn xor_seed")
	}

	if !defn.IsAI || defn.Intent != "Manifesting intent" {
		t.Errorf("Expected defn-ai with intent, got AI:%v Intent:%s", defn.IsAI, defn.Intent)
	}

	if len(defn.Body) != 1 {
		t.Fatalf("Expected 1 expression in defn body, got %d", len(defn.Body))
	}

	call, ok := defn.Body[0].(*Call)
	if !ok || call.Callee != "qurq/bit-xor" {
		t.Errorf("Expected call to qurq/bit-xor")
	}
}
