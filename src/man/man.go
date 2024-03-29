package man

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

const _maxResults = 36

// ManInterface is the interface for working with the man database
type ManInterface struct {
	Commands []*Command
}

type Command struct {
	Name        string
	Description string
	ManArg      string
	ManURI      string
}

// NewManInterface creates a new interface for working with the man database
func NewManInterface() *ManInterface {
	return &ManInterface{}
}

// GetManDatabase retrieves a list of the man database on the local machine
func (m *ManInterface) GetManDatabase() error {
	cmd := exec.Command("/usr/bin/man", "-k", "-P", "cat", ".")

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err != nil {
		return err
	}

	var re = regexp.MustCompile(`(?im)^([^a-z0-9\s]+)$`)

	for _, item := range strings.Split(outb.String(), "\n") {
		pages := strings.Split(item, " - ")
		if len(pages) != 2 {
			continue
		}

		for _, name := range strings.Split(pages[0], ",") {
			command := &Command{
				Name:        removeNonAscii(strings.TrimSpace(name)),
				Description: removeNonAscii(re.ReplaceAllString(pages[1], "")),
				ManURI:      manURI(removeNonAscii(strings.TrimSpace(name))),
				ManArg:      manArg(removeNonAscii(strings.TrimSpace(name))),
			}
			m.Commands = append(m.Commands, command)
		}
	}

	return nil
}

// removeNonAscii removes all non-ascii characters from a string
func removeNonAscii(data string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		return r
	}, data)
}

// ManURI generates a man database URI
func manURI(data string) string {
	return fmt.Sprintf("x-man-page://%s", manArg(data))
}

// ManArg generates a man database URI
func manArg(data string) string {
	var re = regexp.MustCompile(`(?m)(.*)\((.+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)
	return fmt.Sprintf("%s/%s", matches[0][2], matches[0][1])
}
