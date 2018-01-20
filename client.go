package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
var events *ui.Par

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	prompt = ui.NewPar("_")
	prompt.TextFgColor = ui.ColorWhite
	prompt.BorderLabel = promptString(currentPrompt, game)
	prompt.BorderFg = ui.ColorCyan
	prompt.Height = 3

	events = ui.NewPar("")
	events.TextFgColor = ui.ColorGreen
	events.BorderLabel = "Events"
	events.Height = 12

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, events),
		),
		ui.NewRow(
			ui.NewCol(8, 0, prompt),
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
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
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
	case actionPrompt:
		new, err := parseAction(s, *game)
		if err != nil {
			events.Text = events.Text + fmt.Sprintf("%s", err) + "\n"
		}
		game = &new
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
		amtStr = strings.TrimPrefix(input, "R")
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

	boxes := []*ui.Par{}

	playerCount := len(game.Players)

	for i := 0; i < playerCount; i++ {
		row := i / 2
		col := i % 2
		box := playerBox(playerDataForPlayer(game.Players[i], game), row, col)
		boxes = append(boxes, box)
	}
	for _, box := range boxes {
		ui.Render(box)
	}

	belowPlayers := (playerCount + 1) / 2 * pbHeight

	board := boardBox(game.Board, belowPlayers)
	ui.Render(board)

	var chips int
	if len(game.Pots) > 0 {
		chips = game.Pots[0].Value
	}

	pot := potBox(chips, belowPlayers)
	ui.Render(pot)

	prompt.Y = belowPlayers + board.Height + 1
	prompt.Width = pbWidth * 2
	ui.Render(prompt)

	events.X = pbWidth * 2
	events.Y = 0
	events.Height = prompt.Y + prompt.Height
	ui.Render(events)
}

type playerData struct {
	name                  string
	status                goker.PlayerStatus
	chipCount, currentBet int
	isActive, isDealer    bool
	hand                  string
}

func stringForPlayerStatus(status goker.PlayerStatus) string {
	switch status {
	case goker.Active:
		return "Active"
	case goker.AllIn:
		return "All In"
	case goker.Folded:
		return "Folded"
	case goker.Eliminated:
		return "Eliminated"
	default:
		return ""
	}
}

func playerDataForPlayer(player *goker.Player, state goker.GameState) playerData {
	cardStrings := []string{}
	for _, card := range player.HoleCards {
		cardStrings = append(cardStrings, card.String())
	}
	handString := strings.Join(cardStrings, " ")

	data := playerData{
		name:       player.Name,
		status:     player.Status,
		chipCount:  player.Chips,
		currentBet: player.CurrentBet,
		isActive:   state.Action == player,
		isDealer:   state.Dealer == player,
		hand:       handString,
	}

	return data
}

const (
	pbHeight = 6
	pbWidth  = 25
)

func potBox(chips, y int) *ui.Par {
	chipString := strconv.Itoa(chips)
	padding := (pbWidth - len(chipString) - 2) / 2
	padString := strings.Repeat(" ", padding)

	content := "\n" + padString + chipString

	box := ui.NewPar(content)
	box.BorderLabel = "Pot"
	box.BorderBg = ui.ColorGreen
	box.BorderFg = ui.ColorWhite
	box.BorderLabelBg = ui.ColorGreen
	box.BorderLabelFg = ui.ColorWhite
	box.Bg = ui.ColorWhite
	box.TextBgColor = ui.ColorWhite
	box.TextFgColor = ui.ColorBlack
	box.Width = pbWidth
	box.Height = 5
	box.X = pbWidth
	box.Y = y
	return box
}

func boardBox(cards goker.CardSet, y int) *ui.Par {
	box := ui.NewPar(cardsStringForCards(cards))
	box.BorderLabel = "Board"
	box.BorderBg = ui.ColorRed
	box.BorderFg = ui.ColorWhite
	box.BorderLabelBg = ui.ColorRed
	box.BorderLabelFg = ui.ColorWhite
	box.Bg = ui.ColorWhite
	box.TextBgColor = ui.ColorWhite
	box.TextFgColor = ui.ColorBlack
	box.Width = pbWidth
	box.Height = 5
	box.X = 0
	box.Y = y
	return box
}

func cardsStringForCards(cards goker.CardSet) string {
	cardStrings := []string{}
	for _, card := range cards {
		cardStrings = append(cardStrings, card.String())
	}

	cardCount := len(cards)
	for i := 5 - cardCount; i > 0; i-- {
		cardStrings = append(cardStrings, "??")
	}

	padding := "    "

	return "\n" + padding + strings.Join(cardStrings, " ")
}

func playerBox(data playerData, row, col int) *ui.Par {
	p := ui.NewPar(playerInfoString(data))
	p.TextFgColor = ui.ColorWhite
	p.Height = pbHeight
	p.Width = pbWidth
	p.X = pbWidth * col
	p.Y = pbHeight * row

	if data.isDealer {
		p.BorderLabel = data.name + " (Dealer)"
	} else {
		p.BorderLabel = data.name
	}

	if data.isActive {
		p.BorderFg = ui.ColorWhite
	} else {
		switch data.status {
		case goker.Active:
			p.BorderFg = ui.ColorBlue
		case goker.Folded:
			p.BorderFg = ui.ColorRGB(1, 1, 1)
		case goker.Eliminated:
			p.BorderFg = ui.ColorBlack
		case goker.AllIn:
			p.BorderFg = ui.ColorYellow
		}
	}
	return p
}

func playerInfoString(data playerData) string {
	return fmt.Sprintf("Status: %s\nChips: %d\nBet: %d\nHand: %s", stringForPlayerStatus(data.status), data.chipCount, data.currentBet, data.hand)
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
