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
	var a = 0;
	in a;
	loop a < 10 {
		a = a + 1;
	}
	return a;
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

	for _, err := range p.Errors {
		fmt.Println(err)
		return
	}
	program.Print()

	compiler := compiler.New()
	err := compiler.Compile(&program)
	if err != nil {
		fmt.Printf("Compilation error: %s\n", err)
	}

	fmt.Println(compiler.ByteCode().Instructions)

	vm := vm.New(compiler.ByteCode())
	err = vm.Run()
	
	print("Stack top: ")
	stackTop := vm.StackTop()
	if stackTop == nil {
		println("nil")
	} else {
		println(stackTop.Inspect())
	}

	print("Stack ghost: ")
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
