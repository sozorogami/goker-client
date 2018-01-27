package main

import (
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/sozorogami/goker"
)

type PotView struct {
	view *ui.Par
	pots []*goker.Pot
}

func NewPotView() *PotView {
	box := ui.NewPar("0")
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
	pv := PotView{box, []*goker.Pot{}}
	return &pv
}

func (pv *PotView) SetPots(pots []*goker.Pot) {
	pv.pots = pots
	pv.view.Text = stringForPots(pots)
}

func (pv *PotView) SetY(y int) {
	pv.view.Y = y
}

func (pv PotView) Render() {
	ui.Render(pv.view)
}

func stringForPots(pots []*goker.Pot) string {
	var components []string
	if len(pots) == 0 {
		components = []string{"0"}
	} else {
		components = []string{}
		for i, pot := range pots {
			potID := ""
			if len(pots) > 1 {
				potID = "[" + strconv.Itoa(i+1) + "]"
			}

			chipString := strconv.Itoa(pot.Value)
			components = append(components, potID+chipString)
		}
	}

	result := strings.Join(components, " ")
	padding := (pbWidth - len(result) - 2) / 2
	padString := strings.Repeat(" ", padding)

	return "\n" + padString + result
}
