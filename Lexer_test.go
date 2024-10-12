package searchquerylexer_test

import (
	"fmt"
	"testing"

	sql "github.com/adampresley/search-query-lexer"
	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	t.Run("success with config", func(t *testing.T) {
		config := sql.Config{
			ComparatorConfig: sql.DefaultComparatorConfig,
			ConnectiveConfig: sql.DefaultConnectiveConfig,
			FieldNames: []string{
				"title",
				"name",
			},
		}

		lexer, err := sql.NewLexer(config)

		assert.NoError(t, err)
		assert.IsType(t, &sql.Lexer{}, lexer)
	})
}

func TestTokenize(t *testing.T) {
	defaultConfig := sql.Config{
		ComparatorConfig: sql.DefaultComparatorConfig,
		ConnectiveConfig: sql.DefaultConnectiveConfig,
		FieldNames: []string{
			"title",
			"name",
			"age",
			"category",
		},
	}

	alternateConfig := sql.Config{
		ComparatorConfig: sql.ComparatorConfig{
			Equal:              ":",
			NotEqual:           ":!",
			LessThan:           ":<",
			GreaterThan:        ":>",
			LessThanEqualTo:    ":<=",
			GreaterThanEqualTo: ":>=",
			Like:               "LIKE",
			NotLike:            "!LIKE",
		},
		ConnectiveConfig: sql.ConnectiveConfig{
			And: "&&",
			Or:  "||",
		},
		FieldNames: []string{
			"title",
			"name",
			"age",
			"role",
		},
	}

	table := []struct {
		name        string
		input       string
		want        []*sql.Token
		wantErr     bool
		expectedErr error
		config      sql.Config
	}{
		{
			name:  "Equal test (no spaces)",
			input: "title=testing",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "testing"),
			},
			config: defaultConfig,
		},
		{
			name:  "Equal test (spaces)",
			input: "title = testing",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "testing"),
			},
			config: defaultConfig,
		},
		{
			name:  "Not equal test (no spaces)",
			input: "name!=bob",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "!="),
				sql.NewToken(sql.TokenTypeValue, "bob"),
			},
			config: defaultConfig,
		},
		{
			name:  "Not equal test (spaces)",
			input: "name != bob",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "!="),
				sql.NewToken(sql.TokenTypeValue, "bob"),
			},
			config: defaultConfig,
		},
		{
			name:  "OR test",
			input: "title=testing OR name != \"bob\"",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "testing"),
				sql.NewToken(sql.TokenTypeConnective, "or"),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "!="),
				sql.NewToken(sql.TokenTypeValue, "bob"),
			},
			config: defaultConfig,
		},
		{
			name:  "value and value",
			input: "yummy and sweet",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeValue, "yummy"),
				sql.NewToken(sql.TokenTypeConnective, "and"),
				sql.NewToken(sql.TokenTypeValue, "sweet"),
			},
			config: defaultConfig,
		},
		{
			name:  "value or value",
			input: "salty or sweet",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeValue, "salty"),
				sql.NewToken(sql.TokenTypeConnective, "or"),
				sql.NewToken(sql.TokenTypeValue, "sweet"),
			},
			config: defaultConfig,
		},
		{
			name:  "subquery test (no spaces)",
			input: "title =~ testing AND (name=\"Adam\" OR name=\"Bob\")",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "=~"),
				sql.NewToken(sql.TokenTypeValue, "testing"),
				sql.NewToken(sql.TokenTypeConnective, "and"),
				sql.NewToken(sql.TokenTypeSubqueryStart, "("),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "Adam"),
				sql.NewToken(sql.TokenTypeConnective, "or"),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "Bob"),
				sql.NewToken(sql.TokenTypeSubqueryEnd, ")"),
			},
			config: defaultConfig,
		},
		{
			name:  "subquery test (spaces)",
			input: "title =~ testing AND ( name=\"Adam\" OR name = \"Bob\" )",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "=~"),
				sql.NewToken(sql.TokenTypeValue, "testing"),
				sql.NewToken(sql.TokenTypeConnective, "and"),
				sql.NewToken(sql.TokenTypeSubqueryStart, "("),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "Adam"),
				sql.NewToken(sql.TokenTypeConnective, "or"),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, "="),
				sql.NewToken(sql.TokenTypeValue, "Bob"),
				sql.NewToken(sql.TokenTypeSubqueryEnd, ")"),
			},
			config: defaultConfig,
		},
		{
			name:  "starts with subquery",
			input: `(title=~"test" AND age >= 30) OR (category != "bad")`,
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeSubqueryStart, "("),
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, "=~"),
				sql.NewToken(sql.TokenTypeValue, "test"),
				sql.NewToken(sql.TokenTypeConnective, "and"),
				sql.NewToken(sql.TokenTypeFieldName, "age"),
				sql.NewToken(sql.TokenTypeComparator, ">="),
				sql.NewToken(sql.TokenTypeValue, "30"),
				sql.NewToken(sql.TokenTypeSubqueryEnd, ")"),
				sql.NewToken(sql.TokenTypeConnective, "or"),
				sql.NewToken(sql.TokenTypeSubqueryStart, "("),
				sql.NewToken(sql.TokenTypeFieldName, "category"),
				sql.NewToken(sql.TokenTypeComparator, "!="),
				sql.NewToken(sql.TokenTypeValue, "bad"),
				sql.NewToken(sql.TokenTypeSubqueryEnd, ")"),
			},
			config: defaultConfig,
		},
		{
			name:  "alternate config",
			input: "title:1 name:!2 && (age :> 23 || role : \"admin\")",
			want: []*sql.Token{
				sql.NewToken(sql.TokenTypeFieldName, "title"),
				sql.NewToken(sql.TokenTypeComparator, ":"),
				sql.NewToken(sql.TokenTypeValue, "1"),
				sql.NewToken(sql.TokenTypeFieldName, "name"),
				sql.NewToken(sql.TokenTypeComparator, ":!"),
				sql.NewToken(sql.TokenTypeValue, "2"),
				sql.NewToken(sql.TokenTypeConnective, "&&"),
				sql.NewToken(sql.TokenTypeSubqueryStart, "("),
				sql.NewToken(sql.TokenTypeFieldName, "age"),
				sql.NewToken(sql.TokenTypeComparator, ":>"),
				sql.NewToken(sql.TokenTypeValue, "23"),
				sql.NewToken(sql.TokenTypeConnective, "||"),
				sql.NewToken(sql.TokenTypeFieldName, "role"),
				sql.NewToken(sql.TokenTypeComparator, ":"),
				sql.NewToken(sql.TokenTypeValue, "admin"),
				sql.NewToken(sql.TokenTypeSubqueryEnd, ")"),
			},
			config: alternateConfig,
		},
		{
			name:        "invalid escape sequence error",
			input:       `title="\atest"`,
			want:        nil,
			wantErr:     true,
			expectedErr: sql.ErrInvalidEscapeSequence,
			config:      defaultConfig,
		},
		{
			name:        "invalid connective",
			input:       "name != \"no cap\" OR     ",
			want:        nil,
			wantErr:     true,
			expectedErr: sql.ErrInvalidConnective,
			config:      defaultConfig,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			lexer, err := sql.NewLexer(tt.config)
			assert.NoError(t, err)

			got, err := lexer.Tokenize(tt.input)
			fmt.Printf("INPUT:\n%s\n", tt.input)

			if tt.wantErr {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				for _, t := range got {
					fmt.Printf("%s\n", t.String())
				}
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
