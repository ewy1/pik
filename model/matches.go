package model

type Matches interface {
	Matches(input string) bool
}
