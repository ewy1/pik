package identity

import "strings"

type Identity struct {
	Full    string
	Reduced string
}

// I return whether the other identity "means the same" as this one
func (i Identity) I(other Identity) bool {
	return i.Reduced == other.Reduced
}

// Is returns whether the other string "means the same" as this one
func (i Identity) Is(input string) bool {
	reduced := Reduce(input)
	return i.Reduced == reduced
}

// New creates a new Identity with default reduction
func New(input string) Identity {
	reduced := Reduce(input)
	return Identity{
		Full:    input,
		Reduced: reduced,
	}

}

// Reduce normalizes input (commands and filenames) to simplify them
func Reduce(input string) string {
	reduced := input
	reduced = strings.TrimPrefix(input, ".")
	if !strings.HasPrefix(reduced, ".") {
		reduced = strings.Split(reduced, ".")[0]
	}
	reduced = strings.ToLower(reduced)
	return reduced

}
