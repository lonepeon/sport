package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Slug string

func NewSlug(s string) (Slug, error) {
	if err := validateSlugLength(s); err != nil {
		return Slug(""), err
	}

	if err := validateSlugCharacterSet(s); err != nil {
		return Slug(""), err
	}

	if err := validateSlugCharacterPositions(s); err != nil {
		return Slug(""), err
	}

	return Slug(s), nil
}

func (s Slug) String() string {
	return string(s)
}

func validateSlugLength(s string) error {
	if len(s) < 1 {
		return errors.New("must be greater than 1 character")
	}

	if len(s) > 63 {
		return errors.New("must be less than 63 characters")
	}

	return nil
}

func validateSlugCharacterSet(s string) error {
	for i, c := range s {
		if err := validateSlugCharacterInCharacterSet(c); err != nil {
			return fmt.Errorf("invalid chracter \"%c\" at index %d: %v", c, i, err)
		}
	}

	return nil
}

func validateSlugCharacterInCharacterSet(c rune) error {
	if c >= '0' && c <= '9' {
		return nil
	}

	if c >= 'a' && c <= 'z' {
		return nil
	}

	if c == '.' || c == '-' {
		return nil
	}

	return errors.New("must only contain lower case alphanumeric characters, and - or .")
}

func validateSlugCharacterPositions(s string) error {
	if strings.Contains(s, "--") {
		return fmt.Errorf("invalid sequence \"--\"")
	}

	if strings.Contains(s, "..") {
		return fmt.Errorf("invalid sequence \"..\"")
	}

	if strings.HasPrefix(s, ".") || strings.HasPrefix(s, "-") {
		return fmt.Errorf("must start with an alphanumeric character")
	}

	if strings.HasSuffix(s, ".") || strings.HasSuffix(s, "-") {
		return fmt.Errorf("must end with an alphanumeric character")
	}

	return nil
}
