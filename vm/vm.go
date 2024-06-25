package vm

import (
	"atlas/compiler"
	"fmt"
)

var STACK_SIZE int = 2048

var True = &compiler.Boolean{Value: true}
var False = &compiler.Boolean{Value: true}

type VM struct {
	instructions compiler.Instructions
	constants    []compiler.Object
	stack        []compiler.Object
	sp           int
}

func New(byteCode compiler.ByteCode) VM {
	return VM{
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		stack:        make([]compiler.Object, STACK_SIZE),
		sp:           0,
	}
}

func (vm *VM) StackTop() compiler.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) PoppedGhost() compiler.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		operation := compiler.OpCode(vm.instructions[ip])

		switch operation {
		case compiler.CONST:
			constIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case compiler.ADD, compiler.SUB, compiler.MUL, compiler.DIV:
			vm.executeBinaryOp(operation)
		case compiler.TRUE:
			vm.push(True)
		case compiler.FALSE:
			vm.push(False)
		case compiler.POP:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) executeBinaryOp(opCode compiler.OpCode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == compiler.INTEGER || leftType == compiler.UNSIGNED_INTEGER || rightType == compiler.INTEGER || rightType == compiler.UNSIGNED_INTEGER {
		vm.executeBinaryIntegerOp(opCode, left, right)
		return nil
	}
	return fmt.Errorf("cannot do binary operations on operands of type `%s` and `%s`", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOp(opCode compiler.OpCode, left compiler.Object, right compiler.Object) {
	switch opCode {
	case compiler.ADD: // TODO: Refactor this
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.Integer{Value: leftValue + int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				vm.push(&compiler.Integer{Value: leftValue + right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.UnsignedInteger{Value: leftValue + right.(*compiler.UnsignedInteger).Value})
			} else {
				vm.push(&compiler.Integer{Value: int64(leftValue) + right.(*compiler.Integer).Value})
			}
		}
	case compiler.SUB:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.Integer{Value: leftValue - int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				vm.push(&compiler.Integer{Value: leftValue - right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.UnsignedInteger{Value: leftValue - right.(*compiler.UnsignedInteger).Value})
			} else {
				vm.push(&compiler.Integer{Value: int64(leftValue) - right.(*compiler.Integer).Value})
			}
		}
	case compiler.MUL:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.Integer{Value: leftValue * int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				vm.push(&compiler.Integer{Value: leftValue * right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.UnsignedInteger{Value: leftValue * right.(*compiler.UnsignedInteger).Value})
			} else {
				vm.push(&compiler.Integer{Value: int64(leftValue) * right.(*compiler.Integer).Value})
			}
		}
	case compiler.DIV:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.Integer{Value: leftValue / int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				vm.push(&compiler.Integer{Value: leftValue / right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				vm.push(&compiler.UnsignedInteger{Value: leftValue / right.(*compiler.UnsignedInteger).Value})
			} else {
				vm.push(&compiler.Integer{Value: int64(leftValue) / right.(*compiler.Integer).Value})
			}
		}
	}
}

func (vm *VM) push(obj compiler.Object) error {
	if vm.sp >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

func (vm *VM) pop() compiler.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj
}

func (vm *VM) DebugStack() {
	fmt.Println(vm.stack)
}
