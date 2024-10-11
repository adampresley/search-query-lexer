package searchquerylexer

import (
	"fmt"
	"strings"
)

type Config struct {
	ComparatorConfig ComparatorConfig
	ConnectiveConfig ConnectiveConfig
	FieldNames       []string
}

type ComparatorConfig struct {
	Equal              string
	NotEqual           string
	LessThan           string
	GreaterThan        string
	LessThanEqualTo    string
	GreaterThanEqualTo string
	Like               string
	NotLike            string
}

type ConnectiveConfig struct {
	And string
	Or  string
}

func (c Config) validate() error {
	if strings.TrimSpace(c.ComparatorConfig.Equal) == "" {
		return fmt.Errorf("missing EQUAL configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.NotEqual) == "" {
		return fmt.Errorf("missing NOT EQUAL configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.LessThan) == "" {
		return fmt.Errorf("missing LESS THAN configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.GreaterThan) == "" {
		return fmt.Errorf("missing GREATER THAN configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.LessThanEqualTo) == "" {
		return fmt.Errorf("missing LESS THAN EQUAL configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.GreaterThanEqualTo) == "" {
		return fmt.Errorf("missing GREATER THAN EQUAL configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.Like) == "" {
		return fmt.Errorf("missing LIKE configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ComparatorConfig.NotLike) == "" {
		return fmt.Errorf("missing NOT LIKE configuration: %w", ErrInvalidConfigComparator)
	}

	if strings.TrimSpace(c.ConnectiveConfig.And) == "" {
		return fmt.Errorf("missing AND configuration: %w", ErrInvalidConfigConnective)
	}

	if strings.TrimSpace(c.ConnectiveConfig.Or) == "" {
		return fmt.Errorf("missing OR configuration: %w", ErrInvalidConfigConnective)
	}

	return nil
}
