package identity

import "strings"

type Identity struct {
	Full    string
	Reduced string
}

func (i Identity) Is(input string) bool {
	reduced := Reduce(input)
	return i.Reduced == reduced
}

func New(input string) Identity {
	reduced := Reduce(input)
	return Identity{
		Full:    input,
		Reduced: reduced,
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
