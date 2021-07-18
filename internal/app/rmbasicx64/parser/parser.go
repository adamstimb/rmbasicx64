package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/lexer"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // - or NOT
	LOGICAL     // AND OR XOR
	CALL        // MyFunction(X)
)

var precedences = map[string]int{
	token.Equal:              EQUALS,
	token.LessThan:           LESSGREATER,
	token.GreaterThan:        LESSGREATER,
	token.Plus:               SUM,
	token.Minus:              SUM,
	token.Star:               PRODUCT,
	token.ForwardSlash:       PRODUCT,
	token.LeftParen:          CALL,
	token.AND:                LOGICAL,
	token.OR:                 LOGICAL,
	token.XOR:                LOGICAL,
	token.InterestinglyEqual: EQUALS,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[string]prefixParseFn
	infixParseFns  map[string]infixParseFn
	errors         []string // For debugging only - RM Basic-ish error handling to be implemented separately - or not?
	errorMsg       string   // This is for the holding parse error.  We don't collect errors before but fail on the first one.
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.TokenType]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.TokenType]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType string, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType string, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:        l,
		errors:   []string{},
		errorMsg: "",
	}
	// Read 2 tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[string]prefixParseFn)
	p.registerPrefix(token.IdentifierLiteral, p.parseIdentifier)
	p.registerPrefix(token.NumericLiteral, p.parseNumericLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LeftParen, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionDefinition)
	p.registerPrefix(token.StringLiteral, p.parseStringLiteral)
	p.infixParseFns = make(map[string]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.ForwardSlash, p.parseInfixExpression)
	p.registerInfix(token.Star, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.Assign, p.parseInfixExpression)
	p.registerInfix(token.LessThan, p.parseInfixExpression)
	p.registerInfix(token.GreaterThan, p.parseInfixExpression)
	p.registerInfix(token.LeftParen, p.parseCallExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.XOR, p.parseInfixExpression)
	p.registerInfix(token.InterestinglyEqual, p.parseInfixExpression)

	return p
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RightParen) {
		return nil
	}
	return exp
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) Errors() []string {
	// For debugging only - RM Basic-ish error handling to be implemented separately
	return p.errors
}

// GetError returns the current error message and a boolean to indicate if the parser failed
func (p *Parser) GetError() (string, bool) {
	if p.errorMsg != "" {
		// was an error
		return p.errorMsg, true
	}
	return "", false
}

func (p *Parser) peekError(t string) {
	// For debugging only - RM Basic-ish error handling to be implemented separately
	msg := fmt.Sprintf("expected next token to be %s got %s", t, p.peekToken.TokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t string) bool {
	return p.curToken.TokenType == t
}

func (p *Parser) peekTokenIs(t string) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) expectPeek(t string) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// -------------------------------------------------------------------------
// -- If

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.RightCurlyBrace) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LeftParen) {
		return nil
	}
	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RightParen) {
		return nil
	}

	if !p.expectPeek(token.LeftCurlyBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LeftCurlyBrace) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

// -------------------------------------------------------------------------
// -- Function literal
//
// Keywords.78
//
// FUNCTION/ENDFUNC
// 		Define a function
// Defining Syntax
//		FUNCTION name([var1/array1[,var2/array3...]])
//		  :
//		  :
//		  Instruction(s)
//		  :
//		  :
// 		ENDFUN
// Calling Syntax
//		name([exp1[,exp2...]])
// Remarks
//		To declare a function, start with the FUNCTION
//		header, defining the name of the function and its
//		parameters.
//		The action of the function is defined by the
//		instruction(s) on the following line(s), down to the next
//		ENDFUN command.  At least one of these instructions
//		must contain the RESULT command, indicating what
//		value will be returned.
//		To call the function, simple use name in the same way that
//		you use a standard RM Basic function.  The function will
//		be executed and its value will be returns to the expression
//		from which it was called when a RESULT command is executed.
//		You can define function parameters (var1,var2...) in the
//		FUNCTION header.  These are variables or array elements,
//		each one accepting a value from the program to be used in
//		the function.
//		You can specify array names (array1,array2...) in the
//		header.  Using the array reference system (see chapter 5),
//		arrays from the main program can be assigned to the
//		function.
//		Function declarations are usually placed at the end of a
//		program.
// Examples
//		200 REM A function that accepts four numbers and returns
//		their sum
//		210 FUNCTION Add(A, B, C, D)
//		220   Sum := A + B + C + D
//		230   RESULT Sum
//		240 ENDFUN
// Associated Keywords
//		GLOBAL, PROCEDURE/ENDPROC, PROCS, RESULT
//
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RightParen) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()
	ident := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}
	if !p.expectPeek(token.RightParen) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseFunctionBlock() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.curTokenIs(token.ENDFUN) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LeftParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	// parameters can be following be Colon or NewLine
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine)) {
		log.Println("End of instruction expected")
		return nil
	} else {
		p.nextToken()
	}
	lit.Body = p.parseFunctionBlock()
	log.Println("parseFunctionLiteral returning:")
	log.Println(lit.Body.String())
	return lit
}

func (p *Parser) parseFunctionDefinition() ast.Expression {
	lit := &ast.FunctionDefinition{
		Token: p.curToken,
	}
	if !p.peekTokenIs(token.IdentifierLiteral) {
		log.Println("Identifier expected")
	}
	lit.Identifier = &ast.Identifier{
		Token: p.curToken,
		Value: "",
	}
	p.nextToken()
	if !p.expectPeek(token.LeftParen) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine)) {
		log.Println("End of instruction expected")
		return nil
	} else {
		p.nextToken()
	}
	lit.Body = p.parseFunctionBlock()
	return lit
}

func (p *Parser) parseByeStatement() *ast.ByeStatement {
	stmt := &ast.ByeStatement{Token: p.curToken}
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF)) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parseClsStatement() *ast.ClsStatement {
	stmt := &ast.ClsStatement{Token: p.curToken}
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF)) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}
	// Handle PRINT without args
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF)) {
		stmt.Value = &ast.StringLiteral{Value: ""}
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF)) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parseSetModeStatement() *ast.SetModeStatement {
	stmt := &ast.SetModeStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseSetPaperStatement() *ast.SetPaperStatement {
	stmt := &ast.SetPaperStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseSetBorderStatement() *ast.SetBorderStatement {
	stmt := &ast.SetBorderStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseSetPenStatement() *ast.SetPenStatement {
	stmt := &ast.SetPenStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// -------------------------------------------------------------------------
// -- LET

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IdentifierLiteral) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.peekTokenIs(token.Assign) && !p.peekTokenIs(token.Equal) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound)
		return nil
	}
	p.nextToken()

	stmt.BindToken = p.curToken

	p.nextToken()

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.InvalidExpression)
		return nil
	}

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

// -------------------------------------------------------------------------
// -- Bind = or :=

func (p *Parser) parseBindStatement() *ast.BindStatement {
	stmt := &ast.BindStatement{Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}}

	if !p.peekTokenIs(token.Assign) && !p.peekTokenIs(token.Equal) {
		return nil
	}
	p.nextToken()
	stmt.Token = p.curToken
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) {
		p.nextToken()
	}

	return stmt
}

// -------------------------------------------------------------------------
// -- RETURN (TODO: This is Monkey implementation, not Basic)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseResultStatement() *ast.ResultStatement {
	stmt := &ast.ResultStatement{Token: p.curToken}
	p.nextToken()

	stmt.ResultValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) {
		p.nextToken()
	}
	return stmt
}

// -------------------------------------------------------------------------
// -- Numeric Literal

func (p *Parser) parseNumericLiteral() ast.Expression {
	lit := &ast.NumericLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float64", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// -------------------------------------------------------------------------
// -- Expression Statement

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RightParen) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RightParen) {
		return nil
	}
	return args
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) noPrefixParseFnError(t string) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.TokenType]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.TokenType)
		return nil
	}
	leftExp := prefix()

	for !(p.peekTokenIs(token.Semicolon) || p.peekTokenIs(token.NewLine)) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.TokenType]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) {
		p.nextToken()
	}

	return stmt
}

// -------------------------------------------------------------------------
// -- Line

func (p *Parser) ParseLine() *ast.Line {

	statements := []ast.Statement{}

	for !(p.curTokenIs(token.EOF) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.Colon)) {
		statements = append(statements, p.parseStatement())
		p.nextToken()
	}

	return &ast.Line{Statements: statements}
}

// -------------------------------------------------------------------------
// -- Statement

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case token.BYE:
		return p.parseByeStatement()
	case token.CLS:
		return p.parseClsStatement()
	case token.SET:
		p.nextToken()
		switch p.curToken.TokenType {
		case token.MODE:
			return p.parseSetModeStatement()
		case token.PAPER:
			return p.parseSetPaperStatement()
		case token.BORDER:
			return p.parseSetBorderStatement()
		case token.PEN:
			return p.parseSetPenStatement()
		default:
			p.errorMsg = syntaxerror.ErrorMessage((syntaxerror.WrongSetAskAttribute))
			return nil
		}
	case token.PRINT:
		return p.parsePrintStatement()
	case token.LET:
		return p.parseLetStatement() // all these methods need to return when they encounter :
	case token.IdentifierLiteral:
		if p.peekTokenIs(token.Equal) || p.peekTokenIs(token.Assign) {
			return p.parseBindStatement()
		}
	case token.RESULT:
		return p.parseResultStatement()
	default:
		return p.parseExpressionStatement()
	}
	// This should never happen.  It's here because of the mess of supporting optional LET keyword for binding statements
	return p.parseExpressionStatement()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
