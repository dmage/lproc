package classifier

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type strict struct {
	classes map[string][]string
}

func NewStrictFromFile(filename string) (Classifier, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &strict{
		classes: make(map[string][]string),
	}
	account := ""
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == nil || line != "" {
			line = strings.TrimRight(line, "\n")
			if line == "" {
				continue
			}
			if line[0] == ' ' || line[0] == '\t' {
				payee := strings.Trim(line, " \t")
				c.classes[payee] = append(c.classes[payee], account)
			} else {
				account = strings.Trim(line, " \t")
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *strict) Classify(payee string) []string {
	payee = strings.Trim(payee, " \t")
	return c.classes[payee]
}
