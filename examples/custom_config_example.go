package main

import (
	"fmt"
	"os"

	searchquerylexer "github.com/adampresley/search-query-lexer"
)

func main() {
	config := searchquerylexer.Config{
		ComparatorConfig: searchquerylexer.ComparatorConfig{
			Equal:              "EQ",
			NotEqual:           "NEQ",
			LessThan:           "LT",
			GreaterThan:        "GT",
			LessThanEqualTo:    "LTE",
			GreaterThanEqualTo: "GTE",
			Like:               "LIKE",
			NotLike:            "!LIKE",
		},
		ConnectiveConfig: searchquerylexer.ConnectiveConfig{
			And: "&&",
			Or:  "||",
		},
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

	input := `(title LIKE "test" && age GTE 30) || (category NEQ "bad")`

	tokens, err := lexer.Tokenize(input)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	for _, t := range tokens {
		fmt.Printf("%s\n", t.String())
	}
}
