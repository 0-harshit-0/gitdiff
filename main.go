package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("git", "diff")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(stdout))
}
