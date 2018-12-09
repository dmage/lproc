package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/dmage/lproc/pkg/classifier"
	"github.com/dmage/lproc/pkg/rewriter"
)

type Format struct {
	Name       string
	Encoding   string
	Comma      string
	HeaderSkip int
	Reverse    bool
	Columns    []string
	Rewrite    []rewriter.Rule
}

type Rules struct {
	Classifiers []classifier.Config
	Formats     []Format
}

func LoadRulesFromFile(filename string) (*Rules, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rules := &Rules{}
	err = yaml.NewDecoder(f).Decode(rules)
	if err != nil {
		return nil, fmt.Errorf("load rules from %s: %s", filename, err)
	}
	return rules, nil
}

func LoadDefaultRules() (*Rules, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	return LoadRulesFromFile(filepath.Join(usr.HomeDir, ".lproc", "rules.yaml"))
}

func (r Rules) GetFormat(name string) (Format, error) {
	for _, format := range r.Formats {
		if format.Name == name {
			return format, nil
		}
	}
	return Format{}, fmt.Errorf("no configuration for the format: %s", name)
}
