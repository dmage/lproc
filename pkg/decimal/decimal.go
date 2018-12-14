package decimal

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

type Decimal interface {
	String() string
}

type decimal struct {
	n     *big.Int
	scale int
}

var (
	reFloat = regexp.MustCompile(`^[+-]?[0-9]*\.[0-9]+$`)
	reInt   = regexp.MustCompile(`^[+-][0-9]+$`)
)

func New(x string) (Decimal, error) {
	if reFloat.MatchString(x) {
		idx := strings.IndexByte(x, '.')
		d := &decimal{
			n:     big.NewInt(0),
			scale: len(x) - idx - 1,
		}
		if _, ok := d.n.SetString(x[0:idx]+x[idx+1:], 10); !ok {
			return nil, fmt.Errorf("cannot parse %q as a decimal number", x)
		}
		return d, nil
	}
	return nil, fmt.Errorf("cannot parse %q as a decimal number", x)
}

func (d *decimal) String() string {
	s := d.n.String()
	if d.scale == 0 {
		return s
	}
	if d.scale >= len(s) {
		return "0." + strings.Repeat("0", d.scale-len(s)) + s
	}
	idx := len(s) - d.scale
	return s[:idx] + "." + s[idx:]
}
