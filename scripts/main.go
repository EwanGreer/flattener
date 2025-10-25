package main

import (
	"fmt"
	"os"

	"github.com/EwanGreer/flattener"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: flattener <json-string> [delimiter]")
		fmt.Fprintln(os.Stderr, "Example: flattener '{\"user\":{\"name\":\"john\"}}' '.'")
		os.Exit(1)
	}

	input := os.Args[1]
	delimiter := "."
	if len(os.Args) > 2 {
		delimiter = os.Args[2]
	}

	f := flattener.Flattener{
		Delimeter: delimiter,
	}

	result, err := f.JSON([]byte(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(result))
}
