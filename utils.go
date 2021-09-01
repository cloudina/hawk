package main

import (
	// standard
	"fmt"
)

type StringArgs []string

func (s *StringArgs) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StringArgs) Set(value string) error {
	*s = append(*s, value)

	return nil
}
