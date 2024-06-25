package compiler

import "fmt"

type ObjectType string

const (
	INTEGER 			= "INTEGER"
	UNSIGNED_INTEGER 	= "UNSIGNED_INTEGER"
	BOOLEAN				= "BOOLEAN"
	NULL				= "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
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

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (uinteger *Boolean) Inspect() string {
	return fmt.Sprint(uinteger.Value)
}