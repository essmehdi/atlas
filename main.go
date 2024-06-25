package main

import (
	"atlas/compiler"
	"atlas/lexer"
	"atlas/parser"
	"atlas/vm"
	"fmt"
)

func main() {
	// testCode := `
	// var a: bool = false;
	// loop !a {
	// 	b=2*5+hi();
	// }
	// fun hi(msg: int) {
	// }
	// a = true;
	// hi(a);
	// `
	testCode := `
	1+2*3-6;
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

	compiler := compiler.New()
	compiler.Compile(&program)

	fmt.Println(compiler.ByteCode().Instructions)

	vm := vm.New(compiler.ByteCode())
	err := vm.Run()
	
	stackTop := vm.StackTop()
	if stackTop == nil {
		println("nil")
	} else {
		println(stackTop.Inspect())
	}
	stackGhost := vm.PoppedGhost()
	if stackGhost == nil {
		println("nil")
	} else {
		println(stackGhost.Inspect())
	}

	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
	}
}
