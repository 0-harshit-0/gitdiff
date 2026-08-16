package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("running `git diff` & parsing...")

	files := ProcessDiff(strings.Contains(os.Args[1], "staged"))
	GUI(files)
}
