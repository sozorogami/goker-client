package main

import (
	"regexp"

	ui "github.com/gizak/termui"
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar("_")
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan
	p.Height = 3

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, p),
		),
	)

	ui.Body.Align()

	ui.Render(ui.Body)

	var inputString string

	ui.Handle("/sys/kbd", func(a ui.Event) {
		isAlphanumeric := regexp.MustCompile(`^[A-Za-z0-9]$`).MatchString

		kbdEvt := a.Data.(ui.EvtKbd)
		newInput := kbdEvt.KeyStr

		if newInput == "<space>" {
			inputString = inputString + " "
		} else if newInput == "C-8" {
			l := len(inputString)
			inputString = inputString[:l-1]
		} else if isAlphanumeric(newInput) {
			inputString = inputString + newInput
		} else {
			inputString = inputString + "!" + newInput + "!"
		}

		p.Text = inputString + "_"

		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		p.BorderLabel = "Ello"
		ui.Body.Align()
		ui.Render(ui.Body)
		ui.StopLoop()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		p.BorderLabel = "Meow"
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
