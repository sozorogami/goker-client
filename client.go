package main

import (
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

		ui.Render(ui.Body)
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
		new := parseAction(s, *game)
		game = &new
	}
}

func parseAction(input string, state goker.GameState) goker.GameState {
	player := state.Action
	value := 0
	var actionType goker.ActionType
	var err error

	switch {
	case input == "C":
		actionType = goker.CheckCall
	case input == "F":
		actionType = goker.Fold
	case strings.HasPrefix(input, "B"):
		actionType = goker.BetRaise

		amtStr := strings.TrimPrefix(input, "B")
		value, err = strconv.Atoi(amtStr)
		if err != nil {
			panic("Bad numeric value")
		}
	default:
		panic("Unable to parse string")
	}

	state, err = goker.Transition(state, goker.Action{Player: player, ActionType: actionType, Value: value})

	if err != nil {
		panic("Advance failed: " + err.Error())
	}

	return state
}

func draw(game goker.GameState) {
	ui.Clear()
	ui.Body.Rows = []*ui.Row{}

	playerRowCount := (len(game.Players) + 1) / 2
	lastRowSingleColumn := len(game.Players)%2 != 0
	for i := 0; i < playerRowCount; i++ {
		var cols = []*ui.Row{}
		if i == playerRowCount-1 && lastRowSingleColumn {
			player := game.Players[2*i]
			cols = append(cols, ui.NewCol(4, 0, playerBox(playerDataForPlayer(player, game))))
		} else {
			player1 := game.Players[2*i]
			player2 := game.Players[2*i+1]
			cols = append(cols,
				ui.NewCol(4, 0, playerBox(playerDataForPlayer(player1, game))),
				ui.NewCol(4, 0, playerBox(playerDataForPlayer(player2, game))),
			)
		}
		if i == 0 {
			cols = append(cols, ui.NewCol(4, 0, events))
		}
		ui.Body.AddRows(ui.NewRow(cols...))
		ui.Body.AddRows(ui.NewRow(ui.NewCol(8, 0, prompt)))
	}
	ui.Body.Align()
	ui.Render(ui.Body)
}

type playerData struct {
	name, status          string
	chipCount, currentBet int
	isActive, isDealer    bool
	hand                  string
}

func playerDataForPlayer(player *goker.Player, state goker.GameState) playerData {
	var statusString string

	switch player.Status {
	case goker.Active:
		statusString = "Active"
	case goker.AllIn:
		statusString = "All In"
	case goker.Folded:
		statusString = "Folded"
	case goker.Eliminated:
		statusString = "Eliminated"
	}

	cardStrings := []string{}
	for _, card := range player.HoleCards {
		cardStrings = append(cardStrings, card.String())
	}
	handString := strings.Join(cardStrings, " ")

	data := playerData{
		name:       player.Name,
		status:     statusString,
		chipCount:  player.Chips,
		currentBet: player.CurrentBet,
		isActive:   state.Action == player,
		isDealer:   state.Dealer == player,
		hand:       handString,
	}

	return data
}

func playerBox(data playerData) *ui.Par {
	p := ui.NewPar(playerInfoString(data))
	p.TextFgColor = ui.ColorWhite

	if data.isDealer {
		p.BorderLabel = data.name + " (Dealer)"
	} else {
		p.BorderLabel = data.name
	}

	if data.isActive {
		p.BorderFg = ui.ColorWhite
	} else {
		p.BorderFg = ui.ColorBlue
	}
	p.Height = 6
	return p
}

func playerInfoString(data playerData) string {
	return fmt.Sprintf("Status: %s\nChips: %d\nBet: %d\nHand: %s", data.status, data.chipCount, data.currentBet, data.hand)
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
	return "It is " + state.Action.Name + "'s turn. " + possibleActions
}
