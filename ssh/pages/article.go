package pages

import (
	"fmt"
	"github.com/muhwyndhamhp/marknotes/internal"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muhwyndhamhp/marknotes/config"
	"github.com/muhwyndhamhp/marknotes/ssh/base"
)

type Article struct {
	Post *internal.Post
}

// GetAccessKey implements base.Page.
func (a *Article) GetAccessKey() string {
	return "p"
}

// GetName implements base.Page.
func (a *Article) GetName() string {
	return "article"
}

// MatchKeyAction implements base.Page.
func (a *Article) MatchKeyAction(m base.Model, key string, sc base.ScreenMetadata) (base.Model, bool, tea.Cmd) {
	return m, false, nil
}

// RenderPage implements base.Page.
func (a *Article) RenderPage(style lipgloss.Style, screenMeta base.ScreenMetadata) string {
	doc := strings.Builder{}
	doc.WriteString(titleStyle.Padding(0, 1, 0, 6).Width(screenMeta.Width-4).Render(a.Post.Title) + "\n\n")

	tags := strings.Split(a.Post.TagsLiteral, ",")

	str := ""
	for i := range tags {
		if tags[i] == "" {
			continue
		}

		str += fmt.Sprintf("#%s ", tags[i])
	}

	doc.WriteString(hashtagStyle.Padding(0, 1, 0, 6).Width(screenMeta.Width-4).Render(str) + "\n")

	md, err := os.ReadFile(config.Get(config.POST_RENDER_PATH) + "/markdowns/" + a.Post.Slug + ".md")
	if err != nil {
		return doc.String()
	}
	out, err := glamour.Render(string(md), "tokyo-night")
	if err != nil {
		return doc.String()
	}

	doc.WriteString(base.DescStyle.Padding(0, 1, 0, 4).Render(out) + "\n")

	return doc.String()
}

func NewArticle(post *internal.Post) base.Page {
	return &Article{post}
}

var (
	titleStyle = lipgloss.
			NewStyle().
			Foreground(base.Highlight).
			AlignHorizontal(lipgloss.Center)

	hashtagStyle = lipgloss.NewStyle().Foreground(base.Special).AlignHorizontal(lipgloss.Center)
)
