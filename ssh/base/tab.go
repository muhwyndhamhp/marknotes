package base

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Tab struct {
	Title       string
	ShortAction string
	Page        Page
}

func (t Tab) ConstructTitle() string {
	return fmt.Sprintf("%s [%s]", t.Title, t.ShortAction)
}

func (m Model) RenderTabs(doc *strings.Builder) {
	var tabs []string
	for i := range m.Tabs {
		if m.ActiveTab == i {
			tabs = append(tabs, activeTab.Render(m.Tabs[i].ConstructTitle()))
		} else {
			tabs = append(tabs, tab.Render(m.Tabs[i].ConstructTitle()))
		}
	}
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabs...,
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-4)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n\n")
}

func (m Model) MatchTab(key string) (Model, bool) {
	for i := range m.Tabs {
		if m.Tabs[i].ShortAction == key {
			m.ActiveTab = i
			m.Page = m.Tabs[i].Page
			return m, true
		}
	}

	return m, false
}
