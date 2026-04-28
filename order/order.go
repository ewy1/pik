package order

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"pik/describe"
	"pik/identity"
	"strings"
)

type Element struct {
	Identifier  identity.Identity
	Description string
}

type Order struct {
	Elements []Element
}

var Empty = Order{}

func FromFile(f fs.FS, path string) (Order, error) {
	fd, err := os.Open(path)
	if err != nil {
		return Empty, err
	}
	defer fd.Close()
	return FromReader(fd)
}

func FromReader(r io.Reader) (Order, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for _, p := range describe.DescriptionPrefixes {
			if strings.HasPrefix(line, p) {
				continue
			}
		}
	}

}
