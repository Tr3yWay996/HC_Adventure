package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/Tr3yWay996/HC_Adventure/game"
	"github.com/Tr3yWay996/HC_Adventure/player"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/ssh"
)

const (
	host = "127.0.0.1"
	port = "6666"
)

func main() {
	log.SetLevel(log.DebugLevel)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not create server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Server error", "error", err)
			done <- nil
		}
	}()

	<-done

	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server gracefully", "error", err)
	}
}

// teaHandler creates a new Bubble Tea program for each SSH session.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	p := player.New(s)
	m := game.NewModel(p, pty.Window.Width, pty.Window.Height)

	// Force TrueColor for each SSH session.
	// On Windows, bubbletea auto-detects the color profile from the server
	// process environment, which doesn't have COLORTERM or similar set,
	// so it defaults to no-color. We override it here to ensure SSH clients
	// receive full ANSI colors regardless of the server OS.
	return m, []tea.ProgramOption{
		tea.WithColorProfile(colorprofile.TrueColor),
	}
}
