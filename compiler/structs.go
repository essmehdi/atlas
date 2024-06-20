package compiler

type OpCode = byte

type Instruction = byte

const (
	ADD OpCode = iota
	SUB
	MUL
	DIV

	AND
	OR
	XOR

	LOAD
	STORE

	JUMP
	JZ
	JNZ

	IN
	OUT

	HALT
)

func MakeInstruction(op OpCode, arg byte) byte {
	return (op << 4) | arg
}

func GetOpCode(instruction Instruction) OpCode {
	return instruction >> 4
}

func GetArg(instruction Instruction) byte {
	return instruction & 0xF // Bit mask
}