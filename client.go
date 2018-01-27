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

type promptType int

const (
	numberOfPlayersPrompt promptType = iota
	playerNamePrompt
	startingChipsPrompt
	actionPrompt
)

var currentPrompt promptType
var numberOfPlayers int
var currentPlayerToName int
var playerNames []string
var startingChips int
var game *goker.GameState

var prompt *ui.Par
var console *ConsoleView
var playersView *PlayersView
var boardView *BoardView

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	boardView = NewBoardView()

	prompt = ui.NewPar("_")
	prompt.TextFgColor = ui.ColorWhite
	prompt.BorderLabel = promptString(currentPrompt, game)
	prompt.BorderFg = ui.ColorCyan
	prompt.Height = 3
	prompt.Width = 40
	prompt.X = 0
	prompt.Y = 0

	console = NewConsoleView()

	ui.Render(prompt)

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

		prompt.Text = inputString + "_"

		ui.Render(prompt)
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleInput(inputString)
		inputString = ""
		prompt.Text = "_"
		prompt.BorderLabel = promptString(currentPrompt, game)
		if game != nil {
			draw(*game)
		}
		ui.Render(prompt)
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

	pot := NewPotView()
	pot.SetY(belowPlayers)
	pot.SetPots(game.Pots)
	pot.Render()

	prompt.Y = belowPlayers + boardView.Height() + 1
	prompt.Width = pbWidth * 2
	ui.Render(prompt)

	console.SetHeight(prompt.Y + prompt.Height)
	console.Render()
}

const (
	pbHeight = 6
	pbWidth  = 25
)

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
		players := make([]*goker.Player, numberOfPlayers)
		for i, name := range playerNames {
			players[i] = goker.NewPlayer(name)
			players[i].Chips = startingChips
		}
		goker.SeatPlayers(players)

		rules := goker.GameRules{SmallBlind: 25, BigBlind: 50}
		game = goker.NewGame(players, rules, goker.NewDeck())
		currentPrompt++
	}
}

func promptString(p promptType, state *goker.GameState) string {
	switch p {
	case numberOfPlayersPrompt:
		return "How many players?"
	case playerNamePrompt:
		return fmt.Sprintf("What is player %d's name?", currentPlayerToName+1)
	case startingChipsPrompt:
		return "How many chips to start?"
	case actionPrompt:
		return actionPromptForGameState(state)
	default:
		return "Whaa?"
	}
}

func actionPromptForGameState(state *goker.GameState) string {
	var possibleActions string
	if state.BetToMatch > 0 {
		possibleActions = "(C)all, (R)aise or (F)old?"
	} else {
		possibleActions = "(C)heck, (B)et or (F)old?"
	}
	return state.Action.Name + ": " + possibleActions
}
