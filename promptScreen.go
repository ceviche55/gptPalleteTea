package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type paletteViewStruct struct {
}

func (pvs paletteViewStruct) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return pvs, tea.Quit
		case "enter":
			if pvs.prompt == "" {
				pvs.prompt = pvs.textInput.Value()

				pvs.palette = aiPaletteCall(pvs.apiKey, pvs.prompt)

				pvs.state = paletteView
			}
			return pvs, nil
		}
	}

	pvs.textInput, cmd = pvs.textInput.Update(msg)

	return pvs, cmd
}
