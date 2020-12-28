package main

type Expr interface {
	// ?
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value string
}

type Unary struct {
	operator Token
	right    Expr
}
