package parser

import (
	"atlas/lexer"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	tokenizer    lexer.Tokenizer
	currentToken *lexer.Token
	peekToken    *lexer.Token

	Errors []string

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func setupParser(parser *Parser) {
	parser.nextToken()
	parser.nextToken() // Load current & peek

	parser.registerPrefixParser(lexer.IDENTIFIER, parser.parseIdentifierExpression)
	parser.registerPrefixParser(lexer.LITERAL_INT, parser.parseUnsignedIntegerLiteralExpression)
	parser.registerPrefixParser(lexer.TRUE, parser.parseBooleanLiteralExpression)
	parser.registerPrefixParser(lexer.FALSE, parser.parseBooleanLiteralExpression)
	parser.registerPrefixParser(lexer.BANG, parser.parsePrefixExpression)
	parser.registerPrefixParser(lexer.BIT_NOT, parser.parsePrefixExpression)
	parser.registerPrefixParser(lexer.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixParser(lexer.LPAR, parser.parseGroupedExpression)

	parser.registerInfixParser(lexer.LPAR, parser.parseCallExpression)
	parser.registerInfixParser(lexer.MINUS, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.PLUS, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.MULTIPLY, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.DIVIDE, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.EQ, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.NEQ, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.LT, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.GT, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.LEQ, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.GEQ, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.LOGICAL_AND, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.LOGICAL_OR, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.BIT_AND, parser.parseInfixExpression)
	parser.registerInfixParser(lexer.BIT_NOT, parser.parseInfixExpression)
}

func New(code *string) *Parser {
	parser := &Parser{
		tokenizer:    lexer.New(code),
		currentToken: nil,
		peekToken:    nil,
		Errors:       []string{},

		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}

	setupParser(parser)

	return parser
}

func NewFromFile(filePath string) (*Parser, error) {
	tokenizer, err := lexer.NewFromFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := &Parser{
		tokenizer:    *tokenizer,
		currentToken: nil,
		peekToken:    nil,
		Errors:       []string{},

		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}

	setupParser(parser)

	return parser, nil
}

// Moves to next token
func (parser *Parser) nextToken() error {
	if parser.peekTokenIs(lexer.EOF) {
		parser.currentToken = parser.peekToken
		return nil
	}

	parser.currentToken = parser.peekToken
	nextToken, err := parser.tokenizer.NextToken()
	if err != nil {
		return err
	}
	parser.peekToken = nextToken
	return nil
}

// Registers prefix parser for PRATT expression parsing
func (parser *Parser) registerPrefixParser(tokenType lexer.TokenType, parseFn prefixParseFn) {
	parser.prefixParseFns[tokenType] = parseFn
}

// Register infix parsers for PRATT expression parsing
func (parser *Parser) registerInfixParser(tokenType lexer.TokenType, parseFn infixParseFn) {
	parser.infixParseFns[tokenType] = parseFn
}

func (parser *Parser) currentTokenIs(tokenType lexer.TokenType) bool {
	return parser.currentToken != nil && tokenType == parser.currentToken.Type
}

func (parser *Parser) peekTokenIs(tokenType lexer.TokenType) bool {
	return parser.peekToken != nil && tokenType == parser.peekToken.Type
}

func (parser *Parser) currentTokenIsDataType() bool {
	return parser.currentToken != nil && (lexer.TYPE_INT == parser.currentToken.Type || lexer.TYPE_UINT == parser.currentToken.Type || lexer.TYPE_BOOL == parser.currentToken.Type)
}

func (parser *Parser) peekTokenIsDataType() bool {
	return parser.peekToken != nil && (lexer.TYPE_INT == parser.peekToken.Type || lexer.TYPE_UINT == parser.peekToken.Type || lexer.TYPE_BOOL == parser.peekToken.Type)
}

func (parser *Parser) currentTokenPrecedence() int {
	if preced, ok := PRECEDENCE_MAP[parser.currentToken.Type]; ok {
		return preced
	}
	return LOWEST
}

func (parser *Parser) peekTokenPrecedence() int {
	if preced, ok := PRECEDENCE_MAP[parser.peekToken.Type]; ok {
		return preced
	}
	return LOWEST
}

func (parser *Parser) reportError(err string) {
	parser.Errors = append(parser.Errors, err)
}

func (parser *Parser) reportUnexpectedToken(found *lexer.Token, expected ...lexer.TokenType) {
	if len(expected) > 0 {
		expectedMsg := ""
		for i, expect := range expected {
			if i != 0 {
				expectedMsg += " or "
			}
			expectedMsg += fmt.Sprint(expect)
		}
		parser.reportError(fmt.Sprintf("Expected %s, found %s %s", expectedMsg, found.Type, found.FormattedLocation()))
	} else {
		parser.reportError(fmt.Sprintf("Unexpected token %s found %s", found.Type, found.FormattedLocation()))
	}
}

// Parses statement as long as there is a token to read
func (parser *Parser) Parse() Program {
	program := newProgram()

	for {
		if parser.currentToken.Type == lexer.EOF {
			break
		}

		statement := parser.parseStatement()
		if statement != nil {
			program.addStatement(statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() Statement {
	var statement Statement = nil
	switch parser.currentToken.Type {
	case lexer.VAR:
		statement = parser.parseDeclarationOrAssignmentOrExpression(false)
	case lexer.IDENTIFIER:
		statement = parser.parseDeclarationOrAssignmentOrExpression(true)
	case lexer.IN:
		statement = parser.parseInputStatement()
	case lexer.IF:
		statement = parser.parserIfStatement()
	case lexer.LOOP:
		statement = parser.parseLoopStatement()
	case lexer.FUN:
		statement = parser.parseFunctionDeclarationStatement()
	case lexer.RETURN:
		statement = parser.parseReturnStatement()
	default:
		statement = parser.parseExpressionStatement()
	}

	return statement
}

func (parser *Parser) parseDeclarationOrAssignmentOrExpression(assignment bool) Statement {
	startToken := parser.currentToken

	if !assignment {
		parser.nextToken()
	}

	name := parser.parseIdentifier()

	var t DataType = INFERED
	if !assignment && parser.peekTokenIs(lexer.COLON) {
		parser.nextToken()

		if !parser.peekToken.IsTypeKeyword() {
			parser.reportUnexpectedToken(parser.peekToken, lexer.TYPES_KEYWORDS...)
		} else {
			parser.nextToken()
			t = DATA_TYPE_MAP[parser.currentToken.Type]
		}
	} else if assignment && parser.peekTokenIs(lexer.LPAR) {
		return parser.parseExpressionStatement()
	}

	if !parser.peekTokenIs(lexer.ASSIGN) {
		parser.reportUnexpectedToken(parser.currentToken, lexer.ASSIGN)
	} else {
		parser.nextToken()
	}

	parser.nextToken()

	value := parser.parseExpression(LOWEST)

	if !parser.peekTokenIs(lexer.SEMICOLON) {
		parser.reportUnexpectedToken(parser.currentToken, lexer.SEMICOLON)
	} else {
		parser.nextToken()
	}

	if assignment {
		return &AssignmentStatement{
			Token: startToken,
			Name:  name,
			Value: value,
		}
	}
	return &DeclarationStatement{
		Token: startToken,
		Name:  name,
		Type:  t,
		Value: value,
	}
}

func (parser *Parser) parseInputStatement() *InputStatement {
	startToken := parser.currentToken
	parser.nextToken()

	identifier := parser.parseIdentifier()
	if identifier == nil {
		return nil
	}

	if !parser.peekTokenIs(lexer.SEMICOLON) {
		parser.reportUnexpectedToken(parser.currentToken, lexer.SEMICOLON)
	} else {
		parser.nextToken()
	}

	return &InputStatement{
		Token: startToken,
		Name:  identifier,
	}
}

func (parser *Parser) parseIdentifier() *Identifier {
	if parser.currentToken.Type != lexer.IDENTIFIER {
		parser.reportUnexpectedToken(parser.currentToken, lexer.IDENTIFIER)
		return nil
	}
	return &Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Value,
	}
}

func (parser *Parser) parseIdentifierExpression() Expression {
	return parser.parseIdentifier()
}

func (parser *Parser) parseUnsignedIntegerLiteral() *UnsignedIntegerLiteralExpression {
	if parser.currentToken.Type != lexer.LITERAL_INT {
		parser.reportUnexpectedToken(parser.currentToken, lexer.LITERAL_INT)
		return nil
	}
	convertedInteger, err := strconv.ParseUint(parser.currentToken.Value, 0, 64)
	if err != nil {
		parser.reportError(fmt.Sprintf("Could not convert `%s` to integer %s", parser.currentToken.Value, parser.currentToken.FormattedLocation()))
		return nil
	}
	return &UnsignedIntegerLiteralExpression{
		Token: parser.currentToken,
		Value: convertedInteger,
	}
}

func (parser *Parser) parseUnsignedIntegerLiteralExpression() Expression {
	return parser.parseUnsignedIntegerLiteral()
}

func (parser *Parser) parseBooleanLiteral() *BooleanLiteralExpression {
	if parser.currentToken.Type != lexer.TRUE && parser.currentToken.Type != lexer.FALSE {
		parser.reportUnexpectedToken(parser.currentToken, lexer.TRUE, lexer.FALSE)
		return nil
	}
	return &BooleanLiteralExpression{
		Token: parser.currentToken,
		Value: parser.currentTokenIs(lexer.TRUE),
	}
}

func (parser *Parser) parseBooleanLiteralExpression() Expression {
	return parser.parseBooleanLiteral()
}

func (parser *Parser) parseExpression(precedence int) Expression {
	prefix := parser.prefixParseFns[parser.currentToken.Type]
	if prefix == nil {
		parser.reportError(fmt.Sprintf("No prefix expression associated to operator %s %s", parser.currentToken.Value, parser.currentToken.FormattedLocation()))
		return nil
	}

	leftExpr := prefix()

	for !parser.peekTokenIs(lexer.SEMICOLON) && precedence < parser.peekTokenPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExpr
		}
		parser.nextToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (parser *Parser) parseGroupedExpression() Expression {
	parser.nextToken()
	expr := parser.parseExpression(LOWEST)
	if !parser.peekTokenIs(lexer.RPAR) {
		parser.reportUnexpectedToken(parser.peekToken, lexer.RPAR)
		return nil
	}
	parser.nextToken()
	return expr
}

func (parser *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Value,
	}
	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Value,
		Left:     left,
	}
	precedence := parser.currentTokenPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)
	return expression
}

func (parser *Parser) parserIfStatement() *IfStatement {
	startToken := parser.currentToken
	parser.nextToken()

	condition, consequence := parser.parseConditionAndConsequence()
	conditions := []Expression{condition}
	consequences := []*StatementsBlock{consequence}

	println(len(consequences))

	var elseConsequence *StatementsBlock = nil
	for parser.peekTokenIs(lexer.ELSE) {
		parser.nextToken()
		if parser.peekTokenIs(lexer.LBRACE) {
			parser.nextToken()
			elseConsequence = parser.parseStatementsBlock()
			break
		} else if parser.peekTokenIs(lexer.IF) {
			parser.nextToken()
			parser.nextToken()
			condition, consequence := parser.parseConditionAndConsequence()
			conditions = append(conditions, condition)
			consequences = append(consequences, consequence)
		} else {
			parser.reportUnexpectedToken(parser.peekToken, lexer.IF, lexer.LBRACE)
			return nil
		}
	}

	return &IfStatement{
		Token:        startToken,
		Conditions:   conditions,
		Consequences: consequences,
		Else:         elseConsequence,
	}
}

func (parser *Parser) parseConditionAndConsequence() (Expression, *StatementsBlock) {
	expression := parser.parseExpression(LOWEST)
	if expression == nil {
		parser.reportError("Could not parse condition expression")
		return nil, nil
	}

	parser.nextToken()

	if !parser.currentTokenIs(lexer.LBRACE) {
		parser.reportUnexpectedToken(parser.currentToken, lexer.LBRACE)
		return nil, nil
	}

	block := parser.parseStatementsBlock()

	return expression, block
}

func (parser *Parser) parseStatementsBlock() *StatementsBlock {
	startToken := parser.currentToken

	parser.nextToken()

	statements := []Statement{}

	for {
		if parser.currentTokenIs(lexer.RBRACE) || parser.currentTokenIs(lexer.EOF) {
			break
		}

		statement := parser.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		parser.nextToken()
	}

	return &StatementsBlock{
		Token:      startToken,
		Statements: statements,
	}
}

func (parser *Parser) parseLoopStatement() *LoopStatement {
	startToken := parser.currentToken
	parser.nextToken()

	condition := parser.parseExpression(LOWEST)
	if condition == nil {
		parser.reportError("Could not parse condition expression")
		return nil
	}

	parser.nextToken()

	if !parser.currentTokenIs(lexer.LBRACE) {
		parser.reportUnexpectedToken(parser.currentToken, lexer.LBRACE)
		return nil
	}

	block := parser.parseStatementsBlock()

	return &LoopStatement{
		Token:     startToken,
		Condition: condition,
		Block:     block,
	}
}

func (parser *Parser) parseFunctionDeclarationStatement() *FunctionDeclarationStatement {
	startToken := parser.currentToken

	parser.nextToken()

	name := parser.parseIdentifier()
	if name == nil {
		return nil
	}

	if !parser.peekTokenIs(lexer.LPAR) {
		parser.reportUnexpectedToken(parser.peekToken, lexer.LPAR)
		return nil
	}

	parser.nextToken()

	identifiers := []*Identifier{}
	dataTypes := []DataType{}
	if !parser.peekTokenIs(lexer.RPAR) {
		parser.nextToken()
		parsedIdentifiers, parsedDataTypes := parser.parseFunctionArgs()
		if parsedIdentifiers != nil {
			identifiers = *parsedIdentifiers
			dataTypes = *parsedDataTypes
		}
	} else {
		parser.nextToken()
	}

	if !parser.peekTokenIs(lexer.COLON) {
		parser.reportUnexpectedToken(parser.peekToken, lexer.LPAR)
	} else {
		parser.nextToken()
	}

	var returnType *DataType
	if !parser.peekTokenIsDataType() {
		parser.reportUnexpectedToken(parser.peekToken, lexer.LPAR)
	} else {
		parser.nextToken()
		dataType, ok := DATA_TYPE_MAP[parser.currentToken.Type]
		if ok {
			returnType = &dataType
		} else {
			parser.reportError("Invalid data type in function return type")
		}
	}

	var body *StatementsBlock = nil
	if !parser.peekTokenIs(lexer.LBRACE) {
		parser.reportUnexpectedToken(parser.peekToken)
	} else {
		parser.nextToken()
		body = parser.parseStatementsBlock()
	}

	return &FunctionDeclarationStatement{
		Token:      startToken,
		Name:       name,
		ArgsNames:  identifiers,
		ArgsTypes:  dataTypes,
		Body:       body,
		ReturnType: returnType,
	}
}

func (parser *Parser) parseFunctionArgs() (*[]*Identifier, *[]DataType) {
	identifiers := []*Identifier{}
	dataTypes := []DataType{}

	for !parser.currentTokenIs(lexer.RPAR) {
		identifier := parser.parseIdentifier()

		parser.nextToken()

		if identifier == nil || !parser.currentTokenIs(lexer.COLON) {
			return nil, nil
		}

		parser.nextToken()

		dataType := DATA_TYPE_MAP[parser.currentToken.Type]

		identifiers = append(identifiers, identifier)
		dataTypes = append(dataTypes, dataType)

		parser.nextToken()

		if parser.currentTokenIs(lexer.COMMA) {
			parser.nextToken()
		}
	}

	return &identifiers, &dataTypes
}

func (parser *Parser) parseExpressionStatement() *ExpressionStatement {
	startToken := parser.currentToken
	expression := parser.parseExpression(LOWEST)
	if expression == nil {
		return nil
	}
	if parser.peekTokenIs(lexer.SEMICOLON) {
		parser.nextToken()
	}
	return &ExpressionStatement{
		Token:      startToken,
		Expression: expression,
	}
}

func (parser *Parser) parseCall(function Expression) *CallExpression {
	expr := &CallExpression{Token: parser.currentToken, Function: function}
	expr.Arguments = parser.parseCallArguments()
	if parser.peekTokenIs(lexer.RPAR) {
		parser.nextToken()
	}
	return expr
}

func (parser *Parser) parseCallExpression(function Expression) Expression {
	return parser.parseCall(function)
}

func (parser *Parser) parseCallArguments() []Expression {
	args := []Expression{}
	if parser.peekTokenIs(lexer.RPAR) {
		parser.nextToken()
		return args
	}
	parser.nextToken()
	args = append(args, parser.parseExpression(LOWEST))
	for parser.peekTokenIs(lexer.COMMA) {
		parser.nextToken()
		parser.nextToken()
		args = append(args, parser.parseExpression(LOWEST))
	}
	if !parser.peekTokenIs(lexer.RPAR) {
		return nil
	}
	return args
}

func (parser *Parser) parseReturnStatement() *ReturnStatement {
	startToken := parser.currentToken

	parser.nextToken()

	expression := parser.parseExpression(LOWEST)
	if expression == nil {
		parser.reportError(fmt.Sprintf("Could not parse expression for return statement %s", startToken.FormattedLocation()))
		return nil
	}

	if !parser.peekTokenIs(lexer.SEMICOLON) {
		parser.reportUnexpectedToken(parser.peekToken, lexer.SEMICOLON)
	} else {
		parser.nextToken()
	}

	return &ReturnStatement{
		Token:      startToken,
		Expression: expression,
	}
}
