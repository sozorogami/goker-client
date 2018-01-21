package main

import (
	"time"

	ui "github.com/gizak/termui"
	"github.com/sozorogami/goker"
)

type BoardView struct {
	view  *ui.Par
	cards goker.CardSet
}

func NewBoardView() *BoardView {
	box := ui.NewPar("")
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
	bv := BoardView{box, goker.CardSet{}}
	return &bv
}

func (bv *BoardView) SetY(y int) {
	bv.view.Y = y
}

func (bv *BoardView) Height() int {
	return bv.view.Height
}

func (bv *BoardView) SetCards(cards goker.CardSet) {
	bv.cards = cards
	bv.view.Text = cardsStringForCards(cards)
}

func (bv *BoardView) AnimateAppendCards(cards goker.CardSet, delay time.Duration, onEach func()) {
	for _, card := range cards {
		bv.SetCards(append(bv.cards, card))
		bv.Render()
		onEach()
		time.Sleep(delay)
	}
}

func (bv *BoardView) Render() {
	ui.Render(bv.view)
}
