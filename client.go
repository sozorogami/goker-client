package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sozorogami/goker"

	ui "github.com/gizak/termui"
)

var game *goker.GameState

var promptView *PromptView
var console *ConsoleView
var playersView *PlayersView
var boardView *BoardView
var potView *PotView

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	boardView = NewBoardView()
	promptView = NewPromptView()
	console = NewConsoleView()
	potView = NewPotView()

	promptView.Render()

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
			inputString = inputString + "<!" + newInput + "!>"
		}

		promptView.SetText(inputString + "_")

		promptView.Render()
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleInput(inputString)
		inputString = ""
		promptView.SetText("_")
		promptView.SetHeading(promptString(currentPrompt, game))
		if game != nil {
			draw(*game)
		}
		promptView.Render()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Loop()
}

func parseEvents(events []interface{}) {
	for _, event := range events {
		switch event.(type) {
		case goker.DrawEvent:
			draw := event.(goker.DrawEvent)
			boardView.AnimateAppendCards(draw.Cards, time.Second, func() {})
		}
	}
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
	case actionPrompt:
		new, err := parseAction(s, *game)
		if err != nil {
			console.AppendLog(fmt.Sprintf("%s", err))
		} else {
			game = &new
			parseEvents(game.Events)
		}
	}
}

func parseAction(input string, state goker.GameState) (goker.GameState, error) {
	player := state.Action
	value := 0
	var actionType goker.ActionType
	var err error

	switch {
	case input == "C":
		actionType = goker.CheckCall
	case input == "F":
		actionType = goker.Fold
	case strings.HasPrefix(input, "B") || strings.HasPrefix(input, "R"):
		actionType = goker.BetRaise

		amtStr := strings.TrimPrefix(input, "B")
		amtStr = strings.TrimPrefix(amtStr, "R")
		value, err = strconv.Atoi(amtStr)
		if err != nil {
			return state, errors.New("Bad numeric value")
		}
	default:
		return state, errors.New("Unable to parse string")
	}

	state, err = goker.Transition(state, goker.Action{Player: player, ActionType: actionType, Value: value})

	if err != nil {
		return state, errors.New("Advance failed: " + err.Error())
	}

	return state, nil
}

func draw(game goker.GameState) {
	ui.Clear()
	ui.Body.Rows = []*ui.Row{}

	playersView = NewPlayersView(playerDataFromGameState(game))
	playersView.Render()

	belowPlayers := playersView.Height()

	boardView.SetCards(game.Board)
	boardView.SetY(belowPlayers)
	boardView.Render()

	potView.SetY(belowPlayers)
	potView.SetPots(game.Pots)
	potView.Render()

	promptView.SetY(belowPlayers + boardView.Height() + 1)
	promptView.Render()

	console.SetHeight(promptView.GetY() + promptView.Height())
	console.Render()
}

const (
	pbHeight = 6
	pbWidth  = 25
)
