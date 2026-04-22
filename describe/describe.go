package describe

import (
	"bufio"
	"io"
	"os"
	"pik/model"
	"strings"
)

var DescriptionPrefixes = []string{
	"#",
	"//",
}

var descriptions = make(map[model.Target]*string)

func Describe(key model.Target, file string) (string, error) {
	if d := descriptions[key]; d != nil {
		return *d, nil
	}
	fd, err := os.Open(file)
	if err != nil {
		msg := err.Error()
		descriptions[key] = &msg
		return "", err
	}
	defer fd.Close()
	text, err := FromReader(fd)
	if err != nil {
		return text, err
	} else {
		descriptions[key] = &text
	}
	return text, err
}

func FromReader(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	text := scanner.Text()
	if strings.HasPrefix(text, "#!") {
		scanner.Scan()
		text = scanner.Text()
	}
	text = strings.TrimSpace(text)
	hasPrefix := false
	for _, p := range DescriptionPrefixes {
		if strings.HasPrefix(text, p) {
			hasPrefix = true
			break
		}
	}
	if !hasPrefix {
		return "", nil
	}
	for _, c := range DescriptionPrefixes {
		text = strings.TrimPrefix(text, c)
		text = strings.TrimSpace(text)
	}
	return text, nil
}
