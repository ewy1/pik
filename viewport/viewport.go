package viewport

import (
	"strings"
)

const Caret = "⏵"

func NeedsViewport(input string, height int) bool {
	lines := strings.Split(input, "\n")
	return len(lines)-1 > height
}

func Process(input string, height int) string {
	lines := strings.Split(input, "\n")
	if len(lines) > height {
		cropped, top, bottom := Crop(input, lines, height)
		return WithScroll(cropped, int(top*float32(height)), int(bottom*float32(height)))
	}
	return input
}

func Focus(lines []string, needle string) int {
	for i, l := range lines {
		if strings.Contains(l, needle) {
			return i
		}
	}
	return -1
}

func Crop(input string, lines []string, height int) (output string, scrollStart float32, scrollEnd float32) {
	output = input
	selectionIndex := Focus(lines, Caret)
	size := len(lines)
	if size <= height {
		return output, 0, 1
	}

	linesAbove := height / 2
	linesBelow := height - linesAbove
	if linesAbove*2 < selectionIndex {
		linesBelow++
	}

	start := selectionIndex - linesAbove
	end := selectionIndex + linesBelow

	if start < 0 {
		end += -start
		start = 0
	}

	if end >= size {
		diff := size - 1 - end
		start += diff
		end += diff
	}

	scrollStart = float32(start) / float32(size)
	scrollEnd = float32(end)/float32(size) + float32(1)/float32(size)
	if scrollEnd > 1 {
		scrollEnd = 1
	}

	return strings.Join(lines[start:end], "\n"), scrollStart, scrollEnd
}
