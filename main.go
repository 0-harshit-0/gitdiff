package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("running `git diff` & parsing...")

	var files []File

	if len(os.Args) > 1 {
		files = ProcessDiff(strings.Contains(os.Args[1], "staged"))
	} else {
		files = ProcessDiff(false)
	}
	GUI(files)
}
