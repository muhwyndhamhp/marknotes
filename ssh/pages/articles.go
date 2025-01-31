package pages

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muhwyndhamhp/marknotes/internal"
	"github.com/muhwyndhamhp/marknotes/ssh/base"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"strings"
)

type Articles struct {
	App   *internal.Application
	Posts []internal.Post
	Page  int
}

// MatchKeyAction implements base.Page.
func (a *Articles) MatchKeyAction(m base.Model, key string, sc base.ScreenMetadata) (base.Model, bool, tea.Cmd) {
	for i := range a.Posts {
		if key == fmt.Sprintf("%d", i+1) {
			a := NewArticle(&a.Posts[i])
			m.Content = a.RenderPage(m.Style, sc)

			m.ActiveTab = -1
			m.Page = a
			return m, true, nil
		}
	}

	if key == "q" {
		if a.Page > 1 {
			a.Page--
		}

		m.Content = a.RenderPage(m.Style, sc)
		m.Page = a
		return m, true, nil
	}

	if key == "e" {
		if len(a.Posts) == 9 {
			a.Page++
		}

		m.Content = a.RenderPage(m.Style, sc)
		m.Page = a
		return m, true, nil
	}

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
	doc := strings.Builder{}

	curStyle := base.DescStyle.Padding(0, 1)
	centerStyle := curStyle.Align(lipgloss.Center).Width(sm.Width - 30)

	out := lipgloss.JoinHorizontal(lipgloss.Center, curStyle.Render("< Prev [Q]"), centerStyle.Render(fmt.Sprintf("Page %d", a.Page)), curStyle.Render("Next [E] >"))
	doc.WriteString(out + "\n\n")

	s := []scopes.QueryScope{
		scopes.OrderBy("published_at", scopes.Descending),
		scopes.Paginate(a.Page, 9),
		scopes.Where("status = ?", internal.PostStatusPublished),
		scopes.PostIndexedOnly(),
	}

	posts, err := a.App.PostRepository.Get(context.Background(), s...)
	if err != nil {
		panic(errs.Wrap(err))
	}

	a.Posts = posts

	for i, post := range posts {
		st := base.SubduedDescStyle.PaddingTop(1).Width(sm.Width - 2)
		body := lipgloss.JoinVertical(
			lipgloss.Top,
			base.PostTitle.Render(fmt.Sprintf("%s %s", post.Title, base.DescStyle.Render(fmt.Sprintf("[%d]", i+1)))),
			st.Render(
				fmt.Sprintf(
					"PostStatusPublished: %s | Updated: %s",
					post.PublishedAt.Format("Jan, 02 2006"),
					post.UpdatedAt.Format("Jan, 02 2006"),
				),
			),
		)
		doc.WriteString(base.PostItem.Width(sm.Width - 8).Render(body))
		doc.WriteString("\n")
	}

	doc.WriteString("\n\n" + out)

	return lipgloss.NewStyle().Padding(0, 4).Render(doc.String())
}

func NewArticles(app *internal.Application) base.Page {
	return &Articles{App: app, Page: 1}
}
