package model

import (
	"slices"
	"strings"
)

type Tag *string
type TagAction func(src *Source)

func New(input string) Tag {
	result := &input
	TagMap[input] = result
	TagList = append(TagList, result)
	return result
}

var (
	Here     = New("here")
	Pre      = New("pre")
	Post     = New("post")
	Final    = New("final")
	Hidden   = New("hidden")
	Single   = New("single")
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
