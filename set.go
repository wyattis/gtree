package main

import (
	"fmt"
	"strings"
)

func NewSet() *Set {
	return &Set{
		data: make(map[string]bool),
	}
}

type Set struct {
	data map[string]bool
}

func (s *Set) Set(val string) error {
	vals := strings.Split(val, ",")
	for _, v := range vals {
		s.data[v] = true
	}
	return nil
}

func (s Set) Has(val string) bool {
	_, ok := s.data[val]
	return ok
}

func (s Set) String() string {
	return fmt.Sprint(s.data)
}
