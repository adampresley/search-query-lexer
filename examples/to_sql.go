package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	searchquerylexer "github.com/adampresley/search-query-lexer"
)

func main() {
	config := searchquerylexer.Config{
		ComparatorConfig: searchquerylexer.DefaultComparatorConfig,
		ConnectiveConfig: searchquerylexer.DefaultConnectiveConfig,
		FieldNames: []string{
			"title",
			"age",
			"category",
		},
	}

	lexer, err := searchquerylexer.NewLexer(config)

	if err != nil {
		fmt.Printf("error initializing lexer: %s\n", err.Error())
		os.Exit(1)
	}

	input := `(title=~"test" AND age >= 30) OR (category != "bad")`

	tokens, err := lexer.Tokenize(input)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	result := strings.Builder{}
	inLike := false

	for _, t := range tokens {
		switch t.Type {
		case searchquerylexer.TokenTypeComparator:
			switch t.Value {
			case "!=":
				result.WriteString(" <> ")

			case "=~":
				inLike = true
				result.WriteString(" LIKE '%")

			case "!~":
				inLike = true
				result.WriteString(" NOT LIKE '%")

			default:
				result.WriteString(" " + t.Value + " ")
			}

		case searchquerylexer.TokenTypeValue:
			if inLike {
				result.WriteString(t.Value + "%' ")
				inLike = false
			} else {
				// Is this a number?
				if _, err := strconv.Atoi(t.Value); err == nil {
					result.WriteString(t.Value)
				} else {
					result.WriteString("'" + t.Value + "' ")
				}
			}

		case searchquerylexer.TokenTypeConnective:
			result.WriteString(" " + strings.ToUpper(t.Value) + " ")

		default:
			result.WriteString(t.Value)
		}
	}

	fmt.Printf("SQL:\n%s\n", result.String())
}
