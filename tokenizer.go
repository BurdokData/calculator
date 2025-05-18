package main

import (
	"unicode"
)

type TokenType string

type Token struct {
	Type TokenType
	Value string
}

const(
	PLUS TokenType = "PLUS"
	MINUS TokenType = "MINUS"
	MUL TokenType = "MUL"
	DIV TokenType = "DIV"
	LPAREN TokenType = "LPAREN"
	RPAREN TokenType = "RPAREN"
	COMMA TokenType = "COMMA"
	NUM TokenType = "NUM"
	IDENTIFIER TokenType = "IDENTIFIER"
	EOF TokenType = "EOF"
	ANS TokenType = "ANS"
	CLEAR TokenType = "CLEAR"
	EXIT TokenType = "EXIT"
	PI TokenType = "PI"
	E TokenType = "E"
)

func Tokenize(input string) <-chan Token {
	ret := make(chan Token, 10)
	go func() {
		defer close(ret)

		i := 0
		for i < len(input) {
			if unicode.IsSpace(rune(input[i])) {
				i++
				continue
			}

			if unicode.IsDigit(rune(input[i])) {
				start := i
				hasDot := false

				numLoop: for i < len(input) {
					c := rune(input[i])
					switch {
					case unicode.IsDigit(c):
						i++
					case c == '.' && !hasDot:
						hasDot = true
						i++
					default:
						break numLoop
					}
				}

				ret <- Token{Type: NUM, Value: input[start:i]}
				continue
			}

			if unicode.IsLetter(rune(input[i])) {
				start := i
				letterLoop: for i < len(input) {
					c := rune(input[i])
					switch {
					case unicode.IsLetter(c):
						i++
					default:
						break letterLoop
					}
				}
				value := input[start:i]
				switch value {
				case "ans":
					ret <- Token{Type: ANS}
				case "exit":
					ret <- Token{Type: EXIT}
				case "clear":
					ret <- Token{Type: CLEAR}
				case "pi":
					ret <- Token{Type: PI}
				case "e":
					ret <- Token{Type: E}
				default:
					ret <- Token{Type: IDENTIFIER, Value: input[start:i]}
				}
				continue
			}

			switch rune(input[i]) {
			case '+':
				ret <- Token{Type: PLUS}
			case '-':
				ret <- Token{Type: MINUS}
			case '(':
				ret <- Token{Type: LPAREN}
			case ')':
				ret <- Token{Type: RPAREN}
			case '*':
				ret <- Token{Type: MUL}
			case '/':
				ret <- Token{Type: DIV}
			case ',':
				ret <- Token{Type: COMMA}
			}
			i++
		}

		ret<-Token{Type:EOF}
	}()
	return ret
}
