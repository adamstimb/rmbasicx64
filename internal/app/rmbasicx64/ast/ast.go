package ast

import (
	"bytes"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// The root of the AST
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// redundant?
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// This represents a line of code, which code contain several statements
type Line struct {
	Statements []Statement
	LineNumber int
	LineString string
}

func (l *Line) lineNode() {}
func (l *Line) String() string {
	var out bytes.Buffer
	for i, s := range l.Statements {
		out.WriteString(s.String())
		if i < len(l.Statements)-1 {
			out.WriteString(" : ")
		}
	}
	return out.String()
}
func (l *Line) TokenLiteral() string {
	if len(l.Statements) > 0 {
		return l.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// redundant?
type BlockStatement struct {
	Token      token.Token // then { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type IfExpression struct {
	Token       token.Token // the If token
	Condition   Expression
	Consequence *Line
	Alternative *Line
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("IF")
	out.WriteString(ie.Condition.String())
	out.WriteString(" THEN")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("ELSE ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type IfStatement struct {
	Token       token.Token // the If token
	Condition   Expression
	Consequence *Line
	Alternative *Line
}

func (s *IfStatement) statementNode() {}
func (s *IfStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("IF")
	out.WriteString(s.Condition.String())
	out.WriteString(" THEN")
	out.WriteString(s.Consequence.String())
	if s.Alternative != nil {
		out.WriteString("ELSE ")
		out.WriteString(s.Alternative.String())
	}
	return out.String()
}

type UntilStatement struct {
	Token     token.Token
	Condition Expression
}

func (s *UntilStatement) statementNode() {}
func (s *UntilStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *UntilStatement) String() string {
	var out bytes.Buffer
	out.WriteString("UNTIL ")
	out.WriteString(s.Condition.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type NumericLiteral struct {
	Token token.Token
	Value float64
}

func (il NumericLiteral) expressionNode() {}
func (il *NumericLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *NumericLiteral) String() string {
	return il.Token.Literal
}

type Identifier struct {
	Token      token.Token // the token.IdentifierLiteral token
	Value      string
	Subscripts []Expression
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. NOT
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type ByeStatement struct {
	Token token.Token // the token.Bye token
}

func (bs *ByeStatement) statementNode() {}
func (bs *ByeStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *ByeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral())
	return out.String()
}

type ListStatement struct {
	Token token.Token // the token.Bye token
}

func (s *ListStatement) statementNode() {}
func (s *ListStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ListStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type RunStatement struct {
	Token token.Token // the token.Bye token
}

func (s *RunStatement) statementNode() {}
func (s *RunStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *RunStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type NewStatement struct {
	Token token.Token // the token.Bye token
}

func (s *NewStatement) statementNode() {}
func (s *NewStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *NewStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type ClearblockStatement struct {
	Token token.Token // the token.Bye token
}

func (s *ClearblockStatement) statementNode() {}
func (s *ClearblockStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ClearblockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type DelblockStatement struct {
	Token token.Token // the token.Bye token
	Block Expression
}

func (s *DelblockStatement) statementNode() {}
func (s *DelblockStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *DelblockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type KeepStatement struct {
	Token token.Token // the token.Bye token
	Block Expression
	Path  Expression
}

func (s *KeepStatement) statementNode() {}
func (s *KeepStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *KeepStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type RenumberStatement struct {
	Token token.Token // the token.Bye token
}

func (s *RenumberStatement) statementNode() {}
func (s *RenumberStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *RenumberStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type ClsStatement struct {
	Token token.Token
}

func (s *ClsStatement) statementNode() {}
func (s *ClsStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ClsStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetModeStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetModeStatement) statementNode() {}
func (s *SetModeStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetModeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetPaperStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetPaperStatement) statementNode() {}
func (s *SetPaperStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetPaperStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetBorderStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetBorderStatement) statementNode() {}
func (s *SetBorderStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetBorderStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetPenStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetPenStatement) statementNode() {}
func (s *SetPenStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetPenStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetMouseStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetMouseStatement) statementNode() {}
func (s *SetMouseStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetMouseStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type GetStatement struct {
	Token token.Token
	Value Expression
}

func (s *GetStatement) statementNode() {}
func (s *GetStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *GetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetCurposStatement struct {
	Token token.Token
	Col   Expression
	Row   Expression
}

func (s *SetCurposStatement) statementNode() {}
func (s *SetCurposStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetCurposStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetPatternStatement struct {
	Token token.Token
	Slot  Expression
	Row   Expression
	C1    Expression
	C2    Expression
	C3    Expression
	C4    Expression
}

func (s *SetPatternStatement) statementNode() {}
func (s *SetPatternStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetPatternStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type MoveStatement struct {
	Token token.Token
	Cols  Expression
	Rows  Expression
}

func (s *MoveStatement) statementNode() {}
func (s *MoveStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *MoveStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type HomeStatement struct {
	Token token.Token
}

func (s *HomeStatement) statementNode() {}
func (s *HomeStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *HomeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type DirStatement struct {
	Token token.Token
}

func (s *DirStatement) statementNode() {}
func (s *DirStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *DirStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetConfigBootStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetConfigBootStatement) statementNode() {}
func (s *SetConfigBootStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetConfigBootStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetColourStatement struct {
	Token       token.Token
	PaletteSlot Expression
	BasicColour Expression
	FlashSpeed  Expression
	FlashColour Expression
}

func (s *SetColourStatement) statementNode() {}
func (s *SetColourStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetColourStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetDegStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetDegStatement) statementNode() {}
func (s *SetDegStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetDegStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type SetRadStatement struct {
	Token token.Token
	Value Expression
}

func (s *SetRadStatement) statementNode() {}
func (s *SetRadStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SetRadStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type RemStatement struct {
	Token   token.Token
	Comment token.Token
}

func (s *RemStatement) statementNode() {}
func (s *RemStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *RemStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type PrintStatement struct {
	Token     token.Token
	Value     Expression
	PrintList []interface{}
}

func (ps *PrintStatement) statementNode() {}
func (ps *PrintStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *PrintStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	out.WriteString(ps.Value.String())
	return out.String()
}

type PlotStatement struct {
	Token     token.Token
	Value     Expression
	CoordList []Expression
	Brush     Expression
	Direction Expression
	SizeX     Expression
	SizeY     Expression
	Font      Expression
	Over      Expression
}

func (ps *PlotStatement) statementNode() {}
func (ps *PlotStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *PlotStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	out.WriteString(ps.Value.String())
	return out.String()
}

type LineStatement struct {
	Token     token.Token
	CoordList []Expression
	Brush     Expression
	Over      Expression
}

func (ps *LineStatement) statementNode() {}
func (ps *LineStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *LineStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	return out.String()
}

type CircleStatement struct {
	Token        token.Token
	Radius       Expression
	CoordList    []Expression
	Brush        Expression
	Over         Expression
	FillStyle    Expression
	FillHatching Expression
	FillColour2  Expression
}

func (s *CircleStatement) statementNode() {}
func (s *CircleStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *CircleStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type PointsStatement struct {
	Token     token.Token
	CoordList []Expression
	Style     Expression
	Brush     Expression
	Over      Expression
}

func (s *PointsStatement) statementNode() {}
func (s *PointsStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *PointsStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type FloodStatement struct {
	Token         token.Token
	CoordList     []Expression
	Brush         Expression
	UseEdgeColour Expression
	EdgeColour    Expression
	FillStyle     Expression
	FillHatching  Expression
	FillColour2   Expression
}

func (s *FloodStatement) statementNode() {}
func (s *FloodStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *FloodStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type FetchStatement struct {
	Token token.Token
	Block Expression
	Path  Expression
}

func (s *FetchStatement) statementNode() {}
func (s *FetchStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *FetchStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type WriteblockStatement struct {
	Token token.Token
	Block Expression
	X     Expression
	Y     Expression
	Over  Expression
}

func (s *WriteblockStatement) statementNode() {}
func (s *WriteblockStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *WriteblockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type ReadblockStatement struct {
	Token token.Token
	Block Expression
	X1    Expression
	Y1    Expression
	X2    Expression
	Y2    Expression
}

func (s *ReadblockStatement) statementNode() {}
func (s *ReadblockStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ReadblockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type CopyblockStatement struct {
	Token token.Token
	X1    Expression
	Y1    Expression
	X2    Expression
	Y2    Expression
	Dx    Expression
	Dy    Expression
	Over  Expression
}

func (s *CopyblockStatement) statementNode() {}
func (s *CopyblockStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *CopyblockStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type SquashStatement struct {
	Token token.Token
	Block Expression
	X     Expression
	Y     Expression
	Over  Expression
}

func (s *SquashStatement) statementNode() {}
func (s *SquashStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SquashStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type DataStatement struct {
	Token    token.Token
	ItemList []token.Token
}

func (ps *DataStatement) statementNode() {}
func (ps *DataStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *DataStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	return out.String()
}

type ReadStatement struct {
	Token        token.Token
	VariableList []*Identifier
}

func (ps *ReadStatement) statementNode() {}
func (ps *ReadStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *ReadStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	return out.String()
}

type AreaStatement struct {
	Token        token.Token
	CoordList    []Expression
	Brush        Expression
	Over         Expression
	FillStyle    Expression
	FillHatching Expression
	FillColour2  Expression
}

func (ps *AreaStatement) statementNode() {}
func (ps *AreaStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *AreaStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	return out.String()
}

type SetFillStyleStatement struct {
	Token        token.Token
	FillStyle    Expression
	FillHatching Expression
	FillColour2  Expression
}

func (ps *SetFillStyleStatement) statementNode() {}
func (ps *SetFillStyleStatement) TokenLiteral() string {
	return ps.Token.Literal
}
func (ps *SetFillStyleStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	return out.String()
}

type SaveStatement struct {
	Token token.Token
	Value Expression
}

func (s *SaveStatement) statementNode() {}
func (s *SaveStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SaveStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Value.String())
	return out.String()
}

type LoadStatement struct {
	Token token.Token
	Value Expression
}

func (s *LoadStatement) statementNode() {}
func (s *LoadStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *LoadStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Value.String())
	return out.String()
}

type GotoStatement struct {
	Token      token.Token
	Linenumber token.Token
}

func (s *GotoStatement) statementNode() {}
func (s *GotoStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *GotoStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type RestoreStatement struct {
	Token      token.Token
	Linenumber token.Token
}

func (s *RestoreStatement) statementNode() {}
func (s *RestoreStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *RestoreStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type EditStatement struct {
	Token      token.Token
	Linenumber token.Token
}

func (s *EditStatement) statementNode() {}
func (s *EditStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *EditStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type RepeatStatement struct {
	Token           token.Token
	LineNumber      int
	StatementNumber int
}

func (s *RepeatStatement) statementNode() {}
func (s *RepeatStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *RepeatStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	return out.String()
}

type LetStatement struct {
	Token     token.Token // the token.Let token
	BindToken token.Token // either then = or := token
	Name      *Identifier
	Value     Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" " + ls.BindToken.Literal + " ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

type ForStatement struct {
	Token           token.Token // the token.Let token
	Name            *Identifier
	BindToken       token.Token
	Start           Expression
	Stop            Expression
	Step            Expression
	LineNumber      int
	StatementNumber int
	StartValue      float64
	StopValue       float64
	StepValue       float64
}

func (s *ForStatement) statementNode() {}
func (s *ForStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type SubroutineStatement struct {
	Token           token.Token
	Name            *Identifier
	LineNumber      int
	StatementNumber int
}

func (s *SubroutineStatement) statementNode() {}
func (s *SubroutineStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SubroutineStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type GosubStatement struct {
	Token           token.Token
	Name            *Identifier
	LineNumber      int
	StatementNumber int
}

func (s *GosubStatement) statementNode() {}
func (s *GosubStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *GosubStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type FunctionDeclaration struct {
	Token           token.Token
	Name            *Identifier
	LineNumber      int
	StatementNumber int
	ReceiveArgs     []*Identifier
}

func (s *FunctionDeclaration) statementNode() {}
func (s *FunctionDeclaration) TokenLiteral() string {
	return s.Token.Literal
}
func (s *FunctionDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type EndfunStatement struct {
	Token token.Token
}

func (bs *EndfunStatement) statementNode() {}
func (bs *EndfunStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *EndfunStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral())
	return out.String()
}

type ReturnStatement struct {
	Token token.Token
}

func (s *ReturnStatement) statementNode() {}
func (s *ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type ProcedureDeclaration struct {
	Token           token.Token
	Name            *Identifier
	LineNumber      int
	StatementNumber int
	ReturnArgs      []Expression
}

func (s *ProcedureDeclaration) statementNode() {}
func (s *ProcedureDeclaration) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ProcedureDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type LeaveStatement struct {
	Token token.Token
}

func (bs *LeaveStatement) statementNode() {}
func (bs *LeaveStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *LeaveStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral())
	return out.String()
}

type EndprocStatement struct {
	Token token.Token
}

func (bs *EndprocStatement) statementNode() {}
func (bs *EndprocStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *EndprocStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral())
	return out.String()
}

type EndStatement struct {
	Token token.Token
}

func (bs *EndStatement) statementNode() {}
func (bs *EndStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *EndStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.TokenLiteral())
	return out.String()
}

type NextStatement struct {
	Token token.Token // the token.Let token
	Name  *Identifier
}

func (s *NextStatement) statementNode() {}
func (s *NextStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *NextStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type DimStatement struct {
	Token      token.Token
	Name       *Identifier
	Subscripts []Expression
}

func (s *DimStatement) statementNode() {}
func (s *DimStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *DimStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	return out.String()
}

type AskMouseStatement struct {
	Token token.Token // the token.Let token
	XName *Identifier
	YName *Identifier
	BName *Identifier
}

func (s *AskMouseStatement) statementNode() {}
func (s *AskMouseStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *AskMouseStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type AskBlocksizeStatement struct {
	Token  token.Token
	Block  Expression
	Width  *Identifier
	Height *Identifier
	Mode   *Identifier
}

func (s *AskBlocksizeStatement) statementNode() {}
func (s *AskBlocksizeStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *AskBlocksizeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral() + " ")
	return out.String()
}

type BindStatement struct {
	Name  *Identifier
	Token token.Token // the token.Equal or token.Assign
	Value Expression
}

func (bs *BindStatement) statementNode() {}
func (bs *BindStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BindStatement) String() string {
	var out bytes.Buffer
	out.WriteString(bs.Name.String())
	out.WriteString(" " + bs.TokenLiteral() + " ")
	if bs.Value != nil {
		out.WriteString(bs.Value.String())
	}
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

type ResultStatement struct {
	Token       token.Token
	ResultValue Expression
}

func (rs *ResultStatement) statementNode() {}
func (rs *ResultStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ResultStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ResultValue != nil {
		out.WriteString(rs.ResultValue.String())
	}
	out.WriteString(";")
	return out.String()
}
