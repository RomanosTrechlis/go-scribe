package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RomanosTrechlis/go-scribe/writer"
)

func main() {
	l, err := writer.NewLogger("writer", "romanos", 8080, "", "", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating logger writer: %v", err)
		os.Exit(1)
	}
	logger := log.New(l, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("this is a test")

	logger.Printf("this is another test because %s", "reasons")

	lll, _ := writer.NewLogger("newWriter", "romanos", 8080, "", "", "")
	lll.Println("new logger")
}
