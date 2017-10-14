package scribe

import (
	"fmt"
	"strconv"
)

const (
	b  = 1
	kb = 1000 * b
	mb = 1000 * kb
	gb = 1000 * mb
	tb = 1000 * gb
)

// LexicalToNumber converts string sizes to byte sizes
func LexicalToNumber(size string) (int64, error) {
	// this condition is necessary for including infinite file size
	if size == "-1" {
		return -1, nil
	}
	// 5MB
	l := len(size)
	m := 1.0
	switch size[l-2:] {
	case "KB":
		m = kb
	case "MB":
		m = mb
	case "GB":
		m = gb
	case "TB":
		m = tb
	default:
		return 0, fmt.Errorf("couldn't transform size to bytes")
	}
	var s float64
	s, err := strconv.ParseFloat(size[:l-2], 64)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse input to float: %v", err)
	}
	return int64(s * m), nil
}
