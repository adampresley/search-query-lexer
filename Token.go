package searchquerylexer

import "fmt"

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}

func EmptyToken() *Token {
	return &Token{Type: TokenEmpty, Value: ""}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s: '%s'", t.Type, t.Value)
}
