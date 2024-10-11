package searchquerylexer

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Lexer struct {
	input          string
	config         Config
	comparatorList []string
	connectiveList []string

	ch         string
	currentPos int

	currentToken *Token
	prevToken    *Token
	nextToken    *Token
}

func NewLexer(config Config) (*Lexer, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	result := &Lexer{
		config: config,
		comparatorList: []string{
			config.ComparatorConfig.Equal,
			config.ComparatorConfig.NotEqual,
			config.ComparatorConfig.LessThan,
			config.ComparatorConfig.GreaterThan,
			config.ComparatorConfig.LessThanEqualTo,
			config.ComparatorConfig.GreaterThanEqualTo,
			config.ComparatorConfig.Like,
			config.ComparatorConfig.NotLike,
		},
		connectiveList: []string{
			config.ConnectiveConfig.And,
			config.ConnectiveConfig.Or,
		},
	}

	sort.Slice(result.comparatorList, func(i, j int) bool {
		return len(result.comparatorList[i]) > len(result.comparatorList[j])
	})

	sort.Slice(result.connectiveList, func(i, j int) bool {
		return len(result.connectiveList[i]) > len(result.connectiveList[j])
	})

	return result, nil
}

func (l *Lexer) Tokenize(input string) ([]*Token, error) {
	var (
		err error
	)

	l.input = input
	l.currentPos = 0

	result := make([]*Token, 0, 50)

	for !errors.Is(err, io.EOF) {
		if l.currentToken != nil {
			l.prevToken = &Token{
				Type:  l.currentToken.Type,
				Value: l.currentToken.Value,
			}
		}

		l.currentToken, err = l.getNextToken()

		// If we get EOF, all is well, and break
		if errors.Is(err, io.EOF) {
			break
		}

		// If we get here with an error, all is NOT well
		if err != nil {
			return result, err
		}

		// Append to the token list
		result = append(result, l.currentToken)
	}

	return result, nil
}

func (l *Lexer) getNextToken() (*Token, error) {
	var (
		err   error
		value string
	)

	l.skipWhitespace()
	l.readChar()

	if l.ch == "" {
		return NewToken(TokenEOF, ""), io.EOF
	}

	/*
	 * Quoted string
	 */
	if l.isStringStart() {
		value, err = l.captureQuotedValue()

		if err != nil {
			return EmptyToken(), l.captureLinterError(err)
		}

		return NewToken(TokenTypeValue, value), nil
	}

	/*
	 * Subquery
	 */
	if l.isSubqueryStart() {
		return NewToken(TokenTypeSubqueryStart, "("), nil
	}

	if l.isSubqueryEnd() {
		return NewToken(TokenTypeSubqueryEnd, ")"), nil
	}

	/*
	 * Connectives
	 */
	isConnective, connectiveString, err := l.isConnective()

	if err != nil {
		return EmptyToken(), l.captureLinterError(err)
	}

	if isConnective {
		return NewToken(TokenTypeConnective, connectiveString), nil
	}

	/*
	 * Comparators.
	 *
	 * Comparators are configurable, so what we will do is look at the first
	 * character for each one. Then we'll peek for the length of each one
	 * to see if it is a match.
	 */
	isComparator, comparatorString := l.isComparator()

	if isComparator {
		return NewToken(TokenTypeComparator, comparatorString), nil
	}

	/*
	 * If we get here, we have either a value or a field name.
	 * To be a field name, it has to match a registered field
	 * name, and not be preceded by a conmparator. If there are
	 * no registered field names, then it will always be a value.
	 */
	isField, fieldName := l.isField()

	if isField {
		return NewToken(TokenTypeFieldName, fieldName), nil
	}

	/*
	 * If we get here, we have a raw value.
	 */
	value = l.captureRawValue()
	return NewToken(TokenTypeValue, value), nil
}

func (l *Lexer) captureRawValue() string {
	var result strings.Builder

	for l.currentPos <= len(l.input) && l.ch != "" {
		result.WriteString(l.ch)

		if l.isWhitespace(l.currentPos) {
			break
		}

		if l.peekAt(l.currentPos) == ')' {
			break
		}

		l.readChar()
	}

	return result.String()
}

func (l *Lexer) captureQuotedValue() (string, error) {
	var result strings.Builder

	l.ch = string(l.input[l.currentPos])

	for {
		l.readChar()

		// We have an escape sequence
		if l.ch == "\\" {
			l.readChar()

			if l.ch != "\"" && l.ch != "\\" {
				return "", ErrInvalidEscapeSequence
			}
		}

		// We have something to break us out
		if l.ch == "" || l.ch == "\"" {
			break
		}

		result.WriteString(l.ch)
	}

	return result.String(), nil
}

func (l *Lexer) discard(num int) {
	l.currentPos += num
}

func (l *Lexer) readChar() {
	if l.currentPos >= len(l.input) {
		l.ch = ""
	} else {
		ch := l.input[l.currentPos]
		l.ch = string(ch)
	}

	l.currentPos++
}

func (l *Lexer) peek(num int) string {
	if l.currentPos >= len(l.input) {
		return ""
	}

	first := l.currentPos - 1
	last := first + num

	if last >= len(l.input) {
		last = len(l.input) - 1
	}

	result := l.input[first:last]
	return result
}

func (l *Lexer) peekAt(pos int) byte {
	if pos >= len(l.input) {
		return 0
	}

	return l.input[pos]
}

func (l *Lexer) isWhitespace(pos int) bool {
	if pos >= len(l.input) {
		pos = len(l.input) - 1
	}

	ch := l.input[pos]
	return l.chIsWhitespace(ch)
}

func (l *Lexer) chIsWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) skipWhitespace() {
	for l.currentPos < len(l.input) && l.isWhitespace(l.currentPos) {
		l.currentPos++
	}
}

func (l *Lexer) isComparator() (bool, string) {
	for _, comparatorString := range l.comparatorList {
		lenString := len(comparatorString)
		discardLen := lenString - 1
		peek := ""
		areEqual := false

		if lenString < 2 && comparatorString == l.ch {
			areEqual = true
		}

		if lenString >= 2 {
			peek = l.peek(lenString)
			if strings.ToLower(peek) == strings.ToLower(comparatorString) {
				areEqual = true
			}
		}

		if areEqual {
			l.discard(discardLen)
			return true, comparatorString
		}
	}

	return false, ""
}

func (l *Lexer) isStringStart() bool {
	return l.ch == "\""
}

func (l *Lexer) isSubqueryStart() bool {
	return l.ch == "("
}

func (l *Lexer) isSubqueryEnd() bool {
	return l.ch == ")"
}

func (l *Lexer) isConnective() (bool, string, error) {
	isConnective := false
	matchingConnective := ""

	for _, connectiveName := range l.connectiveList {
		peekNum := len(connectiveName) + 1
		peek := strings.ToLower(l.peek(peekNum))

		if peek == strings.ToLower(connectiveName)+" " {
			// This can only be a connective if it is preceeded by a value or subquery
			if l.prevToken != nil && (l.prevToken.Type == TokenTypeValue || l.prevToken.Type == TokenTypeSubqueryEnd) {
				isConnective = true
				matchingConnective = connectiveName

				// There has to be something after a connective. Otherwise
				// it is invalid
				peekStart := l.currentPos + peekNum

				if peekStart >= len(l.input) {
					return false, "", ErrInvalidConnective
				}

				peek = strings.TrimSpace(l.input[peekStart:])

				if peek == "" {
					return false, "", ErrInvalidConnective
				}

				l.discard(len(connectiveName) - 1)
			}
		}
	}

	return isConnective, matchingConnective, nil
}

func (l *Lexer) isField() (bool, string) {
	isFieldName := false
	matchingFieldName := ""

	for _, fieldName := range l.config.FieldNames {
		peekNum := len(fieldName)
		peek := strings.ToLower(l.peek(peekNum))

		if peek == strings.ToLower(fieldName) {
			// We have a potential match. Do we have a preceeding comparator?
			// If so this isn't a field name
			if l.prevToken == nil || l.prevToken.Type != TokenTypeComparator {
				isFieldName = true
				matchingFieldName = fieldName

				peekAt := l.currentPos + peekNum
				discardNum := peekNum - 1

				if l.chIsWhitespace(l.peekAt(peekAt)) {
					discardNum += 1
				}

				l.discard(discardNum)
				break
			}
		}
	}

	return isFieldName, matchingFieldName
}

func (l *Lexer) captureLinterError(originError error) error {
	prefix := "INPUT: "
	s := prefix + l.input + "\n"

	s += fmt.Sprintf("%*s\n", l.currentPos+len(prefix), "│")
	s += fmt.Sprintf("%*s %s\n", l.currentPos+len(prefix), "└", l.prettyError(originError))

	return fmt.Errorf("%s: %w", s, originError)
}

func (l *Lexer) prettyError(err error) string {
	if errors.Is(err, ErrInvalidConnective) {
		return "invalid boolean operator. boolean operators must have two conditions"
	}

	return err.Error()
}
