package compiler

import (
	"atlas/parser"
	"fmt"
)

type Compiler struct {
	instructions Instructions
	constants    []Object
	symbolTable  *SymbolTable
}

func New() Compiler {
	return Compiler{
		instructions: Instructions{},
		symbolTable:  NewSymbolTable(),
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
	case *parser.StatementsBlock:
		for _, stmt := range node.Statements {
			err := compiler.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *parser.DeclarationStatement:
		err := compiler.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := compiler.symbolTable.Define(node.Name.Value)
		compiler.emit(GLOBAL_SET, symbol.Index)
	case *parser.AssignmentStatement:
		err := compiler.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol, ok := compiler.symbolTable.Resolve(node.Name.Value)
		if ok {
			compiler.emit(GLOBAL_SET, symbol.Index)
		} else {
			return fmt.Errorf("cannot assign new value to undeclared variable `%s`", node.Name.Value)
		}
	case *parser.IfStatement:
		blockEndJumpPositions := []int{}
		lastBlockIndex := len(node.Consequences) - 1
		for i, conseq := range node.Consequences {
			err := compiler.Compile(node.Conditions[0])
			if err != nil {
				return err
			}

			jumpOpPosition := compiler.emit(JNT, 0)

			err = compiler.Compile(conseq)
			if err != nil {
				return err
			}

			if node.Else != nil || i != lastBlockIndex {
				jump := compiler.emit(JUMP, 0)
				blockEndJumpPositions = append(blockEndJumpPositions, jump)
			}

			postConsequence := len(compiler.instructions)
			compiler.changeOperand(jumpOpPosition, postConsequence)
		}
		if node.Else != nil {
			err := compiler.Compile(node.Else)
			if err != nil {
				return err
			}
		}
		for _, blockEndJumpPosition := range blockEndJumpPositions {
			postIf := len(compiler.instructions)
			compiler.changeOperand(blockEndJumpPosition, postIf)
		}
	case *parser.LoopStatement:
		offsetPreCondition := len(compiler.instructions)
		err := compiler.Compile(node.Condition)
		if err != nil {
			return err
		}
		offsetPostCondition := len(compiler.instructions)
		offset := offsetPostCondition - offsetPreCondition

		jumpOpPosition := compiler.emit(JNT, 0)

		err = compiler.Compile(node.Block)
		if err != nil {
			return err
		}

		compiler.emit(JUMP, jumpOpPosition - offset)

		postBlock := len(compiler.instructions)
		compiler.changeOperand(jumpOpPosition, postBlock)
	case *parser.ExpressionStatement:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}
		// compiler.emit(POP)
	case *parser.PrefixExpression:
		err := compiler.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			compiler.emit(BANG)
		case "-":
			compiler.emit(MINUS)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *parser.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			err := compiler.Compile(node.Right)
			if err != nil {
				return err
			}

			err = compiler.Compile(node.Left)
			if err != nil {
				return err
			}

			if node.Operator == "<" {
				compiler.emit(GT)
			} else {
				compiler.emit(GEQ)
			}
			return nil
		} else {
			err := compiler.Compile(node.Left)
			if err != nil {
				return err
			}

			err = compiler.Compile(node.Right)
			if err != nil {
				return err
			}
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
		case "==":
			compiler.emit(EQ)
		case "!=":
			compiler.emit(NEQ)
		case ">":
			compiler.emit(GT)
		case ">=":
			compiler.emit(GEQ)
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
	case *parser.Identifier:
		symbol, ok := compiler.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined symbol %s", node.Value)
		}
		compiler.emit(GLOBAL_GET, symbol.Index)
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

func (c *Compiler) replaceInstruction(position int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[position+i] = newInstruction[i]
	}
}

func (compiler *Compiler) changeOperand(opPosition int, operand int) {
	op := OpCode(compiler.instructions[opPosition])
	newInstruction := MakeInstruction(op, operand)
	compiler.replaceInstruction(opPosition, newInstruction)
}

func (compiler *Compiler) ByteCode() ByteCode {
	return ByteCode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}
