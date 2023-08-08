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

// Structure for managing state
type sessionState int

const (
	promptView sessionState = iota
	paletteView
	setUpView
)

// Top level model/struct
type mainModelStruct struct {
	prompt    string
	textInput textinput.Model
	currState sessionState
	palette   []string
	apiKey    string
}

// Switching states with cmd and msg
type changeStateMsg struct{ newState sessionState }

func changeState(newState sessionState) tea.Cmd {
	return func() tea.Msg {
		return changeStateMsg{newState}
	}
}

func initialModel(apiKeyStr string) mainModelStruct {
	newTextInput := textinput.New()
	newTextInput.Placeholder = "I want a vibe like x, fits like y, ivokes the feelings of z"
	newTextInput.Focus()
	newTextInput.CharLimit = 150
	newTextInput.Width = 50

	return mainModelStruct{
		textInput: newTextInput,
		currState: promptView,
		apiKey:    apiKeyStr,
	}
}

func (mms mainModelStruct) Init() tea.Cmd {
	return textinput.Blink
}

func (mms mainModelStruct) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return mms, tea.Quit
		case "enter":
			if mms.currState == promptView {
				mms.prompt = mms.textInput.Value()

				mms.palette = aiPaletteCall(mms.apiKey, mms.prompt)

				// changeState(paletteView)
				mms.currState = paletteView
			}
			return mms, nil
		}
	case changeStateMsg:
		if mms.currState == msg.newState {
			log.Print("State was same")
			break
		}

		mms.currState = msg.newState
		log.Print("Reached state change")

		return mms, nil
	}

	mms.textInput, cmd = mms.textInput.Update(msg)

	return mms, cmd
}

func (mms mainModelStruct) View() string {

	var str string

	switch mms.currState {
	case promptView:
		str += fmt.Sprintf(
			"What vibe do you want the palette to be?\n\n%s\n\n%s",
			mms.textInput.View(),
			"(Ctrl + C or Esc to quit)",
		) + "\n"
	case paletteView:
		str += mms.prompt

		str += "\n\n"

		var renderedArray []string

		for i := 0; i < len(mms.palette); i++ {

			colorElement := lipgloss.NewStyle().
				Background(lipgloss.Color(mms.palette[i])).
				Padding(3, 0).
				// Margin(0, 1).
				// Border(lipgloss.NormalBorder()).
				Render("       ")

			hexCodeElement := lipgloss.NewStyle().
				// Padding(1).
				// Margin(0, 1).
				// Border(lipgloss.NormalBorder()).
				Render(mms.palette[i])

			columnElementRaw := lipgloss.JoinVertical(
				lipgloss.Top,
				hexCodeElement,
				colorElement,
			)

			columnElement := lipgloss.NewStyle().
				// Margin(0, 1).
				Border(lipgloss.NormalBorder()).
				Render(columnElementRaw)

			renderedArray = append(renderedArray, columnElement)
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
