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

func NewPlayersView(data []playerData) *PlayersView {
	boxes := []*ui.Par{}
	playerCount := len(data)

	for i := 0; i < playerCount; i++ {
		row := i / 2
		col := i % 2
		box := playerBox(data[i], row, col)
		boxes = append(boxes, box)
	}

	pv := PlayersView{boxes, data}
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
