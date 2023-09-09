package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/hellracer2007/webCalc/calculator/token"
)

type Lexer struct {
	input		 string
	position	 int
	readPosition int
	ch			 byte
}

func New(input string) *Lexer{
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar(){
	if l.readPosition >= len(l.input){
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '+' :
		tok = token.Token{Type: token.PLUS, Literal: string(l.ch)}
	case '-' :
		tok = token.Token{Type: token.MINUS, Literal: string(l.ch)}
	case '*' :
		tok = token.Token{Type: token.AST, Literal: string(l.ch)}
	case '/' :
		tok = token.Token{Type: token.DIV, Literal: string(l.ch)}
	case 'e' :
		tok = token.Token{Type: token.EULER, Literal: string(l.ch)}
	case '(' :
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case ')' :
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case '!' :
		tok = token.Token{Type: token.FACTORIAL, Literal: string(l.ch)}
	case 'E' :
		tok = token.Token{Type: token.EXP, Literal: string(l.ch)}
	case '^' :
		tok = token.Token{Type: token.ELEVATE, Literal: string(l.ch)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: string(l.ch)}
	default :
		if isDigit(l.ch) {
			lit, tp := l.readNumber()
			tok = token.Token{Type: tp, Literal: lit}
			return tok
		}else if proc := l.isProcedure(); proc != 0{
			if '√' == proc {
				tok = token.Token{Type: token.PROC, Literal: string(proc)}
				return tok
			}
		}else if isLetter(l.ch) {
			word := l.readWord()
			i, _ := token.Keywords[word]
			tok = token.Token{Type: i, Literal: word}
			fmt.Println("Parsing Token: ", tok.Literal)
			return tok
		}
	}

	l.readChar()
	return tok
} 

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func (l *Lexer) isProcedure() rune{
	drune, le := utf8.DecodeRune([]byte(l.input[l.position:]))
	if drune == '√' {
		l.position += le
		l.readPosition += le
		l.ch = l.input[l.position]
		return drune 
	}
	return 0
}

func (l *Lexer) readNumber() (string, token.TokenType){
	float := false
	position := l.position
	for isDigit(l.ch) || (l.ch == '.' && !float) {
		if l.ch == '.' {
			float = true
		}
		l.readChar()
	}
	var tokenType token.TokenType
	if float == false {
		tokenType = token.INT
	} else {
		tokenType = token.FLOAT
	}
	return l.input[position:l.position], tokenType
}

func (l *Lexer) readWord() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
