package ast

import (
	"bytes"

	"github.com/hellracer2007/webCalc/calculator/token"
)

type Node interface {
	TokenLiteral() string
	String()	string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements	[]Statement
}

type PrefixExpression struct {
	Token		token.Token // THe prefix token, e.g. !
	Operator	string
	Right		Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}


func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token		token.Token
	Expression	Expression
}

func (es *ExpressionStatement) statementNode()	{}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}


func (il *IntegerLiteral) expressionNode()	{}
func (il *IntegerLiteral) TokenLiteral() string {return il.Token.Literal}
func (il *IntegerLiteral) String() string {return il.Token.Literal}


type InfixExpression struct {
	Token		token.Token
	Left		Expression
	Operator	string
	Right		Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {return oe.Token.Literal}
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}
func (fl *FloatLiteral) expressionNode()	{}
func (fl *FloatLiteral) TokenLiteral() string {return fl.Token.Literal}
func (fl *FloatLiteral) String() string {return fl.Token.Literal}

var Euler = &FloatLiteral{
	Token: token.Token{Type: token.INT, Literal: "2.7182818284"},
	Value: 2.7182818284,
}

type Procedure struct {
	Token token.Token
	Func	string
	Body	Expression
}

func (pr *Procedure) expressionNode()	{}
func (pr *Procedure) TokenLiteral() string {return pr.Token.Literal}
func (pr *Procedure) String() string {
	var out bytes.Buffer
	out.WriteString(pr.Token.Literal)
	out.WriteString(pr.Body.String())
	return out.String()
}

type PostfixExpression struct {
	Token	token.Token
	Operator string
	Left	Expression
}

func (pof *PostfixExpression) expressionNode()	{}
func (pof *PostfixExpression) TokenLiteral() string {return pof.Token.Literal}
func (pof *PostfixExpression) String() string {
	return pof.Token.Literal+pof.Left.String()
}
