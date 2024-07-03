package compiler

import "fmt"

type ObjectType string

const (
	INTEGER 			= "INTEGER"
	UNSIGNED_INTEGER 	= "UNSIGNED_INTEGER"
	BOOLEAN				= "BOOLEAN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

func IsObjectNumber(object Object) bool {
	return object.Type() == UNSIGNED_INTEGER || object.Type() == INTEGER
}

// Integers object

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType {
	return INTEGER
}

func (integer *Integer) Inspect() string {
	return fmt.Sprint(integer.Value)
}

// Unsigned integer object

type UnsignedInteger struct {
	Value uint64
}

func (uinteger *UnsignedInteger) Type() ObjectType {
	return UNSIGNED_INTEGER
}

func (uinteger *UnsignedInteger) Inspect() string {
	return fmt.Sprint(uinteger.Value)
}

// Boolean object

var True = &Boolean{Value: true}
var False = &Boolean{Value: false}

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprint(boolean.Value)
}

func ParseBooleanFromNative(nativeBool bool) *Boolean {
	if nativeBool {
		return True
	}
	return False
}