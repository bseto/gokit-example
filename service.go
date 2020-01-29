package main

import (
	"errors"
	"strings"
)

type stringService struct{}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
	Lowercase(string) (string, error)
}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func (stringService) Lowercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToLower(s), nil
}

// ServiceMiddleware is a chainable behavior modifier for StringService.
type ServiceMiddleware func(StringService) StringService
