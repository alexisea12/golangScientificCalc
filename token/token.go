package token

type TokenType string

type Token struct{
	Type	TokenType
	Literal	string
}

const (
	ILLEGAL		= "ILLEGAl"
	EOF			= "EOF"
	EULER		= "EULER"
	INT			= "INT"
	FLOAT		= "FLOAT"
	
	PLUS		= "+"
	MINUS		= "-"
	AST			= "*"
	DIV			= "/"
	LPAREN		= "("
	RPAREN		= ")"
	PROC		= "PROCEDURE"
	FACTORIAL	= "!"
	SINE		= "sin"
	EXP			= "EXP"
	ELEVATE		= "ELEVATE"
)

var Keywords = map[string]TokenType{
	"log":	PROC,
	"cos":	PROC,
	"sin":	PROC,
	"√":	PROC,
	"tan":	PROC,
	"ln":	PROC,
	"arcsin":	PROC,
	"arccos":	PROC,
	"arctan":	PROC,
}
