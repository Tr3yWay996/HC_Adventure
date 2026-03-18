package player

import (
	"time"

	"github.com/charmbracelet/ssh"
)

// Player holds per-session data for a connected player.
// Extend this struct to persist progress, inventory, stats, etc.
type Player struct {
	Name        string
	SessionID   string
	ConnectedAt time.Time
	Progress    map[string]any
	Inventory   []string
}

// New creates a Player from the incoming SSH session.
func New(s ssh.Session) *Player {
	name := s.User()
	if name == "" {
		name = "adventurer"
	}
	return &Player{
		Name:        name,
		SessionID:   s.Context().SessionID(),
		ConnectedAt: time.Now(),
		Progress:    make(map[string]any),
		Inventory:   make([]string, 0),
	}
}
