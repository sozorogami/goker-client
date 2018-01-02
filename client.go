package main

import (
	"fmt"
	"regexp"
	"strconv"

	ui "github.com/gizak/termui"
)

type promptType int

const (
	numberOfPlayersPrompt promptType = iota
	playerNamePrompt
	startingChipsPrompt
)

var currentPrompt promptType
var numberOfPlayers int
var currentPlayerToName int
var playerNames []string
var startingChips int

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar("_")
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = promptString(currentPrompt)
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
			if l := len(inputString); l > 0 {
				inputString = inputString[:l-1]
			}
		} else if isAlphanumeric(newInput) {
			inputString = inputString + newInput
		} else {
			inputString = inputString + "!" + newInput + "!"
		}

		p.Text = inputString + "_"

		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleInput(inputString)
		inputString = ""
		p.Text = "_"
		p.BorderLabel = promptString(currentPrompt)
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		p.BorderLabel = "Meow"
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}

func handleInput(s string) {
	if s == "exit" {
		ui.StopLoop()
	}

	switch currentPrompt {
	case numberOfPlayersPrompt:
		getNumberOfPlayers(s)
	case playerNamePrompt:
		getPlayerName(s)
	case startingChipsPrompt:
		getStartingChips(s)
	}

}

func getNumberOfPlayers(s string) {
	val, err := strconv.Atoi(s)
	if err != nil || val < 2 || val > 10 {
		numberOfPlayers = -1
	} else {
		numberOfPlayers = val
		playerNames = make([]string, val)
		currentPrompt++
	}
}

func getPlayerName(s string) {
	l := len(s)
	if l > 0 && l < 20 {
		playerNames[currentPlayerToName] = s
		currentPlayerToName++
	}
	if currentPlayerToName == numberOfPlayers {
		currentPrompt++
	}
}

func getStartingChips(s string) {
	val, err := strconv.Atoi(s)
	if err != nil || val <= 0 {
		startingChips = -1
	} else {
		startingChips = val
		currentPrompt++
	}
}

func promptString(p promptType) string {
	switch p {
	case numberOfPlayersPrompt:
		return "How many players?"
	case playerNamePrompt:
		return fmt.Sprintf("What is player %d's name?", currentPlayerToName+1)
	case startingChipsPrompt:
		return "How many chips to start?"
	default:
		return "Whaa?"
	}
}
