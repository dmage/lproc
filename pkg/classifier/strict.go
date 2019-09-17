package classifier

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type strict struct {
	filename string
	classes  map[string][]string
}

func NewStrictFromFile(filename string) (Classifier, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &strict{
		filename: filename,
		classes:  make(map[string][]string),
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
				c.Learn(line, account)
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

func (c *strict) Learn(payee, class string) {
	payee = strings.Trim(payee, " \t")
	for _, cls := range c.classes[payee] {
		if cls == class {
			return
		}
	}
	c.classes[payee] = append(c.classes[payee], class)
}

func (c *strict) Classify(payee string) []string {
	payee = strings.Trim(payee, " \t")
	classes := c.classes[payee]
	if len(classes) == 0 {
		classes = nil
		bestMatch := ""
		for p, cls := range c.classes {
			if len(p) > len(bestMatch) && strings.Contains(payee, p) {
				bestMatch = p
				classes = cls
			}
		}
		if len(classes) == 1 {
			c.Learn(payee, classes[0])
		}
	}
	return classes
}

func (c *strict) Close() error {
	classToPayees := make(map[string][]string)
	for payee, classes := range c.classes {
		for _, class := range classes {
			classToPayees[class] = append(classToPayees[class], payee)
		}
	}
	var classes []string
	for class := range classToPayees {
		classes = append(classes, class)
	}

	f, err := ioutil.TempFile("", "lproc-rules-*")
	if err != nil {
		return err
	}
	sort.Strings(classes)
	for _, class := range classes {
		fmt.Fprintf(f, "%s\n", class)
		payees := classToPayees[class]
		sort.Strings(payees)
		for _, payee := range payees {
			fmt.Fprintf(f, "\t%s\n", payee)
		}
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), c.filename); err != nil {
		return err
	}

	return nil
}
