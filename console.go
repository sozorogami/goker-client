package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui"
)

type ConsoleView struct {
	view *ui.Par
	logs []string
}

func NewConsoleView() *ConsoleView {
	events := ui.NewPar("")
	events.TextFgColor = ui.ColorGreen
	events.BorderLabel = "Events"
	events.X = pbWidth * 2
	events.Y = 0
	events.Width = 40
	cv := ConsoleView{events, []string{}}
	return &cv
}

func (cv *ConsoleView) SetHeight(h int) {
	cv.view.Height = h
}

func (cv *ConsoleView) AppendLog(log string) {
	cv.logs = append(cv.logs, log)
	fmt.Println(len(cv.logs))
}

func (cv *ConsoleView) Render() {
	cv.view.Text = strings.Join(cv.linesToDisplay(), "\n")
	ui.Render(cv.view)
}

func (cv ConsoleView) linesToDisplay() []string {
	displayHeight := cv.view.Height - 2

	var count int
	var idx int
	for i := len(cv.logs) - 1; i >= 0 && count < displayHeight; i-- {
		count += cv.lineCountWhenRendered(cv.logs[i])
		idx = i
	}

	lines := make([]string, 0, count)
	for j := idx; j < len(cv.logs); j++ {
		sublines := cv.breakString(cv.logs[j])
		lines = append(lines, sublines...)
	}

	if count < displayHeight {
		return lines
	}

	return lines[len(lines)-displayHeight:]
}

func (cv ConsoleView) breakString(s string) []string {
	l := cv.lineCountWhenRendered(s)
	w := cv.displayWidth()
	result := make([]string, l)
	for i := 0; i < l-1; i++ {
		result[i] = s[w*i : w*(i+1)]
	}
	result[l-1] = s[w*(l-1):]
	return result
}

func (cv ConsoleView) lineCountWhenRendered(s string) int {
	return len(s)/cv.displayWidth() + 1
}

func (cv ConsoleView) displayWidth() int {
	return cv.view.Width - 2
}
