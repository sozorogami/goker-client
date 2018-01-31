package main

import (
	"fmt"
	"strconv"

	ui "github.com/gizak/termui"
	"github.com/sozorogami/goker"
)

type promptType int

var currentPrompt promptType
var numberOfPlayers int
var currentPlayerToName int
var playerNames []string
var startingChips int

const (
	numberOfPlayersPrompt promptType = iota
	playerNamePrompt
	startingChipsPrompt
	actionPrompt
)

type PromptView struct {
	view *ui.Par
}

func (pv PromptView) Render() {
	ui.Render(pv.view)
}

func NewPromptView() *PromptView {
	prompt := ui.NewPar("_")
	prompt.TextFgColor = ui.ColorWhite
	prompt.BorderLabel = promptString(currentPrompt, game)
	prompt.BorderFg = ui.ColorCyan
	prompt.Height = 3
	prompt.Width = pbWidth * 2
	prompt.X = 0
	prompt.Y = 0
	pv := PromptView{prompt}
	return &pv
}

func (pv PromptView) SetText(text string) {
	pv.view.Text = text
}

func (pv PromptView) SetHeading(text string) {
	pv.view.BorderLabel = text
}

func (pv PromptView) SetY(y int) {
	pv.view.Y = y
}

func (pv PromptView) GetY() int {
	return pv.view.Y
}

func (pv PromptView) Height() int {
	return pv.view.Height
}

func getNumberOfPlayers(s string) {
	val, err := strconv.Atoi(s)
	if err != nil || val < 2 || val > 10 {
		numberOfPlayers = -1
	} else {
		numberOfPlayers = val
		playerNames = make([]string, val)
		initViews()
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
