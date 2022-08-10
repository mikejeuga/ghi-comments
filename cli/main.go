package main

import (
	"fmt"
	client "github.com/mikejeuga/ghi-comments"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string { return i.title }

func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	config := client.NewConfig()
	ghClient := client.NewGHClient(config)
	issues, err := ghClient.GetIssues()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}
	for _, issue := range issues {
		items = append(items, issue)
	}

	delegate := list.NewDefaultDelegate()
	delegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		issue, ok := m.SelectedItem().(client.Issue)
		if ok {
			return m.NewStatusMessage(issue.Titles)

		}
		return nil
	}
	m := model{list: list.New(items, delegate, 0, 0)}
	m.list.Title = "Mike's Repo Issues"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
