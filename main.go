package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/joho/godotenv"
)

type sessionState int

const (
	promptView sessionState = iota
	paletteView
	setUpView
)

type mainModelStruct struct {
	prompt    string
	textInput textinput.Model
	state     sessionState
	palette   []string
	apiKey    string
}

func initialModel(apiKeyStr string) mainModelStruct {
	newTextInput := textinput.New()
	newTextInput.Placeholder = "I want a vibe like x, fits like y, ivokes the feelings of z"
	newTextInput.Focus()
	newTextInput.CharLimit = 150
	newTextInput.Width = 50

	return mainModelStruct{
		textInput: newTextInput,
		state:     promptView,
		apiKey:    apiKeyStr,
	}
}

func (mm mainModelStruct) Init() tea.Cmd {
	return textinput.Blink
}

func (mm mainModelStruct) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mm, tea.Quit
		}
	}

	mm.textInput, cmd = mm.textInput.Update(msg)

	return mm, cmd
}

func (mm mainModelStruct) View() string {

	var str string

	if mm.state == promptView {
		str += fmt.Sprintf(
			"What vibe do you want the palette to be?\n\n%s\n\n%s",
			mm.textInput.View(),
			"(Ctrl + C or Esc to quit)",
		) + "\n"
	} else {
		str += mm.prompt

		str += "\n\n"

		var renderedArray []string

		for i := 0; i < len(mm.palette); i++ {

			preRenderedElement := lipgloss.NewStyle().
				Background(lipgloss.Color(mm.palette[i])).
				Padding(0, 1, 6).
				Margin(0, 1).
				Border(lipgloss.NormalBorder()).
				Render(mm.palette[i])

			palleteElement := lipgloss.NewStyle().
				Padding(1).
				Margin(0, 1).
				Border(lipgloss.NormalBorder()).
				Render(mm.palette[i])

			renderedElement := lipgloss.JoinVertical(
				lipgloss.Top,
				palleteElement,
				preRenderedElement,
			)

			renderedArray = append(renderedArray, renderedElement)
		}

		str += lipgloss.JoinHorizontal(lipgloss.Center, renderedArray...)
	}

	str += "\n\n"

	return str
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to open .env file.")
	}
	apiKey := os.Getenv("OPENAI_KEY")

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Fatal: ", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(initialModel(apiKey))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
