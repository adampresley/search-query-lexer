package main

import (
	"fmt"
	"os"

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

	for _, t := range tokens {
		fmt.Printf("%s\n", t.String())
	}
}
