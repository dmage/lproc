package rewriter

import (
	"crypto/md5"
	"fmt"
	"io"
)

func init() {
	register("md5", funcMD5)
}

func funcMD5(state *State, args Arguments) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected exactly one argument")
	}
	val, err := args.Evaluate(0, state)
	if err != nil {
		return "", err
	}
	h := md5.New()
	io.WriteString(h, val)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
