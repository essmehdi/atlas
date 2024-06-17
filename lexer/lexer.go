package lexer

import (
	"atlas/utils"
	"fmt"
	"unicode"
)

var KEYWORDS = []string{"if", "else", "return", "var", "int", "uint", "bool", "while", "for", "function", "true", "false"}

var OPERATORS_FIRSTS = []byte{'+', '-', '*', '/', '<', '>', '&', '|', '!', '=', '~'}
var OPERATORS_ASSIGN_MAP = map[string]TokenType{
	"=":  ASSIGN,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	"<":  LT,
	">=": GEQ,
	"<=": LEQ,
	"+":  PLUS,
	"-":  MINUS,
	"*":  MULTIPLY,
	"/":  DIVIDE,
	"!":  BANG,
	"&":  BIT_AND,
	"|":  BIT_OR,
	"~":  BIT_NOT,
	"&&": LOGICAL_AND,
	"||": LOGICAL_OR,
}

var KEYWORDS_MAP = map[string]TokenType{
	"var":      VAR,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"while":    WHILE,
	"loop":     LOOP,
	"int":      TYPE_INT,
	"uint":     TYPE_UINT,
	"bool":     TYPE_BOOL,
	"fun": 		FUN,
	"true":     TRUE,
	"false":    FALSE,
}

var TYPES_KEYWORDS = []TokenType{TYPE_INT, TYPE_UINT, TYPE_BOOL}

type TokenType int

const (
	VAR TokenType = iota // Keywords
	IF
	ELSE
	RETURN
	WHILE
	LOOP
	FUN

	TRUE // Built-in literals
	FALSE

	TYPE_INT // Built-in types
	TYPE_UINT
	TYPE_BOOL

	IDENTIFIER  // An identifier variable
	LITERAL_INT // A LITERAL number
	OPERATOR    // An operator

	EQ  // ==
	NEQ // ==
	GT  // >
	LT  // <
	GEQ // >=
	LEQ // <=

	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /

	BANG // !

	BIT_AND // &
	BIT_OR  // |
	BIT_NOT // ~

	LOGICAL_AND // &&
	LOGICAL_OR  // ||

	ASSIGN    // Assignment =
	LPAR      // Left parenthesis (
	RPAR      // Right parenthesis )
	LBRACE    // Left curly brace {
	RBRACE    // Right curly brace }
	LBRACKET  // Left bracket [
	RBRACKET  // Right bracket ]
	SEMICOLON // Semicolon ;
	COLON     // Colon :
	COMMA	  // Comma ,
	ILLEGAL   // Illegal token
	EOF       // End of file
)

func (d TokenType) String() string {
	return [...]string{
		"var keyword",
		"if keyword",
		"else keyword",
		"return keyword",
		"while keyword",
		"for keyword",
		"function keyword",

		"true keyword",
		"false keyword",

		"Integer type keyword",
		"Unsigned int type keyword",
		"Boolean type keyword",

		"Identifier",
		"Literal number",
		"Operator",

		"Equal",
		"Not Equal",
		"Greater",
		"Less",
		"Greater or equal",
		"Less or equal",

		"Plus",
		"Minus",
		"Multiply",
		"Divide",

		"Bang",

		"Bit AND",
		"Bit OR",
		"Bit NOT",

		"Logical AND",
		"Logical OR",

		"Assign",
		"Left Parenthesis",
		"Right Parenthesis",
		"Left brace",
		"Right brace",
		"Left bracket",
		"Right bracket",
		"Semicolon",
		"Colon",
		"Comma",
		"Illegal",
		"End of file",
	}[d]
}

type Token struct {
	Type  TokenType // The Token type
	Value string    // The lexem/value of this token
	Row   int       // The row in which this token appears
	Col   int       // The column in which this token appears
}

func (token *Token) FormattedLocation() string {
	return fmt.Sprintf("at line %d, column %d", token.Row, token.Col)
}

func (token *Token) IsTypeKeyword() bool {
	return utils.ArrayContains(token.Type, TYPES_KEYWORDS)
}

func createToken(tokenType TokenType, value string, rowIndex int, colIndex int) Token {
	return Token{
		Type:  tokenType,
		Value: value,
		Row:   rowIndex + 1,
		Col:   colIndex + 1,
	}
}

type Tokenizer struct {
	code      *string
	index     int
	line      int
	lineStart int
}

func NewTokenizer(code *string) Tokenizer {
	tokenizer := Tokenizer{
		code:      code,
		index:     0,
		line:      0,
		lineStart: 0,
	}

	return tokenizer
}

// Gets next token in code. When reaching EOF, all cursors will be reset and starts tokenizing from the beginning.
func (tokenizer *Tokenizer) NextToken() (*Token, error) {

	for tokenizer.index < len(*tokenizer.code) {
		currentChar := (*tokenizer.code)[tokenizer.index]

		if currentChar == '@' {
			// Comment
			tokenizer.skipComment()
		} else if currentChar == '\n' {
			// Skip new line
			tokenizer.index++
			tokenizer.line++
			tokenizer.lineStart = tokenizer.index
		} else if unicode.IsSpace(rune(currentChar)) || currentChar == '\t' {
			// Skip whitespaces
			tokenizer.index++
		} else if utils.ArrayContains(currentChar, OPERATORS_FIRSTS) {
			// Read operator or assignment
			value, tokenType, new_i := tokenizer.readOperatorOrAssign()
			token := createToken(tokenType, value, tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index = new_i
			return &token, nil
		} else if currentChar == '(' {
			token := createToken(LPAR, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == ')' {
			token := createToken(RPAR, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == '{' {
			token := createToken(LBRACE, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == '}' {
			token := createToken(RBRACE, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == '[' {
			token := createToken(LBRACKET, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == ']' {
			token := createToken(RBRACKET, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if unicode.IsLetter(rune(currentChar)) || currentChar == '_' {
			value, new_i := tokenizer.readIdentifier()
			var tokenType TokenType
			if keywordTokenType, ok := KEYWORDS_MAP[value]; ok {
				tokenType = keywordTokenType
			} else {
				tokenType = IDENTIFIER
			}
			token := createToken(tokenType, value, tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index = new_i
			return &token, nil
		} else if unicode.IsNumber(rune(currentChar)) {
			value, new_i := tokenizer.readLiteralNumber()
			token := createToken(LITERAL_INT, value, tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index = new_i
			return &token, nil
		} else if currentChar == ';' {
			token := createToken(SEMICOLON, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == ':' {
			token := createToken(COLON, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else if currentChar == ',' {
			token := createToken(COMMA, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		} else {
			// Read illegal if nothing matches
			token := createToken(ILLEGAL, string(currentChar), tokenizer.line, tokenizer.index-tokenizer.lineStart)
			tokenizer.index++
			return &token, nil
		}
	}

	tokenizer.resetCusrors()
	eofToken := createToken(EOF, "", tokenizer.line, tokenizer.index-tokenizer.lineStart)
	return &eofToken, nil
}

func (tokenizer *Tokenizer) resetCusrors() {
	tokenizer.index = 0
	tokenizer.line = 0
	tokenizer.lineStart = 0
}

func (tokenizer *Tokenizer) skipComment() {
	i := tokenizer.index + 1
	skippedChars := 1

	lines := 0
	lineStart := 0

	for {
		peekedChar := (*tokenizer.code)[i]
		if peekedChar == '\n' {
			lines++
			lineStart = i + 1
			break
		}
		i++
		skippedChars++
	}
	tokenizer.index += skippedChars + 1
	tokenizer.line += lines
	if lineStart != 0 {
		tokenizer.lineStart = lineStart
	}
}

func (tokenizer *Tokenizer) readIdentifier() (string, int) {
	i := tokenizer.index
	buffer := string((*tokenizer.code)[i])
	i++

	for i < len(*tokenizer.code) {
		currentChar := (*tokenizer.code)[i]
		if !unicode.IsLetter(rune(currentChar)) && !unicode.IsNumber(rune(currentChar)) && currentChar != '_' {
			break
		}
		buffer = buffer + string(currentChar)
		i++
	}

	return buffer, i
}

func (tokenizer *Tokenizer) readLiteralNumber() (string, int) {
	i := tokenizer.index
	currentChar := (*tokenizer.code)[tokenizer.index]
	buffer := string(currentChar)

	i++
	for {
		if len(*tokenizer.code) == i {
			break
		}
		currentChar = (*tokenizer.code)[i]
		if currentChar != '.' && !unicode.IsNumber(rune(currentChar)) {
			break
		}
		buffer += string(currentChar)
		i++
	}
	return buffer, i
}

func (tokenizer *Tokenizer) readLiteral() (string, int) {
	i := tokenizer.index
	currentChar := (*tokenizer.code)[tokenizer.index]
	buffer := string(currentChar)

	i++
	if unicode.IsNumber(rune(currentChar)) {
		for {
			if len(*tokenizer.code) == i {
				break
			}
			currentChar = (*tokenizer.code)[i]
			if currentChar != '.' && !unicode.IsNumber(rune(currentChar)) {
				break
			}
			buffer += string(currentChar)
			i++
		}
	} else {
		for {
			if len(*tokenizer.code) == i {
				break
			}
			currentChar = (*tokenizer.code)[i]
			buffer += string(currentChar)
			if currentChar == '"' {
				break
			}
			i++
		}
		i++ // Because the loop stops at the closing " not after
	}

	return buffer, i
}

func (tokenizer *Tokenizer) readOperatorOrAssign() (string, TokenType, int) {
	buffer := string((*tokenizer.code)[tokenizer.index])

	i := tokenizer.index + 1

	if i < len(*tokenizer.code) {
		buffer = buffer + string((*tokenizer.code)[i])

		if doubleCharOp, ok := OPERATORS_ASSIGN_MAP[buffer]; ok { // TODO: Fix this shit
			return buffer, doubleCharOp, tokenizer.index + 2
		}
		buffer = string(buffer[0])
	}
	return buffer, OPERATORS_ASSIGN_MAP[buffer], i
}

func (tokenizer *Tokenizer) isIdentifierKeyword(value string) bool {
	return utils.ArrayContains(value, KEYWORDS)
}

func (tokenizer *Tokenizer) updateCursorForNewLine() {
	tokenizer.index++
	tokenizer.line++
	tokenizer.lineStart = tokenizer.index
}
