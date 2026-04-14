package identity

import "strings"

type Identity struct {
	Full    string
	Reduced string
}

func (i Identity) Is(input string) bool {
	return Reduce(input) == i.Reduced
}

func New(input string) Identity {
	return Identity{
		Full:    input,
		Reduced: Reduce(input),
	}

}

func Reduce(input string) string {
	reduced := input
	if !strings.HasPrefix(reduced, ".") {
		reduced = strings.Split(reduced, ".")[0]
	}
	reduced = strings.ToLower(reduced)
	return reduced

}
