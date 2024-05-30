package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muhwyndhamhp/marknotes/ssh/base"
)

type Articles struct{}

// MatchKeyAction implements base.Page.
func (a *Articles) MatchKeyAction(m base.Model, key string, sc base.ScreenMetadata) (base.Model, bool, tea.Cmd) {
	if key != a.GetAccessKey() {
		return m, false, nil
	}

	m.Content = a.RenderPage(m.Style, sc)

	return m, true, nil
}

// GetAccessKey implements base.Page.
func (a *Articles) GetAccessKey() string {
	return "a"
}

// GetName implements base.Page.
func (a *Articles) GetName() string {
	return "Articles"
}

// RenderPage implements base.Page.
func (a *Articles) RenderPage(style lipgloss.Style, sm base.ScreenMetadata) string {
	return "This is Articles page"
}

func NewArticles() base.Page {
	return &Articles{}
}
