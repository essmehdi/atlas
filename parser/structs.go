package parser

import (
	"atlas/lexer"
	"atlas/utils"
	"fmt"
	"strings"
)

// Operator precedance
const (
	_ int = iota
	LOWEST
	BITWISE     // ~ &
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var PRECEDENCE_MAP = map[lexer.TokenType]int{
	lexer.BIT_AND:     	BITWISE,
	lexer.BIT_OR:      	BITWISE,
	lexer.EQ:          	EQUALS,
	lexer.NEQ:         	EQUALS,
	lexer.LOGICAL_AND: 	EQUALS,
	lexer.LOGICAL_OR:  	EQUALS,
	lexer.LT:          	LESSGREATER,
	lexer.GT:          	LESSGREATER,
	lexer.LEQ:         	LESSGREATER,
	lexer.GEQ:         	LESSGREATER,
	lexer.PLUS:        	SUM,
	lexer.MINUS:       	SUM,
	lexer.MULTIPLY:    	PRODUCT,
	lexer.DIVIDE:      	PRODUCT,
	lexer.BANG:        	PREFIX,
	lexer.BIT_NOT:     	PREFIX,
	lexer.LPAR:			CALL,
}

type DataType int
const (
	INFERED DataType = iota
	INT
	UINT
	BOOL
)

func (dataType DataType) String() string {
	return [...]string{
		"Infered",
		"Integer",
		"Unsigned integer",
		"Boolean",
	}[dataType]
}

var DATA_TYPE_MAP = map[lexer.TokenType]DataType{
	lexer.TYPE_INT: INT,
	lexer.TYPE_UINT: UINT,
	lexer.TYPE_BOOL: BOOL,
}

type Node interface {
	GetToken() *lexer.Token      // Gets the initiating token
	StringRepr(level int) string // Gets the Node string representation
}

type Program struct {
	statements []Statement
}

func newProgram() Program {
	return Program{
		statements: []Statement{},
	}
}

func (program *Program) Print() {
	for _, stmt := range program.statements {
		fmt.Println(stmt.StringRepr(0))
		fmt.Println()
	}
}

func (program *Program) addStatement(statement Statement) {
	program.statements = append(program.statements, statement)
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Declaration: var a = 5;

type DeclarationStatement struct {
	Token 	*lexer.Token
	Name  	*Identifier
	Type	DataType		// For now, only builtin types are accepted
	Value 	Expression
}

func (decl *DeclarationStatement) statementNode() {}

func (decl *DeclarationStatement) GetToken() *lexer.Token {
	return decl.Token
}

func (decl *DeclarationStatement) StringRepr(level int) string {
	if decl == nil {
		return ""
	}
	nameRepr := ""
	if decl.Name != nil {
		nameRepr = decl.Name.StringRepr(level + 1)
	}

	valueRepr := ""
	if decl.Value != nil {
		valueRepr = decl.Value.StringRepr(level + 1)
	}
	
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("DeclarationStatement\nName:\n%s\nValue:\n%s\nType: %s", nameRepr, valueRepr, decl.Type),
	)
}

// Assignment statement: a = 9

type AssignmentStatement struct {
	Token 	*lexer.Token
	Name  	*Identifier
	Value 	Expression
}

func (assign *AssignmentStatement) statementNode() {}

func (assign *AssignmentStatement) GetToken() *lexer.Token {
	return assign.Token
}

func (assign *AssignmentStatement) StringRepr(level int) string {
	if assign == nil {
		return ""
	}
	nameRepr := ""
	if assign.Name != nil {
		nameRepr = assign.Name.StringRepr(level + 1)
	}

	valueRepr := ""
	if assign.Value != nil {
		valueRepr = assign.Value.StringRepr(level + 1)
	}

	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("AssignmentStatement\nName:\n%s\nValue:\n%s", nameRepr, valueRepr),
	)
}

// If statement: if () {} else {};

type IfStatement struct {
	Token        *lexer.Token
	Conditions   []Expression
	Consequences []*StatementsBlock
	Else         *StatementsBlock
}

func (ifStmt *IfStatement) statementNode() {}

func (ifStmt *IfStatement) GetToken() *lexer.Token {
	return ifStmt.Token
}

func (ifStmt *IfStatement) StringRepr(level int) string {
	if ifStmt == nil {
		return ""
	}
	buffer := ""
	for i, cond := range ifStmt.Conditions {
		buffer += fmt.Sprintf(
			"\nCondition:\n%s\nConsequence:\n%s",
			cond.StringRepr(level+1),
			ifStmt.Consequences[i].StringRepr(level+1),
		)
	}
	if ifStmt.Else != nil {
		buffer += fmt.Sprintf(
			"Else:\n%s",
			ifStmt.Else.StringRepr(level+1),
		)
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("IfStatement%s", buffer),
	)
}

// Statments block

type StatementsBlock struct {
	Token      *lexer.Token
	Statements []Statement
}

func (block *StatementsBlock) GetToken() *lexer.Token {
	return block.Token
}

func (block *StatementsBlock) StringRepr(level int) string {
	if block == nil {
		return ""
	}
	if len(block.Statements) == 0 {
		return "Empty"
	}
	buffer := ""
	for _, stmt := range block.Statements {
		buffer = buffer + stmt.StringRepr(level) + "\n"
	}
	return buffer
}

// Identifier expression: a

type Identifier struct {
	Token *lexer.Token
	Value string
}

func (iden *Identifier) expressionNode() {}

func (iden *Identifier) GetToken() *lexer.Token {
	return iden.Token
}

func (iden *Identifier) StringRepr(level int) string {
	if iden == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("Identifier: %s", iden.Value),
	)
}

/*
	Unsigned integer literal expression: 5

	Signed literals are not parsed. They will be detected until compilation with	
	the prefix modifier MINUS (-)
*/ 

type UnsignedIntegerLiteralExpression struct {
	Token *lexer.Token
	Value uint64
}

func (liter *UnsignedIntegerLiteralExpression) expressionNode() {}

func (liter *UnsignedIntegerLiteralExpression) GetToken() *lexer.Token {
	return liter.Token
}

func (liter *UnsignedIntegerLiteralExpression) StringRepr(level int) string {
	if liter == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("IntegerLiteral: %d", liter.Value),
	)
}

// Boolean literal expression: false

type BooleanLiteralExpression struct {
	Token *lexer.Token
	Value bool
}

func (boolean *BooleanLiteralExpression) expressionNode() {}

func (boolean *BooleanLiteralExpression) GetToken() *lexer.Token {
	return boolean.Token
}

func (boolean *BooleanLiteralExpression) StringRepr(level int) string {
	if boolean == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("BooleanLiteral: %t", boolean.Value),
	)
}

// Prefix expression: !condition

type PrefixExpression struct {
	Token    *lexer.Token
	Operator string
	Right    Expression
}

func (preExp *PrefixExpression) expressionNode() {}

func (preExp *PrefixExpression) GetToken() *lexer.Token {
	return preExp.Token
}

func (preExp *PrefixExpression) StringRepr(level int) string {
	if preExp == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("PrefixExpression\nOperator: %s\nRight:\n%s", preExp.Operator, preExp.Right.StringRepr(level+1)),
	)
}

// Infix expression: a + b

type InfixExpression struct {
	Token    *lexer.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (infixExp *InfixExpression) expressionNode() {}

func (infixExp *InfixExpression) GetToken() *lexer.Token {
	return infixExp.Token
}

func (infixExp *InfixExpression) StringRepr(level int) string {
	if infixExp == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("InfixExpression\nOperator: %s\nLeft:\n%s\nRight:\n%s", infixExp.Operator, infixExp.Left.StringRepr(level+1), infixExp.Right.StringRepr(level+1)),
	)
}

// Loop expression: loop a > 10 {...}

type LoopStatement struct {
	Token     *lexer.Token
	Condition Expression
	Block      *StatementsBlock
}

func (loop *LoopStatement) statementNode() {}

func (loop *LoopStatement) GetToken() *lexer.Token {
	return loop.Token
}

func (loop *LoopStatement) StringRepr(level int) string {
	if loop == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("LoopStatement:\nCondition:\n%s\nLoop:\n%s", loop.Condition.StringRepr(level+1), loop.Block.StringRepr(level+1)),
	)
}

// Function definition: fun hello() {}

type FunctionDeclarationStatement struct {
	Token 			*lexer.Token
	Name			*Identifier
	ArgsNames		[]*Identifier
	ArgsTypes		[]DataType
	Body			*StatementsBlock
}

func (fun *FunctionDeclarationStatement) statementNode() {}

func (fun *FunctionDeclarationStatement) GetToken() *lexer.Token {
	return fun.Token
}

func (fun *FunctionDeclarationStatement) StringRepr(level int) string {
	if fun == nil {
		return ""
	}
	argsStr := ""
	for i, arg := range fun.ArgsNames {
		if i > 0 {
			argsStr += ","
		}
		argsStr += arg.Value + ": " + fun.ArgsTypes[i].String()
	}

	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("FunctionDeclarationStatement:\nName:\n%s\nArgs: %s\nBody:\n%s", fun.Name.StringRepr(level + 1), argsStr, fun.Body.StringRepr(level + 1)),
	)
}

// Expression statement: hello();

type ExpressionStatement struct {
	Token 		*lexer.Token
	Expression	Expression
}

func (expr *ExpressionStatement) statementNode() {}

func (expr *ExpressionStatement) GetToken() *lexer.Token {
	return expr.Token
}

func (expr *ExpressionStatement) StringRepr(level int) string {
	if expr == nil {
		return ""
	}
	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("ExpressionStatement:\n%s", expr.Expression.StringRepr(level + 1)),
	)
}

// Call expression: hello()

type CallExpression struct {
	Token 		*lexer.Token
	Function 	Expression
	Arguments 	[]Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) GetToken() *lexer.Token { return ce.Token }

func (ce *CallExpression) StringRepr(level int) string {
	if ce == nil {
		return ""
	}

	var argsBuilder strings.Builder
	for _, arg := range ce.Arguments {
		argsBuilder.WriteString(arg.StringRepr(level + 1))
		argsBuilder.WriteRune('\n')
	}

	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("CallExpression:\nName:\n%s\nArgs:\n%s", ce.Function.StringRepr(level + 1), argsBuilder.String()),
	)
}

// Return statement

type ReturnStatement struct {
	Token		*lexer.Token
	Expression	Expression
}

func (ret *ReturnStatement) statementNode() {}

func (ret *ReturnStatement) GetToken() *lexer.Token { return ret.Token }

func (ret *ReturnStatement) StringRepr(level int) string {
	if ret == nil {
		return ""
	}

	return utils.IndentStringByLevel(
		level,
		fmt.Sprintf("ReturnStatement:\nExpression:\n%s", ret.Expression.StringRepr(level + 1)),
	)
}
