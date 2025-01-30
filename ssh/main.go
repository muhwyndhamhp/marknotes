package main

import (
	"context"
	"errors"
	_postRepo "github.com/muhwyndhamhp/marknotes/internal/post"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apsystole/log"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muhwyndhamhp/marknotes/db"
	"github.com/muhwyndhamhp/marknotes/ssh/base"
	"github.com/muhwyndhamhp/marknotes/ssh/pages"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(runTea),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server ", "host ", host, " port ", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func runTea(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	postRepo := _postRepo.NewPostRepository(db.GetLibSQLDB())

	homePage := pages.NewHome(postRepo)
	articlesPage := pages.NewArticles()

	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	}
	return base.Model{
		ActiveTab: 0,
		Tabs:      getTabs(homePage, articlesPage),
		Page:      homePage,
		Style:     newStyle(),
		Viewport:  viewport.New(base.Width-4, base.Height-4),
	}, opts
}

func getTabs(pages ...base.Page) []base.Tab {
	var tabs []base.Tab

	for i := range pages {
		tabs = append(tabs, base.Tab{
			Title:       pages[i].GetName(),
			ShortAction: pages[i].GetAccessKey(),
			Page:        pages[i],
		})
	}
	return tabs
}

func newStyle() lipgloss.Style {
	return lipgloss.
		NewStyle()
}
