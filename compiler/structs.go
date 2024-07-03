package compiler

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instructions []byte

type OpCode byte

const (
	CONST OpCode = iota

	ADD // Adds last two items in stack
	SUB // Subs ...
	MUL // Multiplies ...
	DIV // Divides

	TRUE  // Pushes True to stack
	FALSE // ... False ...

	EQ // Comparison ops of last two items in stack
	NEQ
	GT
	GEQ

	BANG // Prefix modifiers
	MINUS

	JUMP // Branch ops
	JNT

	GLOBAL_SET // Global bindings
	GLOBAL_GET

	IN // Program IO
	OUT

	POP // Pops from stack
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var DEFINITIONS = map[OpCode]*Definition{
	CONST: {"CONST", []int{2}},

	ADD: {"ADD", []int{}},
	SUB: {"SUB", []int{}},
	MUL: {"MUL", []int{}},
	DIV: {"DIV", []int{}},

	TRUE:  {"TRUE", []int{}},
	FALSE: {"FALSE", []int{}},

	EQ:  {"EQ", []int{}},
	NEQ: {"NEQ", []int{}},
	GT:  {"GT", []int{}},
	GEQ: {"GEQ", []int{}},

	BANG:  {"BANG", []int{}},
	MINUS: {"MINUS", []int{}},

	JUMP: {"JUMP", []int{2}},
	JNT:  {"JNT", []int{2}},

	GLOBAL_SET: {"GLOBAL_SET", []int{2}},
	GLOBAL_GET: {"GLOBAL_GET", []int{2}},

	IN: {"IN", []int{2}},
	OUT: {"OUT", []int{}},

	POP: {"POP", []int{}},
}

type ByteCode struct {
	Instructions Instructions
	Constants    []Object
}

func LookupOperation(op byte) (*Definition, error) {
	definition, ok := DEFINITIONS[OpCode(op)]
	if !ok {
		return nil, fmt.Errorf("OpCode %d undefined", op)
	}
	return definition, nil
}

func MakeInstruction(op OpCode, operands ...int) []byte {
	definition, ok := DEFINITIONS[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, width := range definition.OperandWidths {
		instructionLen += width
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for index, operand := range operands {
		width := definition.OperandWidths[index]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}
	return instruction
}

func ReadInstructionOperands(definition *Definition, instruction Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0
	for i, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instruction[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (instruction Instructions) String() string {
	var out strings.Builder
	i := 0
	for i < len(instruction) {
		definition, err := LookupOperation(instruction[i])
		if err != nil {
			fmt.Fprintf(&out, "Error: %s\n", err)
			continue
		}
		operands, read := ReadInstructionOperands(definition, instruction[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, instruction.formatInstruction(definition, operands))
		i += 1 + read
	}
	return out.String()
}

func (instruction Instructions) formatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("Error: operand length %d does not match length in the definition %d\n",
			len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("Error: unhandled operands count for %s\n", def.Name)
}
