package main

type Parser struct {
	tokens    []Token
	position  int
	curToken  Token
	peekToken Token
}

func (p *Parser) New(tokens []Token) {
	p.tokens = tokens
	p.position = 0
	p.peekToken = tokens[p.position]
	// Read a token so curToken and peekToken are both set
	p.nextToken()
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.position++
	p.peekToken = p.tokens[p.position]
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.peekToken.TokenType != EndOfLine {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.TokenType {
	case LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}
	if !p.expectPeek(IdentifierLiteral) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(Assign) {
		return nil
	}

	//p.nextToken()

	// To do: Expressions
	for !p.curTokenIs(Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t int) bool {
	return p.curToken.TokenType == t
}

func (p *Parser) peekTokenIs(t int) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) expectPeek(t int) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
