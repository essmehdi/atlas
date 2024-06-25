package compiler

import (
	"atlas/parser"
	"fmt"
)

type Compiler struct {
	instructions Instructions
	constants    []Object
}

func New() Compiler {
	return Compiler{
		instructions: Instructions{},
	}
}

func (compiler *Compiler) Compile(program parser.Node) error {
	switch node := program.(type) {
	case *parser.Program:
		for _, stmt := range node.Statements {
			err := compiler.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *parser.ExpressionStatement:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}
		compiler.emit(POP)
	case *parser.InfixExpression:
		err := compiler.Compile(node.Left)
		if err != nil {
			return err
		}

		err = compiler.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			compiler.emit(ADD)
		case "-":
			compiler.emit(SUB)
		case "*":
			compiler.emit(MUL)
		case "/":
			compiler.emit(DIV)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *parser.UnsignedIntegerLiteralExpression:
		integer := UnsignedInteger{Value: node.Value}
		compiler.emit(CONST, compiler.registerConstant(&integer))
	case *parser.BooleanLiteralExpression:
		if node.Value {
			compiler.emit(TRUE)
		} else {
			compiler.emit(FALSE)
		}
	}
	return nil
}

func (compiler *Compiler) registerConstant(obj Object) int {
	compiler.constants = append(compiler.constants, obj)
	return len(compiler.constants) - 1
}

func (compiler *Compiler) emit(opCode OpCode, operands ...int) int {
	instruction := MakeInstruction(opCode, operands...)
	position := compiler.addInstruction(instruction)
	return position
}

func (compiler *Compiler) addInstruction(instruction []byte) int {
	newInstPosition := len(compiler.instructions)
	compiler.instructions = append(compiler.instructions, instruction...)
	return newInstPosition
}

func (compiler *Compiler) ByteCode() ByteCode {
	return ByteCode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}
