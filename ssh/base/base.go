package base

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

type Model struct {
	ActiveTab int
	Tabs      []Tab
	Page      Page
	Content   string
	Width     int
	Height    int
	Viewport  viewport.Model
	Style     lipgloss.Style
}

type Page interface {
	RenderPage(style lipgloss.Style, screenMeta ScreenMetadata) string
	GetName() string
	GetAccessKey() string
	MatchKeyAction(m Model, key string, sc ScreenMetadata) (Model, bool, tea.Cmd)
}

type ScreenMetadata struct {
	Width int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	doc := strings.Builder{}
	m.RenderTabs(&doc)
	return lipgloss.
		NewStyle().
		Width(tern.Int(m.Width, width)).
		Height(tern.Int(m.Height, height)).
		Align(lipgloss.Center).
		Render(
			fmt.Sprintf("\n\n\n\n\n%s%s\n\n%s",
				doc.String(),
				m.Viewport.View(),
				footer.Width(width-4).AlignVertical(lipgloss.Bottom).Render(""),
			))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	screenMeta := ScreenMetadata{
		Width: width - 4,
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.Content == "" {
		m.Content = m.Page.RenderPage(m.Style, screenMeta)
		m.Viewport.SetContent(m.Content)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		m, _ = m.MatchTab(string(msg.String()))
		var c tea.Cmd
		var ok bool
		m, ok, c = m.Page.MatchKeyAction(m, string(msg.String()), screenMeta)
		cmds = append(cmds, c)

		if ok {
			m.Viewport.SetContent(m.Content)
		}

	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

		m.Viewport.Width = tern.Min(m.Width, width) - 4
		m.Viewport.Height = tern.Min(m.Height, height) - 4
		m.Viewport.YPosition = 4
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
