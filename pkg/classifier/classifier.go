package classifier

import (
	"fmt"
)

type Config struct {
	Name string
	File string
}

type Classifier interface {
	Classify(in string) []string
	Close() error
}

type Factory interface {
	GetClassifier(name string) (Classifier, error)
	Close() error
}

type factory struct {
	config []Config
	cache  map[string]Classifier
}

func NewFactory(cfg []Config) Factory {
	return &factory{
		config: cfg,
		cache:  make(map[string]Classifier),
	}
}

func (f *factory) GetClassifier(name string) (Classifier, error) {
	if c, ok := f.cache[name]; ok {
		return c, nil
	}
	for _, cfg := range f.config {
		if cfg.Name != name {
			continue
		}
		c, err := NewStrictFromFile(cfg.File)
		if err != nil {
			return c, err
		}
		f.cache[name] = c
		return c, nil
	}
	return nil, fmt.Errorf("unknown classifier: %s", name)
}

func (f *factory) Close() error {
	for name, c := range f.cache {
		if err := c.Close(); err != nil {
			return fmt.Errorf("close classifier %s: %s", name, err)
		}
	}
	return nil
}
