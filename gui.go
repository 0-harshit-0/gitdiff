package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func parseHunkHeader(line string) (int, int, bool) {
	if !strings.HasPrefix(line, "@@ ") {
		return 0, 0, false
	}
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return 0, 0, false
	}
	oldPart := strings.TrimPrefix(strings.Split(parts[1], ",")[0], "-")
	newPart := strings.TrimPrefix(strings.Split(parts[2], ",")[0], "+")

	var oldStart, newStart int
	fmt.Sscanf(oldPart, "%d", &oldStart)
	fmt.Sscanf(newPart, "%d", &newStart)

	return oldStart, newStart, true
}

func GUI(files []File) {
	a := app.New()
	w := a.NewWindow("GitDiff")

	mergeViewer := widget.NewRichText()
	mergeViewer.Wrapping = fyne.TextWrapWord
	mergeScroll := container.NewVScroll(mergeViewer)

	oldViewer := widget.NewRichText()
	oldViewer.Wrapping = fyne.TextWrapWord
	oldScroll := container.NewVScroll(oldViewer)

	newViewer := widget.NewRichText()
	newViewer.Wrapping = fyne.TextWrapWord
	newScroll := container.NewVScroll(newViewer)

	splitLayoutContainer := container.NewHSplit(oldScroll, newScroll)
	splitLayoutContainer.SetOffset(0.5)

	rightViewStack := container.NewStack()

	currentFileID := -1
	currentMode := "Merge View"

	updateView := func() {
		if currentFileID == -1 {
			msg := []widget.RichTextSegment{
				&widget.TextSegment{
					Style: widget.RichTextStyle{ColorName: theme.ColorNameForeground},
					Text:  "Select a file from the sidebar to view its diff layout.",
				},
			}
			mergeViewer.Segments = msg
			oldViewer.Segments = msg
			newViewer.Segments = msg
			mergeViewer.Refresh()
			oldViewer.Refresh()
			newViewer.Refresh()
			rightViewStack.Objects = []fyne.CanvasObject{mergeScroll}
			rightViewStack.Refresh()
			return
		}

		var mergeSegments []widget.RichTextSegment
		var oldSegments []widget.RichTextSegment
		var newSegments []widget.RichTextSegment

		oldLineNum := 1
		newLineNum := 1

		monospaceStyle := fyne.TextStyle{Monospace: true}

		mutedNumStyle := widget.RichTextStyle{
			ColorName: theme.ColorNamePlaceHolder,
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		headerStyle := widget.RichTextStyle{
			ColorName: theme.ColorNamePrimary,
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		errorStyle := widget.RichTextStyle{
			ColorName: theme.ColorNameError,
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		successStyle := widget.RichTextStyle{
			ColorName: theme.ColorNameSuccess,
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		commonCodeStyle := widget.RichTextStyle{
			ColorName: theme.ColorNameForeground,
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		defaultStyle := widget.RichTextStyle{
			TextStyle: monospaceStyle,
			Inline:    true,
		}

		lines := files[currentFileID].Data
		for len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		for i, line := range lines {
			endChar := "\n"
			if i == len(lines)-1 {
				endChar = ""
			}

			if startOld, startNew, ok := parseHunkHeader(line); ok {
				oldLineNum = startOld
				newLineNum = startNew

				// Cleaned up headers
				mergeSegments = append(mergeSegments, &widget.TextSegment{Style: headerStyle, Text: "==== │ " + line + endChar})
				oldSegments = append(oldSegments, &widget.TextSegment{Style: headerStyle, Text: "==== │ " + line + endChar})
				newSegments = append(newSegments, &widget.TextSegment{Style: headerStyle, Text: "==== │ " + line + endChar})
				continue
			}

			if strings.HasPrefix(line, "-") {
				oldStr := fmt.Sprintf("%4d", oldLineNum)
				oldLineNum++

				// Merge View: Show old line number on left, blank on right
				mergeSegments = append(mergeSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s      │ ", oldStr)},
					&widget.TextSegment{Style: errorStyle, Text: line + endChar},
				)

				oldSegments = append(oldSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s │ ", oldStr)},
					&widget.TextSegment{Style: errorStyle, Text: line + endChar},
				)

				newSegments = append(newSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: "     │ "},
					&widget.TextSegment{Style: defaultStyle, Text: endChar},
				)

			} else if strings.HasPrefix(line, "+") {
				newStr := fmt.Sprintf("%4d", newLineNum)
				newLineNum++

				// Merge View: Show blank on left, new line number on right
				mergeSegments = append(mergeSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("     %s │ ", newStr)},
					&widget.TextSegment{Style: successStyle, Text: line + endChar},
				)

				oldSegments = append(oldSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: "     │ "},
					&widget.TextSegment{Style: defaultStyle, Text: endChar},
				)

				newSegments = append(newSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s │ ", newStr)},
					&widget.TextSegment{Style: successStyle, Text: line + endChar},
				)

			} else {
				oldStr := fmt.Sprintf("%4d", oldLineNum)
				newStr := fmt.Sprintf("%4d", newLineNum)
				oldLineNum++
				newLineNum++

				// Merge View: Show both line numbers side-by-side for context lines
				mergeSegments = append(mergeSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s %s │ ", oldStr, newStr)},
					&widget.TextSegment{Style: commonCodeStyle, Text: line + endChar},
				)

				oldSegments = append(oldSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s │ ", oldStr)},
					&widget.TextSegment{Style: commonCodeStyle, Text: line + endChar},
				)

				newSegments = append(newSegments,
					&widget.TextSegment{Style: mutedNumStyle, Text: fmt.Sprintf("%s │ ", newStr)},
					&widget.TextSegment{Style: commonCodeStyle, Text: line + endChar},
				)
			}
		}

		mergeViewer.Segments = mergeSegments
		oldViewer.Segments = oldSegments
		newViewer.Segments = newSegments

		mergeViewer.Refresh()
		oldViewer.Refresh()
		newViewer.Refresh()

		if currentMode == "Merge View" {
			rightViewStack.Objects = []fyne.CanvasObject{mergeScroll}
		} else {
			rightViewStack.Objects = []fyne.CanvasObject{splitLayoutContainer}
		}
		rightViewStack.Refresh()
	}

	updateView()

	viewModeRadio := widget.NewRadioGroup([]string{"Merge View", "Split View"}, func(selected string) {
		currentMode = selected
		updateView()
	})
	viewModeRadio.Horizontal = true
	viewModeRadio.SetSelected("Merge View")

	bottomControlBar := container.NewHBox(layout.NewSpacer(), viewModeRadio)
	rightContentArea := container.NewBorder(nil, bottomControlBar, nil, nil, rightViewStack)

	sidebarBox := container.NewVBox()
	for id, file := range files {
		fileID := id
		fileName := file.Name

		formattedName := strings.ReplaceAll(fileName, "\n", "\n\n")

		rt := widget.NewRichText()
		rt.Wrapping = fyne.TextWrapWord
		rt.ParseMarkdown(formattedName)

		rowButton := widget.NewButton("", nil)
		clickableRow := container.NewMax(rowButton, container.NewPadded(rt))

		rowButton.OnTapped = func() {
			currentFileID = fileID
			updateView()
			mergeScroll.ScrollToTop()
			oldScroll.ScrollToTop()
			newScroll.ScrollToTop()
		}

		sidebarBox.Add(clickableRow)
	}

	leftScroll := container.NewVScroll(sidebarBox)
	splitLayout := container.NewHSplit(leftScroll, rightContentArea)
	splitLayout.SetOffset(0.25)

	mainLayout := container.NewBorder(nil, nil, nil, nil, splitLayout)
	w.SetContent(mainLayout)
	w.Resize(fyne.NewSize(1300, 800))
	w.ShowAndRun()
}
