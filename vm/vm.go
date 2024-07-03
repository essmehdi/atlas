package vm

import (
	"atlas/compiler"
	"fmt"
	"math"
)

const STACK_SIZE int = 2048
const GLOBALS_SIZE int = 65536

type VM struct {
	instructions compiler.Instructions
	constants    []compiler.Object
	stack        []compiler.Object
	sp           int
	globals      []compiler.Object
}

func New(byteCode compiler.ByteCode) VM {
	return VM{
		instructions: byteCode.Instructions,
		constants:    byteCode.Constants,
		stack:        make([]compiler.Object, STACK_SIZE),
		sp:           0,
		globals:      make([]compiler.Object, GLOBALS_SIZE),
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
		
		var err error

		switch operation {
		case compiler.CONST:
			constIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err = vm.push(vm.constants[constIndex])
		case compiler.EQ, compiler.NEQ, compiler.GT, compiler.GEQ:
			err = vm.executeComparison(operation)
		case compiler.ADD, compiler.SUB, compiler.MUL, compiler.DIV:
			err = vm.executeBinaryOp(operation)
		case compiler.BANG:
			err = vm.executeBangOperation()
		case compiler.MINUS:
			err = vm.executeMinusOperation()
		case compiler.TRUE:
			err = vm.push(compiler.True)
		case compiler.FALSE:
			err = vm.push(compiler.False)
		case compiler.JUMP:
			targetInstruction := compiler.ReadUint16(vm.instructions[ip+1:])
			ip = int(targetInstruction - 1)
		case compiler.JNT:
			targetInstruction := compiler.ReadUint16(vm.instructions[ip+1:])
			conditionEval := vm.pop()
			if !conditionEval.(*compiler.Boolean).Value {
				ip = int(targetInstruction - 1)
			} else {
				ip += 2
			}
		case compiler.GLOBAL_SET:
			globalIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[globalIndex] = vm.pop()
		case compiler.GLOBAL_GET:
			globalIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case compiler.IN:
			globalIndex := compiler.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			global := vm.globals[globalIndex]
			switch global.Type() {
			case compiler.UNSIGNED_INTEGER:
				number := &compiler.UnsignedInteger{}
				fmt.Scanln(&number.Value)
				vm.globals[globalIndex] = number
			case compiler.INTEGER:
				number := &compiler.Integer{}
				fmt.Scanln(&number.Value)
				vm.globals[globalIndex] = number
			case compiler.BOOLEAN:
				number := &compiler.Boolean{}
				fmt.Scanln(&number.Value)
				vm.globals[globalIndex] = number
			}
		case compiler.OUT:
			output := vm.pop()
			fmt.Println(output.Inspect())
		case compiler.POP:
			vm.pop()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) executeBangOperation() error {
	operand := vm.pop()
	switch operand {
	// case True:
	// return vm.push(False)
	case compiler.False:
		return vm.push(compiler.True)
	default:
		return vm.push(compiler.False)
	}
}

func (vm *VM) executeMinusOperation() error {
	operand := vm.pop()
	switch oper := operand.(type) {
	case *compiler.Integer:
		return vm.push(&compiler.Integer{Value: -oper.Value})
	case *compiler.UnsignedInteger:
		if oper.Value > math.MaxInt64 {
			return fmt.Errorf("overflow error when trying to apply `-` operator on `%d`", oper.Value)
		}
		return vm.push(&compiler.Integer{Value: -int64(oper.Value)})
	default:
		return vm.push(compiler.False)
	}
}

func (vm *VM) executeComparison(opCode compiler.OpCode) error {
	right := vm.pop()
	left := vm.pop()
	if compiler.IsObjectNumber(left) && compiler.IsObjectNumber(right) {
		return vm.executeIntegerComparison(opCode, left, right)
	}

	switch opCode {
	case compiler.EQ:
		return vm.push(compiler.ParseBooleanFromNative(left == right))
	case compiler.NEQ:
		return vm.push(compiler.ParseBooleanFromNative(left != right))
	default:
		return fmt.Errorf("could not apply operator: %d on (%s %s)", opCode, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(opCode compiler.OpCode, left compiler.Object, right compiler.Object) error {
	leftType := left.Type()
	rightType := right.Type()
	if leftType == compiler.INTEGER && rightType == compiler.UNSIGNED_INTEGER {
		leftValue := left.(*compiler.Integer).Value
		rightValue := right.(*compiler.UnsignedInteger).Value
		switch opCode {
		case compiler.EQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue == int64(rightValue)))
		case compiler.NEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue != int64(rightValue)))
		case compiler.GT:
			return vm.push(compiler.ParseBooleanFromNative(leftValue > int64(rightValue)))
		case compiler.GEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue >= int64(rightValue)))
		}
	} else if leftType == compiler.UNSIGNED_INTEGER && rightType == compiler.INTEGER {
		leftValue := left.(*compiler.UnsignedInteger).Value
		rightValue := right.(*compiler.Integer).Value
		switch opCode {
		case compiler.EQ:
			return vm.push(compiler.ParseBooleanFromNative(int64(leftValue) == rightValue))
		case compiler.NEQ:
			return vm.push(compiler.ParseBooleanFromNative(int64(leftValue) != rightValue))
		case compiler.GT:
			return vm.push(compiler.ParseBooleanFromNative(int64(leftValue) > rightValue))
		case compiler.GEQ:
			return vm.push(compiler.ParseBooleanFromNative(int64(leftValue) >= rightValue))
		}
	} else if leftType == compiler.UNSIGNED_INTEGER && rightType == compiler.UNSIGNED_INTEGER {
		leftValue := left.(*compiler.UnsignedInteger).Value
		rightValue := right.(*compiler.UnsignedInteger).Value
		switch opCode {
		case compiler.EQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue == rightValue))
		case compiler.NEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue != rightValue))
		case compiler.GT:
			return vm.push(compiler.ParseBooleanFromNative(leftValue > rightValue))
		case compiler.GEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue >= rightValue))
		}
	} else /*if leftType == compiler.INTEGER && rightType == compiler.INTEGER*/ {
		leftValue := left.(*compiler.Integer).Value
		rightValue := right.(*compiler.Integer).Value
		switch opCode {
		case compiler.EQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue == rightValue))
		case compiler.NEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue != rightValue))
		case compiler.GT:
			return vm.push(compiler.ParseBooleanFromNative(leftValue > rightValue))
		case compiler.GEQ:
			return vm.push(compiler.ParseBooleanFromNative(leftValue >= rightValue))
		}
	}
	return fmt.Errorf("unknown operator: %d", opCode)
}

func (vm *VM) executeBinaryOp(opCode compiler.OpCode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == compiler.INTEGER || leftType == compiler.UNSIGNED_INTEGER || rightType == compiler.INTEGER || rightType == compiler.UNSIGNED_INTEGER {
		return vm.executeBinaryIntegerOp(opCode, left, right)
	}
	return fmt.Errorf("cannot do binary operations on operands of type `%s` and `%s`", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOp(opCode compiler.OpCode, left compiler.Object, right compiler.Object) error {
	switch opCode {
	case compiler.ADD: // TODO: Refactor this
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.Integer{Value: leftValue + int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				return vm.push(&compiler.Integer{Value: leftValue + right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.UnsignedInteger{Value: leftValue + right.(*compiler.UnsignedInteger).Value})
			} else {
				return vm.push(&compiler.Integer{Value: int64(leftValue) + right.(*compiler.Integer).Value})
			}
		}
	case compiler.SUB:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.Integer{Value: leftValue - int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				return vm.push(&compiler.Integer{Value: leftValue - right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.UnsignedInteger{Value: leftValue - right.(*compiler.UnsignedInteger).Value})
			} else {
				return vm.push(&compiler.Integer{Value: int64(leftValue) - right.(*compiler.Integer).Value})
			}
		}
	case compiler.MUL:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.Integer{Value: leftValue * int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				return vm.push(&compiler.Integer{Value: leftValue * right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.UnsignedInteger{Value: leftValue * right.(*compiler.UnsignedInteger).Value})
			} else {
				return vm.push(&compiler.Integer{Value: int64(leftValue) * right.(*compiler.Integer).Value})
			}
		}
	case compiler.DIV:
		if left.Type() == compiler.INTEGER {
			leftValue := left.(*compiler.Integer).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.Integer{Value: leftValue / int64(right.(*compiler.UnsignedInteger).Value)})
			} else {
				return vm.push(&compiler.Integer{Value: leftValue / right.(*compiler.Integer).Value})
			}
		} else {
			leftValue := left.(*compiler.UnsignedInteger).Value
			if right.Type() == compiler.UNSIGNED_INTEGER {
				return vm.push(&compiler.UnsignedInteger{Value: leftValue / right.(*compiler.UnsignedInteger).Value})
			} else {
				return vm.push(&compiler.Integer{Value: int64(leftValue) / right.(*compiler.Integer).Value})
			}
		}
	}
	return fmt.Errorf("could not do binary integer op: %d", opCode)
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
