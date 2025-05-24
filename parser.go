package main

import (
	"fmt"
	"strconv"
	"math"
)

type UnexpectedTokenError struct {
	UnexpectedToken TokenType
}

func (ute UnexpectedTokenError) Error() string {
	return fmt.Sprintf("unexpected token %v", ute.UnexpectedToken)
}

type Parser struct {
	Tokens <-chan Token
	NextToken Token
}

func (p *Parser) Eat(expected TokenType) (Token, error) {
	if p.NextToken.Type == expected {
		ret := p.NextToken
		p.NextToken = <-p.Tokens
		return ret, nil
	}
	return Token{}, UnexpectedTokenError{UnexpectedToken: p.NextToken.Type}
}

func NewParser(tokens <-chan Token) Parser {
	return Parser{
		Tokens: tokens,
		NextToken: <-tokens,
	}
}

func (p *Parser) Expression() (Expr, error) {
	var ret Expr
	if p.NextToken.Type == PLUS || p.NextToken.Type == MINUS {
		ret = NumberNode(lastAnswer)
	} else {
		expr, err := p.parseMulDiv()

		if err != nil {
			return nil, err
		}
		
		ret = expr
	}

	for p.NextToken.Type == PLUS || p.NextToken.Type == MINUS {
		op := p.NextToken.Type
		p.Eat(op)

		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}

		ret = BinaryOpNode{Op: op, Left: ret, Right: right}
	}

	return ret, nil
}

func (p *Parser) parseMulDiv() (Expr, error) {
	var ret Expr
	if p.NextToken.Type == MUL || p.NextToken.Type == DIV {
		ret = NumberNode(lastAnswer)
	} else {
		expr, err := p.parsePrimitive()
		if err != nil {
			return nil, err
		}
		ret = expr
	}

	for p.NextToken.Type == MUL || p.NextToken.Type == DIV {
		op := p.NextToken.Type
		p.Eat(op)

		right, err := p.parsePrimitive()
		if err != nil {
			return nil, err
		}

		ret = BinaryOpNode{Op: op, Left: ret, Right: right}
	}

	return ret, nil
}

type ClearError struct{}

func (ce ClearError) Error() string {
	return ""
}

type ExitError struct{}

func (ee ExitError) Error() string {
	return ""
}

func (p *Parser) parsePrimitive() (Expr, error) {
	var ret Expr
	if p.NextToken.Type == NUM {
		num, err := p.Eat(NUM)
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseFloat(num.Value, 64)
		if err != nil {
			return nil, err
		}
		ret = NumberNode(val)
	} else if p.NextToken.Type == IDENTIFIER {
		identifier, err := p.Eat(IDENTIFIER)
		_, err = p.Eat(LPAREN)
		if err != nil {
			return nil, err
		}
		args := []Expr{}
		if p.NextToken.Type != RPAREN {
			for {
				arg, err := p.Expression()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				if p.NextToken.Type == COMMA {
					_, err = p.Eat(COMMA)
					if err != nil {
						return nil, err
					}
				} else {
					break
				}
			}
		}
		_, err = p.Eat(RPAREN)
		if err != nil && p.NextToken.Type != EOF {
			return nil, err
		}
		ret = FunctionNode{Name: identifier.Value, Args: args}
	} else if p.NextToken.Type == LPAREN {
		p.Eat(LPAREN)
		exp, err := p.Expression()
		if err != nil {
			return nil, err
		}
		ret = exp
		p.Eat(RPAREN)
	} else if p.NextToken.Type == ANS {
		p.Eat(ANS)
		ret = NumberNode(lastAnswer)
	} else if p.NextToken.Type == CLEAR {
		return nil, ClearError{}
	} else if p.NextToken.Type == EXIT {
		return nil, ExitError{}
	} else if p.NextToken.Type == PI {
		p.Eat(PI)
		ret = NumberNode(math.Pi)
	} else if p.NextToken.Type == E {
		p.Eat(E)
		ret = NumberNode(math.E)
	} else {
		return nil, UnexpectedTokenError{UnexpectedToken: p.NextToken.Type}
	}
	return ret, nil
}
