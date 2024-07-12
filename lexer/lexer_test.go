package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	code := `var a = 0;
in a;

var result = 1;
loop a > 1 {
	result = result * a;
	a = a - 1;
}

return result;`

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{VAR, "var"},
		{IDENTIFIER, "a"},
		{ASSIGN, "="},
		{LITERAL_INT, "0"},
		{SEMICOLON, ";"},
		{IN, "in"},
		{IDENTIFIER, "a"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "result"},
		{ASSIGN, "="},
		{LITERAL_INT, "1"},
		{SEMICOLON, ";"},
		{LOOP, "loop"},
		{IDENTIFIER, "a"},
		{GT, ">"},
		{LITERAL_INT, "1"},
		{LBRACE, "{"},
		{IDENTIFIER, "result"},
		{ASSIGN, "="},
		{IDENTIFIER, "result"},
		{MULTIPLY, "*"},
		{IDENTIFIER, "a"},
		{SEMICOLON, ";"},
		{IDENTIFIER, "a"},
		{ASSIGN, "="},
		{IDENTIFIER, "a"},
		{MINUS, "-"},
		{LITERAL_INT, "1"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{RETURN, "return"},
		{IDENTIFIER, "result"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	tokenizer := New(&code)

	for i, exp := range expected {
		token, err := tokenizer.NextToken()
		if err != nil {
			t.Fatalf("Error getting next token: %v", err)
		}

		if token.Type != exp.tokenType {
			t.Errorf("Test case %d: expected token type %v, got %v", i, exp.tokenType, token.Type)
		}

		if token.Value != exp.value {
			t.Errorf("Test case %d: expected token value '%s', got '%s'", i, exp.value, token.Value)
		}
	}
}

func TestLexerWithComments(t *testing.T) {
	code := `var x = 10; @ This is a comment
@ This is another comment
x = x + 5; @ Inline comment`

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{VAR, "var"},
		{IDENTIFIER, "x"},
		{ASSIGN, "="},
		{LITERAL_INT, "10"},
		{SEMICOLON, ";"},
		{IDENTIFIER, "x"},
		{ASSIGN, "="},
		{IDENTIFIER, "x"},
		{PLUS, "+"},
		{LITERAL_INT, "5"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	tokenizer := New(&code)

	for i, exp := range expected {
		token, err := tokenizer.NextToken()
		if err != nil {
			t.Fatalf("Error getting next token: %v", err)
		}

		if token.Type != exp.tokenType {
			t.Errorf("Test case %d: expected token type %v, got %v", i, exp.tokenType, token.Type)
		}

		if token.Value != exp.value {
			t.Errorf("Test case %d: expected token value '%s', got '%s'", i, exp.value, token.Value)
		}
	}
}

func TestLexerIllegalTokens(t *testing.T) {
	code := `var x = 10;
$illegal_token
x = 20;`

	expected := []struct {
		tokenType TokenType
		value     string
	}{
		{VAR, "var"},
		{IDENTIFIER, "x"},
		{ASSIGN, "="},
		{LITERAL_INT, "10"},
		{SEMICOLON, ";"},
		{ILLEGAL, "$"},
		{IDENTIFIER, "illegal_token"},
		{IDENTIFIER, "x"},
		{ASSIGN, "="},
		{LITERAL_INT, "20"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	tokenizer := New(&code)

	for i, exp := range expected {
		token, err := tokenizer.NextToken()
		if err != nil {
			t.Fatalf("Error getting next token: %v", err)
		}

		if token.Type != exp.tokenType {
			t.Errorf("Test case %d: expected token type %v, got %v", i, exp.tokenType, token.Type)
		}

		if token.Value != exp.value {
			t.Errorf("Test case %d: expected token value '%s', got '%s'", i, exp.value, token.Value)
		}
	}
}

func TestLexerAdvanced(t *testing.T) {
	input := `
	fun fibonacci(n: int): int {
		if n <= 1 {
			return n;
		}
		return fibonacci(n - 1) + fibonacci(n - 2);
	}

	var result = fibonacci(10);
	@ This is a comment
	var flag = true;
	var unsignedNum: uint = 4294967295;
	var bitwise = 5 & 3 | 2;
	`

	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{FUN, "fun"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{COLON, ":"},
		{TYPE_INT, "int"},
		{RPAR, ")"},
		{COLON, ":"},
		{TYPE_INT, "int"},
		{LBRACE, "{"},
		{IF, "if"},
		{IDENTIFIER, "n"},
		{LEQ, "<="},
		{LITERAL_INT, "1"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENTIFIER, "n"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{RETURN, "return"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{MINUS, "-"},
		{LITERAL_INT, "1"},
		{RPAR, ")"},
		{PLUS, "+"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{MINUS, "-"},
		{LITERAL_INT, "2"},
		{RPAR, ")"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{VAR, "var"},
		{IDENTIFIER, "result"},
		{ASSIGN, "="},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{LITERAL_INT, "10"},
		{RPAR, ")"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "flag"},
		{ASSIGN, "="},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "unsignedNum"},
		{COLON, ":"},
		{TYPE_UINT, "uint"},
		{ASSIGN, "="},
		{LITERAL_INT, "4294967295"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "bitwise"},
		{ASSIGN, "="},
		{LITERAL_INT, "5"},
		{BIT_AND, "&"},
		{LITERAL_INT, "3"},
		{BIT_OR, "|"},
		{LITERAL_INT, "2"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	lexer := New(&input)

	for i, tt := range tests {
		tok, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("test %d - NextToken() returned error: %v", i, err)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("test %d - literal wrong. expected=%q, got=%q", i, tt.expectedValue, tok.Value)
		}
	}
}

func TestLexerFromFileWithAdvancedCode(t *testing.T) {
	tests := []struct {
		expectedType  TokenType
		expectedValue string
	}{
		{FUN, "fun"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{COLON, ":"},
		{TYPE_INT, "int"},
		{RPAR, ")"},
		{COLON, ":"},
		{TYPE_INT, "int"},
		{LBRACE, "{"},
		{IF, "if"},
		{IDENTIFIER, "n"},
		{LEQ, "<="},
		{LITERAL_INT, "1"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENTIFIER, "n"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{RETURN, "return"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{MINUS, "-"},
		{LITERAL_INT, "1"},
		{RPAR, ")"},
		{PLUS, "+"},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{IDENTIFIER, "n"},
		{MINUS, "-"},
		{LITERAL_INT, "2"},
		{RPAR, ")"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{VAR, "var"},
		{IDENTIFIER, "result"},
		{ASSIGN, "="},
		{IDENTIFIER, "fibonacci"},
		{LPAR, "("},
		{LITERAL_INT, "10"},
		{RPAR, ")"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "flag"},
		{ASSIGN, "="},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "unsignedNum"},
		{COLON, ":"},
		{TYPE_UINT, "uint"},
		{ASSIGN, "="},
		{LITERAL_INT, "4294967295"},
		{SEMICOLON, ";"},
		{VAR, "var"},
		{IDENTIFIER, "bitwise"},
		{ASSIGN, "="},
		{LITERAL_INT, "5"},
		{BIT_AND, "&"},
		{LITERAL_INT, "3"},
		{BIT_OR, "|"},
		{LITERAL_INT, "2"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	lexer, err := NewFromFile("../tests/example.atl")
	if err != nil {
		t.Fatalf("Failed to create lexer: %s", err)
	}

	for i, tt := range tests {
		tok, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("test %d - NextToken() returned error: %v", i, err)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedValue {
			t.Fatalf("test %d - literal wrong. expected=%q, got=%q", i, tt.expectedValue, tok.Value)
		}
	}
}