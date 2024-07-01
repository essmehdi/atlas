package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (symbolTable *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: symbolTable.numDefinitions, Scope: GlobalScope}
	symbolTable.store[name] = symbol
	symbolTable.numDefinitions++
	return symbol
}

func (symbolTable *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := symbolTable.store[name]
	return obj, ok
}
