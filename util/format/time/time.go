package ftime

import "time"

// PrintTime formats time.Now() with a given layout
func PrintTime(layout string) string {
	t := time.Now()
	return t.Local().Format(layout)
}
