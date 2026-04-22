package model

import (
	"slices"
	"strings"
)

// Tag is some text which is contained in a filename which triggers pik functionality
type Tag *string
type TagAction func(src *Source)

// New creates a new tag and registers it in the subsystems
func New(input string) Tag {
	result := &input
	TagMap[input] = result
	TagList = append(TagList, result)
	return result
}

var (
	// Here will force the target to run in the current directory instead of the source directory
	Here = New("here")
	// Pre turns the target into a trigger, causing it to be triggered before another target gets ran
	Pre = New("pre")
	// Post turns the target into a trigger, causing it to be triggered after another target gets ran and exits succesfully
	Post = New("post")
	// Final turns the target into a trigger, causing it to be triggered after another target gets ran
	Final = New("final")
	// Hidden means the target is not visible in the menu
	Hidden = New("hidden")
	// Single means this target will not use any triggers
	Single = New("single")
	// Override means this should be selected instead of a non-override target, if possible
	Override = New("override")
)

var TagList []Tag

var TagMap = map[string]Tag{}

type Tags []Tag

func (t Tags) AnyOf(expected ...Tag) bool {
	if len(expected) > 1 && len(t) == 0 {
		return false
	}
	if len(expected) == 0 {
		return true
	}
	for _, e := range expected {
		if slices.Contains(t, e) {
			return true
		}
	}
	return false
}

func (t Tags) Has(expected Tag) bool {
	return slices.Contains(t, expected)
}

func TagsFromFilename(filename string) Tags {
	var tags Tags
	// if hidden
	if strings.HasPrefix(filename, ".") {
		filename = strings.TrimPrefix(filename, ".")
		tags = append(tags, Hidden)
	}

	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return nil
	}

	for _, p := range parts {
		p = strings.ToLower(p)
		if TagMap[p] != nil {
			tags = append(tags, TagMap[p])
		}
	}

	return tags
}

func (t Tags) Visible() bool {
	return !t.AnyOf(Hidden, Pre, Post, Final)
}
