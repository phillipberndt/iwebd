package util

import (
	"errors"
	"strings"
)

// A cobra.Value implementation accepting and parsing a val1:val2 combination
type ColonSeparated struct {
	value1 string
	value2 string
}

func (v *ColonSeparated) String() string {
	if v.value1 == "" && v.value2 == "" {
		return ""
	}
	return v.value1 + ":" + v.value2
}

func (v *ColonSeparated) Set(arg string) error {
	s := strings.Split(arg, ":")
	if len(s) != 2 || len(s[0]) == 0 || len(s[1]) == 0 {
		return errors.New("Syntax is " + v.Type())
	}
	v.value1 = s[0]
	v.value2 = s[1]
	return nil
}

func (v *ColonSeparated) Type() string {
	return "value1:value2"
}

func (v *ColonSeparated) IsSet() bool {
	return v.value1 != "" || v.value2 != ""
}

// An argument for user:password pairs.
type UserPass struct {
	ColonSeparated
}

func (v *UserPass) Type() string {
	return "user:password"
}

func (v *UserPass) User() string {
	return v.value1
}

func (v *UserPass) Password() string {
	return v.value2
}
