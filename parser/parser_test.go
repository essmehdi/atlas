package parser

import (
	"testing"
)

func TestParseDeclarationStatement(t *testing.T) {
	input := "var x = 5;"
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*DeclarationStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *DeclarationStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}

	if stmt.Type != INFERED {
		t.Errorf("stmt.Type not INFERED. got=%v", stmt.Type)
	}

	literal, ok := stmt.Value.(*UnsignedIntegerLiteralExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *UnsignedIntegerLiteralExpression. got=%T", stmt.Value)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not '5'. got=%d", literal.Value)
	}
}

func TestParseInputStatement(t *testing.T) {
	input := "in x;"
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*InputStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *InputStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}
}

func TestParseAssignmentStatement(t *testing.T) {
	input := "x = 10;"
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}

	literal, ok := stmt.Value.(*UnsignedIntegerLiteralExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *UnsignedIntegerLiteralExpression. got=%T", stmt.Value)
	}

	if literal.Value != 10 {
		t.Errorf("literal.Value not '10'. got=%d", literal.Value)
	}
}

func TestParseIfStatement(t *testing.T) {
	input := `
	if x > 5 {
		y = 10;
	} else if x < 3 {
		y = 5;
	} else {
		y = 1;
	}`
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *IfStatement. got=%T", program.Statements[0])
	}

	if len(stmt.Conditions) != 2 {
		t.Fatalf("statement does not have 2 conditions. got=%d", len(stmt.Conditions))
	}

	if len(stmt.Consequences) != 2 {
		t.Fatalf("statement does not have 2 consequences. got=%d", len(stmt.Consequences))
	}

	if stmt.Else == nil {
		t.Fatalf("statement does not have else block")
	}
}

func TestParseLoopStatement(t *testing.T) {
	input := `
	loop x > 0 {
		x = x - 1;
	}`
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*LoopStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *LoopStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatalf("loop statement condition is nil")
	}

	if stmt.Block == nil {
		t.Fatalf("loop statement block is nil")
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	input := `
	fun add(x: int, y: int): int {
		return x + y;
	}`
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*FunctionDeclarationStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *FunctionDeclarationStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.ArgsNames) != 2 {
		t.Fatalf("function has wrong number of parameters. got=%d", len(stmt.ArgsNames))
	}

	if stmt.ArgsNames[0].Value != "x" || stmt.ArgsNames[1].Value != "y" {
		t.Errorf("parameter names are not 'x' and 'y'. got=%s, %s", stmt.ArgsNames[0].Value, stmt.ArgsNames[1].Value)
	}

	if stmt.ArgsTypes[0] != INT || stmt.ArgsTypes[1] != INT {
		t.Errorf("parameter types are not both INT. got=%v, %v", stmt.ArgsTypes[0], stmt.ArgsTypes[1])
	}

	if *stmt.ReturnType != INT {
		t.Errorf("return type is not INT. got=%v", stmt.ReturnType)
	}
}

func TestParseCallExpression(t *testing.T) {
	input := "add(1, 2);"
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *CallExpression. got=%T", stmt.Expression)
	}

	if exp.Function.(*Identifier).Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", exp.Function.(*Identifier).Value)
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("wrong number of arguments. got=%d", len(exp.Arguments))
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := "return x + y;"
	parser := New(&input)
	program := parser.Parse()

	if len(program.Statements) > 1 {
		t.Fatalf("program has too many statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ReturnStatement. got=%T", program.Statements[0])
	}

	if stmt.Expression == nil {
		t.Fatalf("return statement expression is nil")
	}

	_, ok = stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("return statement expression is not *InfixExpression. got=%T", stmt.Expression)
	}
}

func TestParseInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  uint64
		operator   string
		rightValue uint64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		parser := New(&tt.input)
		program := parser.Parse()

		if len(program.Statements) > 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*InfixExpression)
		if !ok {
			t.Fatalf("stmt is not *InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testUnsignedIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if !testUnsignedIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func testUnsignedIntegerLiteral(t *testing.T, il Expression, value uint64) bool {
	integ, ok := il.(*UnsignedIntegerLiteralExpression)
	if !ok {
		t.Errorf("il not *UnsignedIntegerLiteralExpression. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}