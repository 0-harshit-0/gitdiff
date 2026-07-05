package main

import (
	"fmt"
)

func main() {
	fmt.Println("running `git diff` & parsing...")

	files := ProcessDiff()
	GUI(files)
}
