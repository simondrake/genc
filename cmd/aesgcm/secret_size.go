// secretSize implements a custom type to be used with Cobra.
//
// It ensures that the size value is set to one of 16, 24, or 32.

package aesgcm

import (
	"errors"
	"fmt"
	"strconv"
)

type secretSize int

const (
	secretSize16 secretSize = 16
	secretSize24 secretSize = 24
	secretSize32 secretSize = 32
)

func (s *secretSize) String() string {
	return strconv.Itoa(int(*s))
}

func (s *secretSize) Set(v string) error {
	switch v {
	case "16", "24", "32":
		i, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error converting to int: %w", err)
		}

		*s = secretSize(i)

		return nil
	default:
		return errors.New(`must be one of 16, 24, or 32`)
	}
}

func (s *secretSize) Type() string {
	return "[16,24,32]"
}
