package main

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/sozorogami/goker"
)

type playerData struct {
	name                  string
	status                goker.PlayerStatus
	chipCount, currentBet int
	isActive, isDealer    bool
	hand                  string
}

type PlayersView struct {
	views []*ui.Par
	data  []playerData
}

func NewPlayersView(numberOfPlayers int) *PlayersView {
	boxes := []*ui.Par{mainPlayerBox(numberOfPlayers)}
	boxes = append(boxes, leftPlayerBoxes(numberOfPlayers)...)
	boxes = append(boxes, rightPlayerBoxes(numberOfPlayers)...)

	pv := PlayersView{boxes, make([]playerData, numberOfPlayers)}
	return &pv
}

func (pv PlayersView) Render() {
	for _, box := range pv.views {
		ui.Render(box)
	}
}

func (pv PlayersView) Height() int {
	return (len(pv.data) + 1) / 2 * pbHeight
}

func playerDataFromGameState(state goker.GameState) []playerData {
	data := make([]playerData, len(state.Players))
	for i, player := range state.Players {
		data[i] = playerDataForPlayer(player, state)
	}
	return data
}

func mainPlayerBox(numberOfPlayers int) *ui.Par {
	y := (numberOfPlayers - 1) / 2 * pbHeight
	var x int
	if numberOfPlayers%2 == 1 {
		x = pbWidth / 2
	}
	return emptyPlayerBox(x, y)
}

func leftPlayerBoxes(numberOfPlayers int) []*ui.Par {
	count := (numberOfPlayers+1)/2 - 1
	boxes := make([]*ui.Par, count)
	for i := count - 1; i >= 0; i-- {
		y := pbHeight * (count - 1 - i)
		boxes[i] = emptyPlayerBox(0, y)
	}
	return boxes
}

func rightPlayerBoxes(numberOfPlayers int) []*ui.Par {
	count := numberOfPlayers / 2
	boxes := make([]*ui.Par, count)
	for i := 0; i < count; i++ {
		y := pbHeight * i
		x := pbWidth
		boxes[i] = emptyPlayerBox(x, y)
	}
	return boxes
}

func emptyPlayerBox(x, y int) *ui.Par {
	p := ui.NewPar("")
	p.TextFgColor = ui.ColorWhite
	p.Height = pbHeight
	p.Width = pbWidth
	p.X = x
	p.Y = y
	return p
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

func (pv PlayersView) setData(data []playerData) {
	pv.data = data
	for i, datum := range data {
		pv.setDataForBox(datum, i)
	}
}

func (pv PlayersView) setDataForBox(data playerData, boxIdx int) {
	p := pv.views[boxIdx]
	p.Text = playerInfoString(data)

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
}

func playerInfoString(data playerData) string {
	return fmt.Sprintf("Status: %s\nChips: %d\nBet: %d\nHand: %s", stringForPlayerStatus(data.status), data.chipCount, data.currentBet, data.hand)
}
