package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/ast"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/game"
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
	token.Equal:               EQUALS,
	token.LessThan:            LESSGREATER,
	token.GreaterThan:         LESSGREATER,
	token.LessThanEqualTo1:    LESSGREATER,
	token.GreaterThanEqualTo1: LESSGREATER,
	token.Inequality1:         LESSGREATER,
	token.Plus:                SUM,
	token.Minus:               SUM,
	token.Star:                PRODUCT,
	token.ForwardSlash:        PRODUCT,
	token.LeftParen:           CALL,
	token.AND:                 LOGICAL,
	token.OR:                  LOGICAL,
	token.XOR:                 LOGICAL,
	token.InterestinglyEqual:  EQUALS,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l               *lexer.Lexer
	curToken        token.Token
	peekToken       token.Token
	prefixParseFns  map[string]prefixParseFn
	infixParseFns   map[string]infixParseFn
	errors          []string // For debugging only - RM Basic-ish error handling to be implemented separately - or not?
	errorMsg        string   // This is for the holding parse error.  We don't collect errors before but fail on the first one.
	ErrorTokenIndex int      // the index of the token where an error occured
	g               *game.Game
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

func New(l *lexer.Lexer, g *game.Game) *Parser {
	p := &Parser{
		l:               l,
		g:               g,
		errors:          []string{},
		errorMsg:        "",
		ErrorTokenIndex: -1,
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
	p.registerPrefix(token.StringLiteral, p.parseStringLiteral)
	p.infixParseFns = make(map[string]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.ForwardSlash, p.parseInfixExpression)
	p.registerInfix(token.Star, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.Assign, p.parseInfixExpression)
	p.registerInfix(token.LessThan, p.parseInfixExpression)
	p.registerInfix(token.LessThanEqualTo1, p.parseInfixExpression)
	p.registerInfix(token.LessThanEqualTo2, p.parseInfixExpression)
	p.registerInfix(token.Inequality1, p.parseInfixExpression)
	p.registerInfix(token.Inequality2, p.parseInfixExpression)
	p.registerInfix(token.GreaterThan, p.parseInfixExpression)
	p.registerInfix(token.GreaterThanEqualTo1, p.parseInfixExpression)
	p.registerInfix(token.GreaterThanEqualTo2, p.parseInfixExpression)
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
	retVal := -1.0 // true
	if p.curTokenIs(token.FALSE) {
		retVal = 0
		// Need to hack the literal as well
		p.curToken.Literal = "0"
	} else {
		p.curToken.Literal = "-1.0"
	}
	return &ast.NumericLiteral{
		Token: p.curToken,
		Value: retVal,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// If immediately followed by ( then it's an array or function
	if p.peekTokenIs(token.LeftParen) {
		p.nextToken()
		p.nextToken()
		for {
			val, ok := p.requireExpression()
			if !ok {
				return nil
			} else {
				ident.Subscripts = append(ident.Subscripts, val)
			}
			if p.curTokenIs(token.RightParen) {
				p.nextToken()
				break
			}
			if !p.requireComma() {
				return nil
			}
			if p.onEndOfInstruction() {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
				return nil
			}
		}
	}
	return ident
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

func (p *Parser) endOfInstruction() bool {
	if !(p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF)) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return false
	}
	return true
}

func (p *Parser) requireComma() bool {
	if !p.curTokenIs(token.Comma) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.CommaSeparatorIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		p.nextToken()
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) requireOpenBracket() bool {
	if !p.curTokenIs(token.LeftParen) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.OpeningBracketIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		p.nextToken()
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) requireClosingBracket() bool {
	if !p.curTokenIs(token.RightParen) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		p.nextToken()
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) requireTo() bool {
	if !p.curTokenIs(token.TO) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ToIsNeededBeforeValue)
		p.ErrorTokenIndex = p.curToken.Index
		p.nextToken()
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) requireExpression() (val ast.Expression, ok bool) {
	val = p.parseExpression(LOWEST)
	p.nextToken()
	if _, hasError := p.GetError(); hasError {
		return val, false
	}
	return val, true
}

func (p *Parser) onEndOfInstruction() bool {
	if p.curTokenIs(token.Colon) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.EOF) {
		return true
	}
	return false
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

func (p *Parser) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{
		Token: p.curToken,
	}

	p.nextToken() // consume IF

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		// this will be a syntax error
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ThenExpected)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken() // consume THEN

	stmt.Consequence = p.parseIfConsequence()

	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		stmt.Alternative = p.ParseLine()
	}
	return stmt
}

func (p *Parser) parseUntilStatement() ast.Statement {
	stmt := &ast.UntilStatement{
		Token: p.curToken,
	}

	p.nextToken() // consume UNTIL

	// Ensure condition follows
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	stmt.Condition = p.parseExpression(LOWEST)

	// Ensure end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseRemStatement() *ast.RemStatement {
	remToken := p.curToken
	p.nextToken()
	return &ast.RemStatement{Token: remToken, Comment: p.curToken}
}

func (p *Parser) parseByeStatement() *ast.ByeStatement {
	stmt := &ast.ByeStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseEndStatement() *ast.EndStatement {
	stmt := &ast.EndStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseListStatement() *ast.ListStatement {
	stmt := &ast.ListStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseRunStatement() *ast.RunStatement {
	stmt := &ast.RunStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseNewStatement() *ast.NewStatement {
	stmt := &ast.NewStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseClearblockStatement() *ast.ClearblockStatement {
	stmt := &ast.ClearblockStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseRenumberStatement() *ast.RenumberStatement {
	stmt := &ast.RenumberStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseClsStatement() *ast.ClsStatement {
	stmt := &ast.ClsStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseHomeStatement() *ast.HomeStatement {
	stmt := &ast.HomeStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseDirStatement() *ast.DirStatement {
	stmt := &ast.DirStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}
	stmt.PrintList = make([]interface{}, 0)
	// Handle PRINT without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		stmt.PrintList = append(stmt.PrintList, &ast.StringLiteral{Value: ""})
		return stmt
	}
	p.nextToken()
	for !(p.curTokenIs(token.Colon) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.EOF)) {
		// ; -> noSpace
		if p.curTokenIs(token.Semicolon) {
			stmt.PrintList = append(stmt.PrintList, "noSpace")
			p.nextToken()
			continue
		}
		// , -> nextPrintZone
		if p.curTokenIs(token.Comma) {
			stmt.PrintList = append(stmt.PrintList, "nextPrintZone")
			p.nextToken()
			continue
		}
		// ! -> newLine
		if p.curTokenIs(token.Exclamation) {
			stmt.PrintList = append(stmt.PrintList, "newLine")
			p.nextToken()
			continue
		}
		val := p.parseExpression(LOWEST)
		stmt.PrintList = append(stmt.PrintList, val)
		p.nextToken()
		if p.curTokenIs(token.Colon) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.EOF) {
			break
		}
	}
	return stmt
}

func (p *Parser) parsePlotStatement() *ast.PlotStatement {
	stmt := &ast.PlotStatement{Token: p.curToken}
	// Handle PLOT without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericOrStringExpressionNeeded)
		return nil
	}
	// Get Value
	val, ok := p.requireExpression()
	if ok {
		stmt.Value = val
	} else {
		return nil
	}
	// ,
	if !p.requireComma() {
		return nil
	}
	// Get coordinate list
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.DIRECTION:
			p.nextToken()
			stmt.Direction = p.parseExpression(LOWEST)
		case token.FONT:
			p.nextToken()
			stmt.Font = p.parseExpression(LOWEST)
		case token.OVER:
			p.nextToken()
			stmt.Over = p.parseExpression(LOWEST)
		case token.SIZE:
			p.nextToken()
			stmt.SizeX = p.parseExpression(LOWEST)
			if p.peekTokenIs(token.Comma) {
				p.nextToken()
				p.nextToken()
				stmt.SizeY = p.parseExpression(LOWEST)
			} else {
				stmt.SizeY = stmt.SizeX
			}
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLineStatement() *ast.LineStatement {
	stmt := &ast.LineStatement{Token: p.curToken}
	// Handle LINE without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	// Get coordinate list
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.OVER:
			p.nextToken()
			stmt.Over = p.parseExpression(LOWEST)
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseCircleStatement() *ast.CircleStatement {
	stmt := &ast.CircleStatement{Token: p.curToken}
	// Handle CIRCLE without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	// Get radius
	val, ok := p.requireExpression()
	if !ok {
		return nil
	}
	stmt.Radius = val
	// ,
	if !p.requireComma() {
		return nil
	}
	// Get coordinate list
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() && !p.curTokenIs(token.ELSE) {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.OVER:
			p.nextToken()
			stmt.Over = p.parseExpression(LOWEST)
		case token.STYLE:
			p.nextToken()
			// get required fill style
			if val, ok := p.requireExpression(); ok {
				stmt.FillStyle = val
			} else {
				return nil
			}
			if p.curTokenIs(token.Comma) {
				// get required fill hatching
				p.nextToken()
				if val, ok := p.requireExpression(); ok {
					stmt.FillHatching = val
				} else {
					return nil
				}
				if p.curTokenIs(token.Comma) {
					// get required fill colour2
					p.nextToken()
					if val, ok := p.requireExpression(); ok {
						stmt.FillColour2 = val
					} else {
						return nil
					}
				}
			}
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parsePointsStatement() *ast.PointsStatement {
	stmt := &ast.PointsStatement{Token: p.curToken}
	// Handle POINTS without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	// Get coordinate list
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() && !p.curTokenIs(token.ELSE) {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.STYLE:
			p.nextToken()
			stmt.Style = p.parseExpression(LOWEST)
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.OVER:
			p.nextToken()
			stmt.Over = p.parseExpression(LOWEST)
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseFloodStatement() *ast.FloodStatement {
	stmt := &ast.FloodStatement{Token: p.curToken}
	// Handle FLOOD without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	// Get coordinate list
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() && !p.curTokenIs(token.ELSE) {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.EDGE:
			p.nextToken()
			if val, ok := p.requireExpression(); ok {
				stmt.UseEdgeColour = val
			} else {
				return nil
			}
			if !p.requireComma() {
				return nil
			}
			if val, ok := p.requireExpression(); ok {
				stmt.EdgeColour = val
			} else {
				return nil
			}
		case token.STYLE:
			p.nextToken()
			// get required fill style
			if val, ok := p.requireExpression(); ok {
				stmt.FillStyle = val
			} else {
				return nil
			}
			if p.curTokenIs(token.Comma) {
				// get required fill hatching
				p.nextToken()
				if val, ok := p.requireExpression(); ok {
					stmt.FillHatching = val
				} else {
					return nil
				}
				if p.curTokenIs(token.Comma) {
					// get required fill colour2
					p.nextToken()
					if val, ok := p.requireExpression(); ok {
						stmt.FillColour2 = val
					} else {
						return nil
					}
				}
			}
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseAreaStatement() *ast.AreaStatement {
	stmt := &ast.AreaStatement{Token: p.curToken}
	// Handle LINE without args
	p.nextToken()
	if p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return nil
	}
	// Get coordinate list
	for !p.onEndOfInstruction() {
		// Get X
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// ,
		if !p.requireComma() {
			return nil
		}
		// Get y
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.CoordList = append(stmt.CoordList, val)
		// Require ; if more coordinates to follow
		if p.curTokenIs(token.Comma) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.SemicolonSeparatorIsNeeded)
			return nil
		}
		// Break loop if no more coordinates to follow
		if !p.curTokenIs(token.Semicolon) {
			break
		}
		p.nextToken() // consume ;
	}
	// Handle no options list
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle options list
	for !p.onEndOfInstruction() {
		tokenType := p.curToken.TokenType
		switch tokenType {
		case token.BRUSH:
			p.nextToken()
			stmt.Brush = p.parseExpression(LOWEST)
		case token.OVER:
			p.nextToken()
			stmt.Over = p.parseExpression(LOWEST)
		case token.STYLE:
			p.nextToken()
			// get required fill style
			if val, ok := p.requireExpression(); ok {
				stmt.FillStyle = val
			} else {
				return nil
			}
			if p.curTokenIs(token.Comma) {
				// get required fill hatching
				p.nextToken()
				if val, ok := p.requireExpression(); ok {
					stmt.FillHatching = val
				} else {
					return nil
				}
				if p.curTokenIs(token.Comma) {
					// get required fill colour2
					p.nextToken()
					if val, ok := p.requireExpression(); ok {
						stmt.FillColour2 = val
					} else {
						return nil
					}
				}
			}
		default:
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownSetAskAttribute)
			return nil
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseDataStatement() *ast.DataStatement {
	stmt := &ast.DataStatement{Token: p.curToken}
	// Handle DATA without args
	p.nextToken()
	if p.onEndOfInstruction() {
		return nil
	}
	// Get item list
	for !p.onEndOfInstruction() {
		// Commas without a preceding value are regarded as representing an item with 0 value
		if p.curTokenIs(token.Comma) {
			stmt.ItemList = append(stmt.ItemList, token.Token{TokenType: token.NumericLiteral, Literal: "0", Index: p.curToken.Index})
			p.nextToken()
			continue
		}
		// Token must be string literal, numeric literal or an identifier literal without $ or % postfix
		// (this is interpreted as a string literal)
		if (p.curTokenIs(token.NumericLiteral) || p.curTokenIs(token.StringLiteral)) || (p.curTokenIs(token.IdentifierLiteral) && (p.curToken.Literal[len(p.curToken.Literal)-1] != '%' && p.curToken.Literal[len(p.curToken.Literal)-1] != '$')) {
			stmt.ItemList = append(stmt.ItemList, p.curToken)
			// each item must be followed by either end of instruction or comma
			p.nextToken()
			if p.onEndOfInstruction() {
				return stmt
			}
			if !p.requireComma() {
				return nil
			}
		} else {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnableToReadExcessData)
			return nil
		}
	}
	return stmt
}

func (p *Parser) parseReadStatement() *ast.ReadStatement {
	stmt := &ast.ReadStatement{Token: p.curToken}
	p.nextToken() // consume READ
	// TODO: Accept and parse arrays
	// Require variable name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	// Collect variable list
	for !p.onEndOfInstruction() {
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		// parse array subscripts
		if subscripts, ok := p.getArraySubscripts(); ok {
			ident.Subscripts = subscripts
		} else {
			return nil
		}
		stmt.VariableList = append(stmt.VariableList, ident)
		if !p.onEndOfInstruction() {
			// require comma
			if !p.requireComma() {
				return nil
			}
		}
	}
	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}
	p.nextToken() // consume FOR
	// Require variable name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// Require assignment
	if !p.peekTokenIs(token.Assign) && !p.peekTokenIs(token.Equal) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	stmt.BindToken = p.curToken
	p.nextToken()
	// Require start value
	val, ok := p.requireExpression()
	if !ok {
		return nil
	}
	stmt.Start = val
	// Require TO token
	if !p.curTokenIs(token.TO) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ToIsNeededBeforeValue)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Require stop value
	val, ok = p.requireExpression()
	if !ok {
		return nil
	}
	stmt.Stop = val
	// Detect optional STEP
	stmt.Step = nil
	if p.curTokenIs(token.STEP) {
		p.nextToken()
		// Require step value
		val, ok = p.requireExpression()
		if !ok {
			return nil
		}
		stmt.Step = val
	}
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseNextStatement() *ast.NextStatement {
	stmt := &ast.NextStatement{Token: p.curToken}
	p.nextToken() // consume NEXT
	// Require variable name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSubroutineStatement() *ast.SubroutineStatement {
	stmt := &ast.SubroutineStatement{Token: p.curToken}
	p.nextToken() // consume SUBROUTINE
	// Require name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NameOfDefinitionRequired)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseGosubStatement() *ast.GosubStatement {
	stmt := &ast.GosubStatement{Token: p.curToken}
	p.nextToken() // consume SUBROUTINE
	// Require name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NameOfDefinitionRequired)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken() // consume RETURN
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{Token: p.curToken}
	p.nextToken() // consume FUNCTION
	// Require name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NameOfDefinitionRequired)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	if !p.requireOpenBracket() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.OpeningBracketIsNeeded)
		return nil
	}
	for {
		if p.curTokenIs(token.RightParen) {
			p.nextToken()
			break
		}
		// Require variable name
		if !p.curTokenIs(token.IdentifierLiteral) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
			return nil
		} else {
			stmt.ReceiveArgs = append(stmt.ReceiveArgs, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			p.nextToken()
		}
		if p.curTokenIs(token.RightParen) {
			p.nextToken()
			break
		}
		if !p.requireComma() {
			return nil
		}
		if p.onEndOfInstruction() {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
			return nil
		}
	}
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseEndfunStatement() *ast.EndfunStatement {
	stmt := &ast.EndfunStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseDimStatement() *ast.DimStatement {
	stmt := &ast.DimStatement{Token: p.curToken}
	p.nextToken() // consume DIM
	// Require variable name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	// Require at least one subscript
	p.nextToken()
	if !p.requireOpenBracket() {
		return nil
	}
	for {
		val, ok := p.requireExpression()
		if !ok {
			return nil
		}
		stmt.Subscripts = append(stmt.Subscripts, val)
		if p.curTokenIs(token.RightParen) {
			p.nextToken()
			break
		}
		if !p.requireComma() {
			return nil
		}
		if p.onEndOfInstruction() {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
			return nil
		}
	}
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseProcedureDeclaration() *ast.ProcedureDeclaration {
	stmt := &ast.ProcedureDeclaration{Token: p.curToken}
	p.nextToken() // consume PROCEDURE
	// Require name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NameOfDefinitionRequired)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	log.Printf("got name")
	// Optional end of instruction
	if p.onEndOfInstruction() {
		return stmt
	}
	// Optional receive args
	log.Printf("get receive args")
	for !p.onEndOfInstruction() && !p.curTokenIs(token.RETURN) {
		// Require variable name
		if !p.curTokenIs(token.IdentifierLiteral) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
			return nil
		} else {
			log.Printf("got receive arg")
			stmt.ReceiveArgs = append(stmt.ReceiveArgs, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			p.nextToken()
		}
		// Require comma, end of instruction or return
		if p.curTokenIs(token.Comma) {
			p.nextToken()
			continue
		}
		if p.onEndOfInstruction() || p.curTokenIs(token.RETURN) {
			break
		} else {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
			return nil
		}
	}
	if p.onEndOfInstruction() {
		return stmt
	}
	// Optional RETURN token following by required args
	log.Printf("get return args")
	if p.curTokenIs(token.RETURN) {
		log.Printf("got return")
		p.nextToken()
		for !p.onEndOfInstruction() {
			// Require variable name
			if !p.curTokenIs(token.IdentifierLiteral) {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
				return nil
			} else {
				stmt.ReturnArgs = append(stmt.ReturnArgs, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
				p.nextToken()
			}
			// Require comma or end of instruction
			if p.curTokenIs(token.Comma) {
				p.nextToken()
				continue
			}
			if p.onEndOfInstruction() {
				break
			} else {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
				return nil
			}
		}
	} else {
		log.Printf("end of instruction expected")
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	log.Printf("return stmt")
	return stmt
}

func (p *Parser) parseProcedureCallStatement() *ast.ProcedureCallStatement {
	stmt := &ast.ProcedureCallStatement{Token: p.curToken}
	// Require name
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NameOfDefinitionRequired)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Case of no arguments
	if p.onEndOfInstruction() {
		return stmt
	}
	// Optional receive args
	log.Printf("get receive args")
	for !p.onEndOfInstruction() && !p.curTokenIs(token.RECEIVE) {
		// Require variable name
		if !p.curTokenIs(token.IdentifierLiteral) {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
			return nil
		} else {
			log.Printf("got receive arg")
			stmt.ReceiveArgs = append(stmt.ReceiveArgs, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			p.nextToken()
		}
		// Require comma, end of instruction or return
		if p.curTokenIs(token.Comma) {
			p.nextToken()
			continue
		}
		if p.onEndOfInstruction() || p.curTokenIs(token.RECEIVE) {
			break
		} else {
			p.ErrorTokenIndex = p.curToken.Index
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
			return nil
		}
	}
	if p.onEndOfInstruction() {
		return stmt
	}
	// Optional RETURN token following by required args
	log.Printf("get return args")
	if p.curTokenIs(token.RECEIVE) {
		log.Printf("got receive")
		p.nextToken()
		for !p.onEndOfInstruction() {
			// Require variable name
			if !p.curTokenIs(token.IdentifierLiteral) {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
				return nil
			} else {
				stmt.ReceiveArgs = append(stmt.ReceiveArgs, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
				p.nextToken()
			}
			// Require comma or end of instruction
			if p.curTokenIs(token.Comma) {
				p.nextToken()
				continue
			}
			if p.onEndOfInstruction() {
				break
			} else {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
				return nil
			}
		}
	} else {
		log.Printf("end of instruction expected")
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return nil
}

func (p *Parser) parseAskMouseStatement() *ast.AskMouseStatement {
	stmt := &ast.AskMouseStatement{Token: p.curToken}
	p.nextToken() // consume MOUSE
	// Handle no arguments
	if p.onEndOfInstruction() {
		return stmt
	}
	// Otherwise require e1, e2
	// e1
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.XName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// ,
	if !p.requireComma() {
		return nil
	}
	// e2
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.YName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Handle no button argument
	if p.onEndOfInstruction() {
		return stmt
	}
	// Handle optional button argument
	if !p.requireComma() {
		return nil
	}
	// e3
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.BName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseAskBlocksizeStatement() *ast.AskBlocksizeStatement {
	stmt := &ast.AskBlocksizeStatement{Token: p.curToken}
	p.nextToken() // consume BLOCKSIZE
	// Handle no arguments
	if p.onEndOfInstruction() {
		return stmt
	}
	// Require block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Require width variable
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Width = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Handle no more args
	if p.onEndOfInstruction() {
		return stmt
	}
	// Require ,height variable
	if !p.requireComma() {
		return nil
	}
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Height = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Handle no more args
	if p.onEndOfInstruction() {
		return stmt
	}
	// Require , mode variable
	if !p.requireComma() {
		return nil
	}
	if !p.curTokenIs(token.IdentifierLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		return nil
	}
	stmt.Mode = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	// Require end of instruction
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSaveStatement() *ast.SaveStatement {
	stmt := &ast.SaveStatement{Token: p.curToken}
	// Handle SAVE without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseLoadStatement() *ast.LoadStatement {
	stmt := &ast.LoadStatement{Token: p.curToken}
	// Handle LOAD without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseFetchStatement() *ast.FetchStatement {
	stmt := &ast.FetchStatement{Token: p.curToken}
	// Handle FETCH without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Get block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get path
	if val, ok := p.requireExpression(); ok {
		stmt.Path = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseWriteblockStatement() *ast.WriteblockStatement {
	stmt := &ast.WriteblockStatement{Token: p.curToken}
	// Handle WRITEBLOCK without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Get block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get X,
	if val, ok := p.requireExpression(); ok {
		stmt.X = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y,
	if val, ok := p.requireExpression(); ok {
		stmt.Y = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Over,
	if val, ok := p.requireExpression(); ok {
		stmt.Over = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseReadblockStatement() *ast.ReadblockStatement {
	stmt := &ast.ReadblockStatement{Token: p.curToken}
	// Handle READBLOCK without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Get block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get X1,
	if val, ok := p.requireExpression(); ok {
		stmt.X1 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y1,
	if val, ok := p.requireExpression(); ok {
		stmt.Y1 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get X2,
	if val, ok := p.requireExpression(); ok {
		stmt.X2 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y2,
	if val, ok := p.requireExpression(); ok {
		stmt.Y2 = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseCopyblockStatement() *ast.CopyblockStatement {
	stmt := &ast.CopyblockStatement{Token: p.curToken}
	// Handle READBLOCK without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Get X1,
	if val, ok := p.requireExpression(); ok {
		stmt.X1 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y1,
	if val, ok := p.requireExpression(); ok {
		stmt.Y1 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get X2,
	if val, ok := p.requireExpression(); ok {
		stmt.X2 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y2,
	if val, ok := p.requireExpression(); ok {
		stmt.Y2 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Dx,
	if val, ok := p.requireExpression(); ok {
		stmt.Dx = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Dy,
	if val, ok := p.requireExpression(); ok {
		stmt.Dy = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Over
	if val, ok := p.requireExpression(); ok {
		stmt.Over = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSquashStatement() *ast.SquashStatement {
	stmt := &ast.SquashStatement{Token: p.curToken}
	// Handle SQUASH without args
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	// Get block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get X,
	if val, ok := p.requireExpression(); ok {
		stmt.X = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Y,
	if val, ok := p.requireExpression(); ok {
		stmt.Y = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get Over,
	if val, ok := p.requireExpression(); ok {
		stmt.Over = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetModeStatement() *ast.SetModeStatement {
	stmt := &ast.SetModeStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetPaperStatement() *ast.SetPaperStatement {
	stmt := &ast.SetPaperStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetBorderStatement() *ast.SetBorderStatement {
	stmt := &ast.SetBorderStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetPenStatement() *ast.SetPenStatement {
	stmt := &ast.SetPenStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetMouseStatement() *ast.SetMouseStatement {
	stmt := &ast.SetMouseStatement{Token: p.curToken}
	p.nextToken()
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetCurposStatement() *ast.SetCurposStatement {
	stmt := &ast.SetCurposStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	// Get col
	stmt.Col = p.parseExpression(LOWEST)
	p.nextToken()
	// Must have ,
	if !p.curTokenIs(token.Comma) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.CommaSeparatorIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken() // consume ,
	if p.curTokenIs(token.Colon) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	// Get row
	stmt.Row = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetPatternStatement() *ast.SetPatternStatement {
	stmt := &ast.SetPatternStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	// Get Slot, Row
	if val, ok := p.requireExpression(); ok {
		stmt.Slot = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	if val, ok := p.requireExpression(); ok {
		stmt.Row = val
	} else {
		return nil
	}
	// Get TO
	if !p.requireTo() {
		return nil
	}
	// Get c1,
	if val, ok := p.requireExpression(); ok {
		stmt.C1 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get c2,
	if val, ok := p.requireExpression(); ok {
		stmt.C2 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get c3,
	if val, ok := p.requireExpression(); ok {
		stmt.C3 = val
	} else {
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get c4
	if val, ok := p.requireExpression(); ok {
		stmt.C4 = val
	} else {
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseMoveStatement() *ast.MoveStatement {
	stmt := &ast.MoveStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	// Get cols
	stmt.Cols = p.parseExpression(LOWEST)
	p.nextToken()
	// Must have ,
	if !p.curTokenIs(token.Comma) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.CommaSeparatorIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	p.nextToken() // consume ,
	if p.curTokenIs(token.Colon) || p.curTokenIs(token.NewLine) || p.curTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	// Get rows
	stmt.Rows = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseDelblockStatement() *ast.DelblockStatement {
	stmt := &ast.DelblockStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	// Get block
	stmt.Block = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseKeepStatement() *ast.KeepStatement {
	stmt := &ast.KeepStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	// Get block,
	if val, ok := p.requireExpression(); ok {
		stmt.Block = val
	} else {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	if !p.requireComma() {
		return nil
	}
	// Get path
	if val, ok := p.requireExpression(); ok {
		stmt.Path = val
	} else {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.StringExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetColourStatement() *ast.SetColourStatement {
	stmt := &ast.SetColourStatement{Token: p.curToken}
	p.nextToken()
	// Get required e1
	val, ok := p.requireExpression()
	if !ok {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	} else {
		stmt.PaletteSlot = val
	}
	// Get required TO
	if !p.requireTo() {
		return nil
	}
	// Get required e2
	val, ok = p.requireExpression()
	if !ok {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	} else {
		stmt.BasicColour = val
	}
	// Return if no more args
	if p.onEndOfInstruction() {
		stmt.FlashSpeed = nil
		stmt.FlashColour = nil
		return stmt
	}
	// Otherwise get required comma
	if !p.requireComma() {
		return nil
	}
	// Get required e3
	val, ok = p.requireExpression()
	if !ok {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	} else {
		stmt.FlashSpeed = val
	}
	// Get required comma
	if !p.requireComma() {
		return nil
	}
	// Get required e4
	val, ok = p.requireExpression()
	if !ok {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	} else {
		stmt.FlashColour = val
	}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetDegStatement() *ast.SetDegStatement {
	stmt := &ast.SetDegStatement{Token: p.curToken}
	p.nextToken()
	if p.onEndOfInstruction() {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetConfigBootStatement() *ast.SetConfigBootStatement {
	stmt := &ast.SetConfigBootStatement{Token: p.curToken}
	p.nextToken()
	if p.onEndOfInstruction() {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseSetFillStyleStatement() *ast.SetFillStyleStatement {
	stmt := &ast.SetFillStyleStatement{Token: p.curToken}
	p.nextToken()
	if p.onEndOfInstruction() {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}
	// get required fill style
	if val, ok := p.requireExpression(); ok {
		stmt.FillStyle = val
	} else {
		return nil
	}
	if p.curTokenIs(token.Comma) {
		// get required fill hatching
		p.nextToken()
		if val, ok := p.requireExpression(); ok {
			stmt.FillHatching = val
		} else {
			return nil
		}
		if p.curTokenIs(token.Comma) {
			// get required fill colour2
			p.nextToken()
			if val, ok := p.requireExpression(); ok {
				stmt.FillColour2 = val
			} else {
				return nil
			}
		}
	}
	return stmt
}

func (p *Parser) parseSetRadStatement() *ast.SetRadStatement {
	stmt := &ast.SetRadStatement{Token: p.curToken}
	if p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.EOF) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseGotoStatement() *ast.GotoStatement {
	stmt := &ast.GotoStatement{Token: p.curToken}
	p.nextToken()
	if !p.curTokenIs(token.NumericLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.LineNumberLabelNeeded)
		return nil
	}
	stmt.Linenumber = p.curToken
	p.nextToken()
	if !p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parseEditStatement() *ast.EditStatement {
	stmt := &ast.EditStatement{Token: p.curToken}
	p.nextToken()
	// TODO: No line number passed so try to get line number of last error
	if p.onEndOfInstruction() {
		return stmt
	}
	// Line number
	if !p.curTokenIs(token.NumericLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.LineNumberLabelNeeded)
		return nil
	}
	stmt.Linenumber = p.curToken
	p.nextToken()
	if !p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parseRestoreStatement() *ast.RestoreStatement {
	stmt := &ast.RestoreStatement{Token: p.curToken}
	p.nextToken()
	if p.onEndOfInstruction() {
		return stmt
	}
	// Line number
	if !p.curTokenIs(token.NumericLiteral) {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.LineNumberLabelNeeded)
		return nil
	}
	stmt.Linenumber = p.curToken
	p.nextToken()
	if !p.onEndOfInstruction() {
		p.ErrorTokenIndex = p.curToken.Index
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		return nil
	}
	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	stmt := &ast.RepeatStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

// -------------------------------------------------------------------------
// -- LET

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IdentifierLiteral) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.VariableNameIsNeeded)
		p.ErrorTokenIndex = p.curToken.Index + 1
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(token.LeftParen) {
		p.nextToken()
		p.nextToken()
		for {
			val, ok := p.requireExpression()
			if !ok {
				return nil
			} else {
				stmt.Name.Subscripts = append(stmt.Name.Subscripts, val)
			}
			if p.curTokenIs(token.RightParen) {
				p.nextToken()
				break
			}
			if !p.requireComma() {
				return nil
			}
			if p.onEndOfInstruction() {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
				return nil
			}
		}
	} else {
		p.nextToken()
	}
	if !p.curTokenIs(token.Assign) && !p.curTokenIs(token.Equal) {
		p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound)
		p.ErrorTokenIndex = p.curToken.Index
		return nil
	}

	stmt.BindToken = p.curToken

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.endOfInstruction() {
		return stmt
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

func (p *Parser) getArraySubscripts() (subscripts []ast.Expression, ok bool) {
	if p.peekTokenIs(token.LeftParen) {
		p.nextToken()
		p.nextToken()
		for {
			val, ok := p.requireExpression()
			if !ok {
				return subscripts, false
			} else {
				subscripts = append(subscripts, val)
			}
			if p.curTokenIs(token.RightParen) {
				p.nextToken()
				break
			}
			if !p.requireComma() {
				return subscripts, false
			}
			if p.onEndOfInstruction() {
				p.ErrorTokenIndex = p.curToken.Index
				p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.ClosingBracketIsNeeded)
				return subscripts, false
			}
		}
	}
	return subscripts, true
}

func (p *Parser) parseBindArrayStatement() *ast.BindStatement {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt := &ast.BindStatement{Name: ident}
	// parse subscripts
	if subscripts, ok := p.getArraySubscripts(); ok {
		stmt.Name.Subscripts = subscripts
	} else {
		return nil
	}
	// then get assign token and value expression
	if !p.curTokenIs(token.Assign) && !p.curTokenIs(token.Equal) {
		return nil
	}
	stmt.Token = p.curToken
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

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

func (p *Parser) parseEndprocStatement() *ast.EndprocStatement {
	stmt := &ast.EndprocStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
}

func (p *Parser) parseLeaveStatement() *ast.LeaveStatement {
	stmt := &ast.LeaveStatement{Token: p.curToken}
	if p.endOfInstruction() {
		return stmt
	}
	return nil
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
	msg := fmt.Sprintf("no prefix parse function for %s found, current token=%s", t, p.curToken.Literal)
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

func (p *Parser) JumpToToken(i int) {
	p.l.JumpToToken(i)
	p.nextToken()
	p.nextToken()
}

// -------------------------------------------------------------------------
// -- Line

func (p *Parser) PrettyPrint() string {
	indent := "  "
	lineString := p.g.PrettyPrintIndent
	// Add indent for following lines
	if p.curTokenIs(token.REPEAT) || p.curTokenIs(token.FOR) || p.curTokenIs(token.FUNCTION) || p.curTokenIs(token.PROCEDURE) {
		p.g.PrettyPrintIndent += indent
	}
	// Remove indent for this and following lines
	if p.curTokenIs(token.UNTIL) || p.curTokenIs(token.NEXT) || p.curTokenIs(token.ENDFUN) || p.curTokenIs(token.ENDPROC) {
		if len(p.g.PrettyPrintIndent) >= 2 {
			p.g.PrettyPrintIndent = p.g.PrettyPrintIndent[:len(p.g.PrettyPrintIndent)-2]
			lineString = p.g.PrettyPrintIndent
		}
	}
	// Build string
	for {
		if p.ErrorTokenIndex > 0 {
			if p.curToken.Index == p.ErrorTokenIndex {
				lineString += ">> "
			}
		}
		// Break if end of line
		if p.curTokenIs(token.EOF) || p.curTokenIs(token.NewLine) {
			break
		}
		// Put string literals in double-quotes
		if p.curToken.TokenType == token.StringLiteral {
			lineString += fmt.Sprintf("%q", p.curToken.Literal) + " "
			p.nextToken()
			continue
		}
		// Never a space after (
		if p.curToken.TokenType == token.LeftParen {
			lineString += "("
			p.nextToken()
			continue
		}
		// Remove trailing space if ) or (
		if len(lineString) > 0 {
			if lineString[len(lineString)-1] == ' ' && (p.curToken.TokenType == token.RightParen || p.curToken.TokenType == token.LeftParen) {
				lineString = lineString[0 : len(lineString)-1]
				lineString += ") "
				p.nextToken()
				continue
			}
		}
		// Remove trailing space if ,
		if len(lineString) > 0 {
			if lineString[len(lineString)-1] == ' ' && p.curToken.TokenType == token.Comma {
				lineString = lineString[0 : len(lineString)-1]
				lineString += ", "
				p.nextToken()
				continue
			}
		}
		// Remove trailing space if ;
		if len(lineString) > 0 {
			if lineString[len(lineString)-1] == ' ' && p.curToken.TokenType == token.Comma {
				lineString = lineString[0 : len(lineString)-1]
				lineString += "; "
				p.nextToken()
				continue
			}
		}
		// Otherwise add literal with trailing space
		curLiteral := p.curToken.Literal
		_, ok := lexer.Builtins[curLiteral]
		if ok {
			lineString += p.curToken.Literal
		} else {
			lineString += p.curToken.Literal + " "
		}
		p.nextToken()
	}
	return strings.TrimRight(lineString, " ")
}

func (p *Parser) ParseLine() *ast.Line {

	statements := []ast.Statement{}

	// Catch new line for stored program
	if p.curTokenIs(token.NumericLiteral) {
		// Get line number
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		lineNumber := int(val)
		p.nextToken()
		// Extract and pretty print the line string
		lineString := p.PrettyPrint()
		return &ast.Line{Statements: nil, LineNumber: lineNumber, LineString: lineString}
	}

	for !(p.curTokenIs(token.EOF) || p.curTokenIs(token.NewLine)) {
		statements = append(statements, p.parseStatement())
		p.nextToken()
		if p.curTokenIs(token.Colon) {
			p.nextToken()
		}
	}

	return &ast.Line{Statements: statements}
}

func (p *Parser) parseIfConsequence() *ast.Line {

	statements := []ast.Statement{}

	for !(p.curTokenIs(token.EOF) || p.curTokenIs(token.NewLine)) {
		statements = append(statements, p.parseStatement())
		if p.curTokenIs(token.ELSE) {
			break
		}
		p.nextToken()
		if p.curTokenIs(token.Colon) {
			p.nextToken()
		}
	}

	return &ast.Line{Statements: statements}
}

// -------------------------------------------------------------------------
// -- Statement

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case token.ELSE:
		return nil
	case token.REM:
		return p.parseRemStatement()
	case token.BYE:
		return p.parseByeStatement()
	case token.END:
		return p.parseEndStatement()
	case token.LIST:
		return p.parseListStatement()
	case token.RUN:
		return p.parseRunStatement()
	case token.NEW:
		return p.parseNewStatement()
	case token.CLS:
		return p.parseClsStatement()
	case token.HOME:
		return p.parseHomeStatement()
	case token.DIR:
		return p.parseDirStatement()
	case token.SAVE:
		return p.parseSaveStatement()
	case token.LOAD:
		return p.parseLoadStatement()
	case token.GOTO:
		return p.parseGotoStatement()
	case token.EDIT:
		return p.parseEditStatement()
	case token.RENUMBER:
		return p.parseRenumberStatement()
	case token.REPEAT:
		return p.parseRepeatStatement()
	case token.UNTIL:
		return p.parseUntilStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.NEXT:
		return p.parseNextStatement()
	case token.SUBROUTINE:
		return p.parseSubroutineStatement()
	case token.GOSUB:
		return p.parseGosubStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionDeclaration()
	case token.ENDFUN:
		return p.parseEndfunStatement()
	case token.PROCEDURE:
		return p.parseProcedureDeclaration()
	case token.ENDPROC:
		return p.parseEndprocStatement()
	case token.LEAVE:
		return p.parseLeaveStatement()
	case token.DIM:
		return p.parseDimStatement()
	case token.ASK:
		p.nextToken()
		switch p.curToken.TokenType {
		case token.MOUSE:
			return p.parseAskMouseStatement()
		case token.BLOCKSIZE:
			return p.parseAskBlocksizeStatement()
		}
	case token.SET:
		p.nextToken()
		switch p.curToken.TokenType {
		case token.MOUSE:
			return p.parseSetMouseStatement()
		case token.MODE:
			return p.parseSetModeStatement()
		case token.PAPER:
			return p.parseSetPaperStatement()
		case token.BORDER:
			return p.parseSetBorderStatement()
		case token.PEN:
			return p.parseSetPenStatement()
		case token.DEG:
			return p.parseSetDegStatement()
		case token.RAD:
			return p.parseSetRadStatement()
		case token.CURPOS:
			return p.parseSetCurposStatement()
		case token.COLOUR:
			return p.parseSetColourStatement()
		case token.PATTERN:
			return p.parseSetPatternStatement()
		case token.CONFIG:
			p.nextToken()
			switch p.curToken.TokenType {
			case token.BOOT:
				return p.parseSetConfigBootStatement()
			}
		case token.FILL:
			p.nextToken()
			switch p.curToken.TokenType {
			case token.STYLE:
				return p.parseSetFillStyleStatement()
			}
		default:
			p.errorMsg = syntaxerror.ErrorMessage((syntaxerror.WrongSetAskAttribute))
			p.ErrorTokenIndex = p.curToken.Index
			return nil
		}
	case token.DATA:
		return p.parseDataStatement()
	case token.READ:
		return p.parseReadStatement()
	case token.PRINT:
		return p.parsePrintStatement()
	case token.PLOT:
		return p.parsePlotStatement()
	case token.LINE:
		return p.parseLineStatement()
	case token.AREA:
		return p.parseAreaStatement()
	case token.CIRCLE:
		return p.parseCircleStatement()
	case token.POINTS:
		return p.parsePointsStatement()
	case token.FLOOD:
		return p.parseFloodStatement()
	case token.FETCH:
		return p.parseFetchStatement()
	case token.WRITEBLOCK:
		return p.parseWriteblockStatement()
	case token.READBLOCK:
		return p.parseReadblockStatement()
	case token.COPYBLOCK:
		return p.parseCopyblockStatement()
	case token.SQUASH:
		return p.parseSquashStatement()
	case token.CLEARBLOCK:
		return p.parseClearblockStatement()
	case token.DELBLOCK:
		return p.parseDelblockStatement()
	case token.KEEP:
		return p.parseKeepStatement()
	case token.MOVE:
		return p.parseMoveStatement()
	case token.LET:
		return p.parseLetStatement() // all these methods need to return when they encounter :
	case token.IdentifierLiteral:
		if p.peekTokenIs(token.Equal) || p.peekTokenIs(token.Assign) {
			return p.parseBindStatement()
		}
		if p.peekTokenIs(token.LeftParen) {
			return p.parseBindArrayStatement()
		}
		if p.peekTokenIs(token.EOF) || p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) || p.peekTokenIs(token.IdentifierLiteral) {
			return p.parseProcedureCallStatement()
		}
		// Catch unknown command/procedure-->this needs to depend on the above
		if p.peekTokenIs(token.EOF) || p.peekTokenIs(token.Colon) || p.peekTokenIs(token.NewLine) {
			p.errorMsg = syntaxerror.ErrorMessage(syntaxerror.UnknownCommandProcedure)
			p.ErrorTokenIndex = p.curToken.Index
		}
	case token.RESULT:
		return p.parseResultStatement()
	case token.RESTORE:
		return p.parseRestoreStatement()
	case token.IF:
		return p.parseIfStatement()
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
