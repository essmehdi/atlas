package main

import (
	"atlas/lexer"
	"atlas/parser"
	"fmt"
)

func main() {
	testCode := "var a: bool = false;loop !a {b=2*5;} fun hi() {} a = true;"
	// testTokenizer(&testCode)
	testParser(&testCode)
}

func testTokenizer(code *string) {
	tokenizer := lexer.NewTokenizer(code)

	for {
		token, err := tokenizer.NextToken()
		if err != nil {
			fmt.Printf("Tokenizer error: %v\n", err)
			break
		}
		if token.Type == lexer.EOF {
			break
		}
		fmt.Println(token)
	}
}

func testParser(code *string) {
	parser := parser.NewParser(code);
	program := parser.Parse()
	program.Print()

	for _, err := range parser.Errors {
		fmt.Println(err)
	}
}