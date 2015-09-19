package token

type TokenID int

const (
	ENDMARKER TokenID = iota
	NAME
	NUMBER
	STRING
	NEWLINE
	INDENT
	DEDENT
	LPAR
	RPAR
	LSQB
	RSQB
	COLON
	COMMA
	SEMI
	PLUS
	MINUS
	STAR
	SLASH
	VBAR
	AMPER
	LESS
	GREATER
	EQUAL
	DOT
	PERCENT
	LBRACE
	RBRACE
	EQEQUAL
	NOTEQUAL
	LESSEQUAL
	GREATEREQUAL
	TILDE
	CIRCUMFLEX
	LEFTSHIFT
	RIGHTSHIFT
	DOUBLESTAR
	PLUSEQUAL
	MINEQUAL
	STAREQUAL
	SLASHEQUAL
	PERCENTEQUAL
	AMPEREQUAL
	VBAREQUAL
	CIRCUMFLEXEQUAL
	LEFTSHIFTEQUAL
	RIGHTSHIFTEQUAL
	DOUBLESTAREQUAL
	DOUBLESLASH
	DOUBLESLASHEQUAL
	AT
	ATEQUAL
	RARROW
	ELLIPSIS
	OP
	AWAIT
	ASYNC
	ERRORTOKEN
	N_TOKENS
)

func (id TokenID) String() string {
	return TokenNames[id]
}
