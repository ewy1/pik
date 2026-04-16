package describe

import (
	"bufio"
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
	scanner := bufio.NewScanner(fd)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	text := scanner.Text()
	if strings.HasPrefix(text, "#!") {
		scanner.Scan()
		text = scanner.Text()
	}
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "#") {
		return "", nil
	}
	for _, c := range DescriptionPrefixes {
		text = strings.TrimPrefix(text, c)
		text = strings.TrimSpace(text)
	}
	descriptions[key] = &text
	return text, nil
}
