package base

import "github.com/charmbracelet/lipgloss"

const (
	width       = 96
	height      = 35
	columnWidth = 30

	Width       = width
	Height      = height
	ColumnWidth = columnWidth
)

var (
	// General.
	subtle     = lipgloss.AdaptiveColor{Light: "#2F373D", Dark: "#9fb9d0"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#B387FA"}
	special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#FF865B"}
	background = lipgloss.AdaptiveColor{Light: "", Dark: "#191724"}

	Subtle     = subtle
	Highlight  = highlight
	Special    = special
	Background = background

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	// Tabs.

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	footerBorder = lipgloss.Border{
		Top: "─",
	}

	footer = lipgloss.NewStyle().
		Border(footerBorder).
		BorderForeground(highlight)

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		Foreground(subtle).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Foreground(lipgloss.NoColor{}).Border(activeTabBorder, true)

	postItemBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	PostItem = lipgloss.
			NewStyle().
			Border(postItemBorder, true).
			BorderForeground(highlight).
			Padding(0, 1)

	PostTitle = lipgloss.
			NewStyle()

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Page
	docStyle = lipgloss.NewStyle().Padding(0, 2, 0, 2)

	DescStyle = lipgloss.NewStyle().Foreground(subtle)

	SubduedDescStyle = DescStyle.Foreground(special)
)
