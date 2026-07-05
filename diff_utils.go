package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func diffOutput() []byte {
	cmd := exec.Command("git", "diff")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return stdout
}

type File struct {
	Name string
	Data []string
}

func ProcessDiff() []File {
	stdout := string(diffOutput())

	filelREG := regexp.MustCompile(`(?m)^diff --git `)
	chunks := filelREG.Split(stdout, -1)

	files := []File{}

	for _, chunk := range chunks {
		if chunk == "" {
			continue
		}

		hunkREG := regexp.MustCompile(`(?m)^@@ .*? @@`)
		loc := hunkREG.FindStringIndex(chunk)

		var filename string
		var lines []string

		if loc != nil {
			filename = chunk[:loc[0]]
			header := chunk[loc[0]:loc[1]]
			rest := chunk[loc[1]:]
			if !strings.HasPrefix(rest, "\n") && !strings.HasPrefix(rest, "\r\n") {
				// Git usually adds a single space before the context code.
				// Trim that space and inject a newline so the context becomes its own line.
				rest = "\n" + strings.TrimPrefix(rest, " ")
			}

			// Recombine and split
			lines = strings.Split(header+rest, "\n")
		} else {
			filename = "none"
			lines = []string{}
		}

		file := File{
			Name: filename,
			Data: lines,
		}
		files = append(files, file)
	}

	return files
}
