package parser

import (
	"fmt"
	"strconv"

	"github.com/hellracer2007/webCalc/calculator/ast"
	"github.com/hellracer2007/webCalc/calculator/lexer"
	"github.com/hellracer2007/webCalc/calculator/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	SUM
	MULT
	PRODUCT
	EULER
	LPAR
	PROC
)

var precedences = map[token.TokenType]int {

	token.PLUS:		SUM,
	token.MINUS:	SUM,
	token.AST:		MULT,	
	token.DIV:		MULT,
	token.EULER:	EULER,
	token.LPAREN:	LPAR,
	token.RPAREN:	LPAR,
	token.PROC:		PROC,
	token.FACTORIAL:PROC,
	token.EXP:		MULT,
	token.ELEVATE:	MULT,
}

type Parser struct {
	l	*lexer.Lexer

	curToken	token.Token
	peekToken	token.Token
	errors		[]string
	prefixParseFns	map[token.TokenType]prefixParseFn
	infixParseFns	map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

type (
	prefixParseFn	func() ast.Expression
	infixParseFn	func(ast.Expression) ast.Expression
	postfixParseFn	func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn){
	p.prefixParseFns[tokenType] = fn
}


func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:		l,
		errors:	[]string{},
	}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.EULER, p.parseEulerLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.PROC, p.parseProcedure)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.AST, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
//	p.registerInfix(token.LPAREN, p.parseGroupedExpression) 
	p.registerInfix(token.EULER, p.parseEulerExp)
	p.registerInfix(token.EXP, p.parseInfixExpression)
	p.registerInfix(token.ELEVATE, p.parseInfixExpression)
	p.registerInfix(token.PROC, p.parseInfixExpression)
	
	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.FACTORIAL, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser)  nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() * ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement()ast.Statement{
	switch p.curToken.Type {
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	p.nextToken()

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			postfix := p.postfixParseFns[p.peekToken.Type]
			if postfix == nil {
				return leftExp
			}
			p.nextToken()
			leftExp = postfix(leftExp)
			continue
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	fmt.Println("parsing infix")
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 0)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseEulerExp(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: token.Token{Type: token.AST, Literal: "*"},
		Operator: "*",
		Left: left,
		Right: ast.Euler,
	}
	return expression
}

func (p *Parser) parseEulerLiteral() ast.Expression {
	return ast.Euler
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	fmt.Println("Parsing grouped expression")
	p.nextToken()
	result := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return result
}

func (p *Parser) parseProcedure() ast.Expression {
	result := &ast.Procedure{
		Token: p.curToken,
		Func: p.curToken.Literal,
	}
	p.nextToken()
	result.Body = p.parseGroupedExpression()
	return result
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	result := &ast.PostfixExpression{
		Token: p.curToken,
		Left: left,
	}
	return result
}

func (p *Parser) peekPrecedence()int{
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}


func (p *Parser) expectPeek(t token.TokenType)bool{
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(tok token.TokenType) {
	msg := fmt.Sprintf("next token expected to be %s, go %s instead", tok, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

