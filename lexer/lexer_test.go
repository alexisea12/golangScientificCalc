package lexer

import (
	"testing"

	"github.com/hellracer2007/webCalc/calculator/token"
)

func TestNextToken(t *testing.T) {
	imput := `234.333`
	type test struct{
		expectedType	token.TokenType
		expectedLiteral	string
	}
	var floatTest = test{expectedType: token.FLOAT, expectedLiteral: "234.333"}

	lex := New(imput)
	tok := lex.NextToken()
	if tok.Type != floatTest.expectedType {
		t.Fatalf("expected type %s got %s", floatTest.expectedType, tok.Type)
	}
	if tok.Literal != floatTest.expectedLiteral {
		t.Fatalf("expected literal %s got %s", floatTest.expectedLiteral, tok.Literal)
	}
}
