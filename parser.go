package main

type Parser struct {
	tokens    []Token
	position  int
	curToken  Token
	peekToken Token
}

func New(tokens []Token) *Parser {
	p := &Parser{}
	p.tokens = tokens
	p.position = 0
	p.peekToken = tokens[p.position]
	// Read a token so curToken and peekToken are both set
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.position++
	p.peekToken = p.tokens[p.position]
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.TokenType != EndOfLine {
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
