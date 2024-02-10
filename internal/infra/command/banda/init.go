package banda_command

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Banda Lumaksa server.",
	Long: `Setting up basic configuration file to make Banda Lumaksa server works.
Including database setup, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompts := []Prompt{
			{
				Prompt: "Please, input postgresql database host to store the banda lumaksa data!",
				Answer: "",
				Type:   InputDBHost,
				Hidden: false,
			},
			{
				Prompt: "Now, input postgresql database port!",
				Answer: "",
				Type:   InputDBPort,
				Hidden: false,
			},
			{
				Prompt: "Input postgresql user!",
				Answer: "",
				Type:   InputDBUser,
				Hidden: false,
			},
			{
				Prompt: "Input postgresql password!",
				Answer: "",
				Type:   InputDBPassword,
				Hidden: true,
			},
			{
				Prompt: "Input database name!",
				Answer: "",
				Type:   InputDBName,
				Hidden: true,
			},
			{
				Prompt: "Input postgresql ssl mode!",
				Answer: "",
				Type:   InputDBSSLMode,
				Hidden: true,
			},
		}

		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		p := tea.NewProgram(NewModel(prompts))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

type Type int

const (
	InputDBHost Type = iota
	InputDBPort
	InputDBUser
	InputDBPassword
	InputDBName
	InputDBSSLMode
)

type Prompt struct {
	Prompt string
	Answer string
	Type   Type
	Hidden bool
}

func NewPrompt(t Type, prompt string, hidden bool) Prompt {
	return Prompt{
		Prompt: prompt,
		Answer: "",
		Type:   t,
		Hidden: hidden,
	}
}

type Model struct {
	Input   textinput.Model
	Prompts []Prompt
	State   Type
}

func (m Model) Prompt() Prompt {
	return m.Prompts[m.State]
}

func (Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			current := m.Prompts[m.State]
			current.Answer = m.Input.Value()
			m.Prompts[m.State] = current

			log.Printf("%s, %s", m.State, InputDBHost)
			log.Printf("%s, %s", m.State, InputDBSSLMode)

			if m.State == InputDBSSLMode {
				return m, tea.Quit
			}

			m.State += 1
			next := m.Prompt()
			m.Input.SetValue("")
			if next.Hidden {
				m.Input.EchoMode = textinput.EchoPassword
			} else {
				m.Input.EchoMode = textinput.EchoNormal
			}

			return m, cmd
		}
	}

	m.Input, cmd = m.Input.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n",
		m.Prompts[m.State].Prompt,
		m.Input.View(),
	)
}

func NewModel(prompts []Prompt) Model {
	ti := textinput.New()
	ti.Focus()

	return Model{
		Input:   ti,
		Prompts: prompts,
		State:   InputDBHost,
	}
}
