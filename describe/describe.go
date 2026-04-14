package describe

import (
	"bufio"
	"os"
	"strings"
)

var DescriptionPrefixes = []string{
	"#",
	"//",
}

func Describe(file string) (string, error) {
	fd, err := os.Open(file)
	if err != nil {
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
	return text, nil
}
