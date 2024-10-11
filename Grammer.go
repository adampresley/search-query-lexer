package searchquerylexer

type TokenType string

const (
	TokenEmpty             TokenType = "[empty]"
	TokenTypeValue         TokenType = "[value]"
	TokenTypeComparator    TokenType = "[comparator]"
	TokenTypeFieldName     TokenType = "[fieldName]"
	TokenTypeSubqueryStart TokenType = "[subQueryStart]"
	TokenTypeSubqueryEnd   TokenType = "[subQueryEnd]"
	TokenTypeConnective    TokenType = "[connective]"
	TokenEOF               TokenType = "[eof]"
)
