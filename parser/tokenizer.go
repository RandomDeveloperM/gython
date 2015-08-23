package parser

import (
	"bufio"
	"os"

	"github.com/brettlangdon/gython/token"
)

var EOF rune = 0
var MAXINDENT int = 100

type TokenizerState struct {
	atBol               bool
	buffer              *bufio.Reader
	curColumn           int
	curLine             int
	curLiteral          string
	fp                  *os.File
	indentationLevel    int
	indentationAltStack []int
	indentationCurrent  int
	indentationPending  int
	indentationStack    []int
	offset              int
	tabsize             int
	tabsizeAlt          int
}

func newTokenizerState() *TokenizerState {
	return &TokenizerState{
		atBol:               true,
		curColumn:           0,
		curLine:             1,
		curLiteral:          "",
		indentationAltStack: make([]int, MAXINDENT),
		indentationCurrent:  0,
		indentationLevel:    0,
		indentationPending:  0,
		indentationStack:    make([]int, MAXINDENT),
		offset:              0,
		tabsize:             8,
	}
}

func TokenizerFromFileName(filename string) (*TokenizerState, error) {
	state := newTokenizerState()
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	state.fp = fp

	state.buffer = bufio.NewReader(state.fp)
	return state, nil
}

func (tokenizer *TokenizerState) readNext() rune {
	next, _, err := tokenizer.buffer.ReadRune()
	if err != nil {
		next = EOF
	}
	tokenizer.offset += 1
	tokenizer.curColumn += 1
	if next != EOF {
		tokenizer.curLiteral += string(next)
	}
	return next
}

func (tokenizer *TokenizerState) unread() error {
	err := tokenizer.buffer.UnreadRune()
	tokenizer.offset -= 1
	tokenizer.curColumn -= 1
	if len(tokenizer.curLiteral) > 0 {
		tokenizer.curLiteral = tokenizer.curLiteral[0 : len(tokenizer.curLiteral)-1]
	}
	return err
}

func (tokenizer *TokenizerState) finalizeToken(tok *token.Token, tokId token.TokenID) *token.Token {
	tok.ID = tokId
	tok.LineEnd = tokenizer.curLine
	tok.ColumnEnd = tokenizer.curColumn
	tok.Literal = tokenizer.curLiteral
	return tok
}

func (tokenizer *TokenizerState) newToken() *token.Token {
	tokenizer.curLiteral = ""
	return &token.Token{
		ID:          token.ERRORTOKEN,
		LineStart:   tokenizer.curLine,
		ColumnStart: tokenizer.curColumn,
		LineEnd:     tokenizer.curLine,
		ColumnEnd:   tokenizer.curColumn,
		Literal:     tokenizer.curLiteral,
	}
}

func (tokenizer *TokenizerState) parseQuoted(curTok *token.Token, nextChar rune) *token.Token {
	quote := nextChar
	quoteSize := 1
	endQuoteSize := 0
	nextChar = tokenizer.readNext()
	if nextChar == quote {
		nextChar = tokenizer.readNext()
		if nextChar == quote {
			quoteSize = 3
		} else {
			endQuoteSize = 1
		}
	}

	if nextChar != quote {
		tokenizer.unread()
	}

	for {
		if endQuoteSize == quoteSize {
			break
		}
		nextChar = tokenizer.readNext()
		if nextChar == EOF {
			return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
		}
		if quoteSize == 1 && nextChar == '\n' {
			return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
		}
		if nextChar == quote {
			endQuoteSize += 1
		} else {
			endQuoteSize = 0
			if nextChar == '\\' {
				nextChar = tokenizer.readNext()
			}
		}
	}
	return tokenizer.finalizeToken(curTok, token.STRING)
}

func (tokenizer *TokenizerState) parseNumber(curTok *token.Token, nextChar rune, fraction bool) *token.Token {
	if fraction {
		goto fraction
	}
	if nextChar == '0' {
		nextChar = tokenizer.readNext()
		if nextChar == '.' {
			tokenizer.unread()
			goto fraction
		}
		if nextChar == 'j' || nextChar == 'J' {
			tokenizer.unread()
			goto imaginary
		}
		if nextChar == 'x' || nextChar == 'X' {
			// Hex
			nextChar = tokenizer.readNext()
			if !IsXDigit(nextChar) {
				return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			}
			for IsXDigit(nextChar) {
				nextChar = tokenizer.readNext()
			}
			tokenizer.unread()
		} else if nextChar == 'o' || nextChar == 'O' {
			// Octal
			nextChar = tokenizer.readNext()
			if nextChar < '0' || nextChar >= '8' {
				return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			}
			for nextChar >= '0' && nextChar < '8' {
				nextChar = tokenizer.readNext()
			}
			tokenizer.unread()
		} else if nextChar == 'b' || nextChar == 'B' {
			// Binary
			nextChar = tokenizer.readNext()
			if nextChar != '0' && nextChar != '1' {
				return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			}
			for nextChar == '0' || nextChar == '1' {
				nextChar = tokenizer.readNext()
			}
			tokenizer.unread()
		} else {
			nonzero := false
			for nextChar == '0' {
				nextChar = tokenizer.readNext()
			}
			for IsDigit(nextChar) {
				nonzero = true
				nextChar = tokenizer.readNext()
			}
			tokenizer.unread()

			if nextChar == '.' {
				goto fraction
			} else if nextChar == 'e' || nextChar == 'E' {
				goto exponent
			} else if nextChar == 'j' || nextChar == 'J' {
				goto imaginary
			} else if nonzero {
				return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			}
			goto end
		}
	} else {
		// Decimal
		for IsDigit(nextChar) {
			nextChar = tokenizer.readNext()
		}
		tokenizer.unread()
		goto fraction
	}
fraction:
	if nextChar == '.' {
		nextChar = tokenizer.readNext()
		nextChar = tokenizer.readNext()
		for IsDigit(nextChar) {
			nextChar = tokenizer.readNext()
		}
		tokenizer.unread()
	}
exponent:
	if nextChar == 'e' || nextChar == 'E' {
		nextChar = tokenizer.readNext()
		nextChar = tokenizer.readNext()
		if nextChar == '+' || nextChar == '-' {
			nextChar = tokenizer.readNext()
			if !IsDigit(nextChar) {
				return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			}
		} else if !IsDigit(nextChar) {
			tokenizer.unread()
			tokenizer.unread()
			return tokenizer.finalizeToken(curTok, token.NUMBER)
		}
		for IsDigit(nextChar) {
			nextChar = tokenizer.readNext()
		}
		tokenizer.unread()
	}
imaginary:
	if nextChar == 'j' || nextChar == 'J' {
		nextChar = tokenizer.readNext()
		nextChar = tokenizer.readNext()
	}
end:
	return tokenizer.finalizeToken(curTok, token.NUMBER)
}

func (tokenizer *TokenizerState) Next() *token.Token {
next_line:
	curTok := tokenizer.newToken()
	nextChar := EOF
	blankline := false

	if tokenizer.atBol {
		// Get indentation level
		col := 0
		altcol := 0
		tokenizer.atBol = false
		for {
			nextChar = tokenizer.readNext()
			if nextChar == ' ' {
				col++
				altcol++
			} else if nextChar == '\t' {
				col = (col/tokenizer.tabsize + 1) * tokenizer.tabsize
				altcol = (altcol/tokenizer.tabsizeAlt + 1) * tokenizer.tabsizeAlt
			} else {
				break
			}
		}
		tokenizer.unread()

		if nextChar == '#' || nextChar == '\n' {
			// Lines with only newline or comment, shouldn't affect indentation
			if col == 0 && nextChar == '\n' {
				blankline = false
			} else {
				blankline = true
			}
		}
		if !blankline && tokenizer.indentationLevel == 0 {
			if col == tokenizer.indentationStack[tokenizer.indentationCurrent] {
				if altcol != tokenizer.indentationAltStack[tokenizer.indentationCurrent] {
					return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
				}
			} else if col > tokenizer.indentationStack[tokenizer.indentationCurrent] {
				if tokenizer.indentationCurrent+1 >= MAXINDENT {
					return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
				}
				if altcol <= tokenizer.indentationAltStack[tokenizer.indentationCurrent] {
					return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
				}
				tokenizer.indentationPending++
				tokenizer.indentationCurrent++
				tokenizer.indentationStack[tokenizer.indentationCurrent] = col
				tokenizer.indentationAltStack[tokenizer.indentationCurrent] = altcol

			} else {
				for tokenizer.indentationCurrent > 0 && col < tokenizer.indentationStack[tokenizer.indentationCurrent] {
					tokenizer.indentationPending--
					tokenizer.indentationCurrent--
				}
				if col != tokenizer.indentationStack[tokenizer.indentationCurrent] {
					return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
				}
				if altcol != tokenizer.indentationAltStack[tokenizer.indentationCurrent] {
					return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
				}
			}
		}
	}

	if tokenizer.indentationPending != 0 {
		if tokenizer.indentationPending < 0 {
			tokenizer.indentationPending++
			return tokenizer.finalizeToken(curTok, token.DEDENT)
		} else {
			tokenizer.indentationPending--
			return tokenizer.finalizeToken(curTok, token.INDENT)
		}
	}

again:
	// Skip spaces
	for {
		nextChar = tokenizer.readNext()
		if !(nextChar == ' ' || nextChar == '\t') {
			break
		}
	}
	curTok.LineStart = tokenizer.curLine
	curTok.ColumnStart = tokenizer.curColumn - 1
	tokenizer.curLiteral = string(nextChar)

	// Skip comments
	if nextChar == '#' {
		for {
			nextChar = tokenizer.readNext()
			if nextChar == EOF || nextChar == '\n' {
				break
			}
		}
	}

	// Check for EOF
	if nextChar == EOF {
		tokenizer.curLiteral = ""
		return tokenizer.finalizeToken(curTok, token.ENDMARKER)
	}

	if IsIdentifierStart(nextChar) {
		saw_b, saw_r, saw_u := false, false, false
		for {
			if !(saw_b || saw_u) && (nextChar == 'b' || nextChar == 'B') {
				saw_b = true
			} else if !(saw_b || saw_u || saw_r) && (nextChar == 'u' || nextChar == 'U') {
				saw_u = true
			} else if !(saw_r || saw_u) && (nextChar == 'r' || nextChar == 'R') {
				saw_r = true
			} else {
				break
			}
			nextChar = tokenizer.readNext()
			if IsQuote(nextChar) {
				goto letter_quote
			}
		}
		for IsIdentifierChar(nextChar) {
			nextChar = tokenizer.readNext()
		}
		tokenizer.unread()
		return tokenizer.finalizeToken(curTok, token.NAME)
	}

	// Newline
	if nextChar == '\n' {
		tokenizer.atBol = true
		if blankline || tokenizer.indentationLevel > 0 {
			goto next_line
		}
		tokenizer.curLine += 1
		tokenizer.curColumn = 0
		return tokenizer.finalizeToken(curTok, token.NEWLINE)
	}

	// Dot or number starting with dot
	if nextChar == '.' {
		nextChar = tokenizer.readNext()
		if IsDigit(nextChar) {
			return tokenizer.parseNumber(curTok, nextChar, true)
		} else if nextChar == '.' {
			nextChar = tokenizer.readNext()
			if nextChar == '.' {
				return tokenizer.finalizeToken(curTok, token.ELLIPSIS)
			} else {
				tokenizer.unread()
			}
			tokenizer.unread()
		} else {
			tokenizer.unread()
		}

		return tokenizer.finalizeToken(curTok, token.DOT)
	}

	// Number
	if IsDigit(nextChar) {
		return tokenizer.parseNumber(curTok, nextChar, false)
	}

letter_quote:
	// String
	if IsQuote(nextChar) {
		return tokenizer.parseQuoted(curTok, nextChar)
	}

	// Line continuation
	if nextChar == '\\' {
		nextChar = tokenizer.readNext()
		if nextChar != '\n' {
			return tokenizer.finalizeToken(curTok, token.ERRORTOKEN)
			goto again
		}
	}

	{
		// Check for two character tokens
		curChar := nextChar
		nextChar = tokenizer.readNext()
		tokId := GetTwoCharTokenID(curChar, nextChar)
		if tokId != token.OP {
			thirdChar := tokenizer.readNext()
			nextTokId := GetThreeCharTokenID(curChar, nextChar, thirdChar)
			if nextTokId != token.OP {
				tokId = nextTokId
			} else {
				tokenizer.unread()
			}
			return tokenizer.finalizeToken(curTok, tokId)
		}
		tokenizer.unread()
		nextChar = curChar
		tokenizer.curLiteral = string(curChar)
	}

	switch nextChar {
	case '(', '[', '{':
		tokenizer.indentationLevel++
		break
	case ')', ']', '}':
		tokenizer.indentationLevel--
		break
	}

	tokId := GetOneCharTokenID(nextChar)
	return tokenizer.finalizeToken(curTok, tokId)
}
