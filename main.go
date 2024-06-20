package main

import (
	"atlas/lexer"
	"atlas/parser"
	"fmt"
)

func main() {
	testCode := `
	var a: bool = false;
	loop !a {
		b=2*5+hi();
	}
	fun hi(msg: int) {
	}
	a = true;
	hi(a);
	`
	// testTokenizer(&testCode)
	testParser(&testCode)
}

func testTokenizer(code *string) {
	tokenizer := lexer.New(code)

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
	p := parser.New(code)
	program := p.Parse()
	program.Print()

	for _, err := range p.Errors {
		fmt.Println(err)
	}
}
