package game

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tr3yWay996/HC_Adventure/player"
)

// ── Game States ────────────────────────────────────────────────────────

type state int

const (
	stateMenu state = iota
	stateGame
	stateQuit
)

// ── Menu Options ───────────────────────────────────────────────────────

var menuItems = []string{"New Game", "Continue", "Quit"}

// ── ASCII Title ────────────────────────────────────────────────────────

const asciiTitle = `
 ██╗  ██╗ ██████╗     █████╗ ██████╗ ██╗   ██╗███████╗███╗   ██╗████████╗██╗   ██╗██████╗ ███████╗
 ██║  ██║██╔════╝    ██╔══██╗██╔══██╗██║   ██║██╔════╝████╗  ██║╚══██╔══╝██║   ██║██╔══██╗██╔════╝
 ███████║██║         ███████║██║  ██║██║   ██║█████╗  ██╔██╗ ██║   ██║   ██║   ██║██████╔╝█████╗
 ██╔══██║██║         ██╔══██║██║  ██║╚██╗ ██╔╝██╔══╝  ██║╚██╗██║   ██║   ██║   ██║██╔══██╗██╔══╝
 ██║  ██║╚██████╗    ██║  ██║██████╔╝ ╚████╔╝ ███████╗██║ ╚████║   ██║   ╚██████╔╝██║  ██║███████╗
 ╚═╝  ╚═╝ ╚═════╝    ╚═╝  ╚═╝╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝`

// ── Model ──────────────────────────────────────────────────────────────

// Model is the Bubble Tea model for the game.
type Model struct {
	state    state
	cursor   int
	width    int
	height   int
	player   *player.Player
	quitting bool
}

// NewModel creates a game model for the given player and terminal size.
func NewModel(p *player.Player, width, height int) Model {
	return Model{
		state:  stateMenu,
		cursor: 0,
		width:  width,
		height: height,
		player: p,
	}
}

// ── Bubble Tea Interface ───────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Global quit
		if key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case stateMenu:
			return m.updateMenu(key)
		case stateGame:
			return m.updateGame(key)
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	var s string

	switch m.state {
	case stateMenu:
		s = m.viewMenu()
	case stateGame:
		s = m.viewGame()
	case stateQuit:
		s = m.viewQuit()
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

// ── Menu ───────────────────────────────────────────────────────────────

func (m Model) updateMenu(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(menuItems)-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0: // New Game
			m.state = stateGame
			m.cursor = 0
			return m, tea.ClearScreen
		case 1: // Continue (placeholder)
			m.state = stateGame
			m.cursor = 0
			return m, tea.ClearScreen
		case 2: // Quit
			m.quitting = true
			return m, tea.Quit
		}
	case "q":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) viewMenu() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(asciiTitle))
	b.WriteString("\n\n")

	// Subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true).
		Render(fmt.Sprintf("Welcome, %s", m.player.Name))
	b.WriteString(subtitle)
	b.WriteString("\n\n")

	// Menu items
	for i, item := range menuItems {
		cursor := "  "
		style := MenuItemStyle
		if i == m.cursor {
			cursor = MenuCursorStyle.Render("▸ ")
			style = MenuSelectedStyle
		}
		b.WriteString(cursor + style.Render(item) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/↓ navigate • enter select • q quit"))

	// Center everything
	content := b.String()
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// ── Game ───────────────────────────────────────────────────────────────

type Choice struct {
	Text         string
	NextID       string
	GiveItem     string
	RequiredItem string
}

type Room struct {
	ID          string
	Title       string
	Description string
	Choices     []Choice
}

// Hardcoded chambers for now
var chambers = map[string]Room{
	"start": {
		ID:    "start",
		Title: "The beginning",
		Description: "You find yourself in a dimly lit bedroom.\n\n" +
			"You are curently sitting on a bed, you can see a door to your upper right, a chest with a lock close to the left wall where your bed is and a window to your right." +
			"\n\n" +
			"You start to recognize where you are, it is your first childhood bedroom.\n" +
			"There is this weird, lingering feeling, that something is off about it.\n\n" +
			"It is as if the room you always knew in your childhood had suddenly changed and became noticibly unfamiliar to you.\n",
		Choices: []Choice{
			{Text: "You don't move and observe", NextID: "observe"},
			{Text: "You get up and try to open the door", NextID: "door"},
			{Text: "You get up and try to open the chest", NextID: "chest"},
			{Text: "You get up and try to look out the window", NextID: "window"},
			{Text: "You get up and look under the bed", NextID: "under_bed"},
		},
	},
	"observe": {
		ID:          "observe",
		Title:       "Observation",
		Description: "You look around but nothing changes.",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
	"door": {
		ID:          "door",
		Title:       "Bedroom door, locked from the outside",
		Description: "You grab the handle, but the door is locked tight.",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
	"chest": {
		ID:          "chest",
		Title:       "Solid gold-ornamented chest",
		Description: "The chest is locked with a heavy padlock. You need a key.",
		Choices: []Choice{
			{Text: "Unlock the chest", NextID: "chest_open", RequiredItem: "Ornate Key"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"chest_open": {
		ID:          "chest_open",
		Title:       "Open Chest",
		Description: "The ornate key fits perfectly. You turn it and the heavy padlock clicks open! Inside you find a mysterious glowing orb.",
		Choices: []Choice{
			{Text: "Take the orb and go back", NextID: "start", GiveItem: "Glowing Orb"},
			{Text: "Leave it and go back", NextID: "start"},
		},
	},
	"window": {
		ID:          "window",
		Title:       "The Window",
		Description: "You look outside and see nothing but a swirling gray mist. This is not the outside world you remember.",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
	"under_bed": {
		ID:          "under_bed",
		Title:       "Under the Bed",
		Description: "You look under the bed and find a small, ornate key. It doesn't look like it could fit the chest's lock though.",
		Choices: []Choice{
			{Text: "Take the key and go back", NextID: "take_key", GiveItem: "Ornate Key"},
			{Text: "Leave the key and go back", NextID: "start"},
		},
	},
	"take_key": {
		ID:          "take_key",
		Title:       "Key",
		Description: "You take the key and put it in your pocket. It feels cold to the touch.",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
}

func (m Model) updateGame(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "q", "escape":
		m.state = stateMenu
		m.cursor = 0
		return m, tea.ClearScreen
	}

	// For the game state, the cursor controls the current choices
	currentRoomID, _ := m.player.Progress["current_room"].(string)
	if currentRoomID == "" {
		currentRoomID = "start"
		m.player.Progress["current_room"] = currentRoomID
	}
	room := chambers[currentRoomID]

	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(room.Choices)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor >= 0 && m.cursor < len(room.Choices) {
			choice := room.Choices[m.cursor]

			// Check for required item
			if choice.RequiredItem != "" {
				hasItem := false
				for _, item := range m.player.Inventory {
					if item == choice.RequiredItem {
						hasItem = true
						break
					}
				}
				if !hasItem {
					return m, nil // Block progression if item is missing
				}
			}

			m.player.Progress["current_room"] = choice.NextID

			// Handle item pickup
			if choice.GiveItem != "" {
				hasItem := false
				for _, item := range m.player.Inventory {
					if item == choice.GiveItem {
						hasItem = true
						break
					}
				}
				if !hasItem {
					m.player.Inventory = append(m.player.Inventory, choice.GiveItem)
				}
			}
			m.cursor = 0 // Reset cursor for the next room
			return m, tea.ClearScreen // Wipe old content before drawing new room
		}
	}

	return m, nil
}

func (m Model) viewGame() string {
	var b strings.Builder

	// Player info bar
	invStr := strings.Join(m.player.Inventory, ", ")
	if invStr == "" {
		invStr = "Empty"
	}
	info := PlayerInfoStyle.Render(fmt.Sprintf("⚔ %s  │  Session: %s  │  Inventory: %s", m.player.Name, m.player.SessionID[:8], invStr))
	b.WriteString(info)
	b.WriteString("\n\n")

	// Get current room
	currentRoomID, _ := m.player.Progress["current_room"].(string)
	if currentRoomID == "" {
		currentRoomID = "start"
	}
	room := chambers[currentRoomID]

	roomTitle := RoomTitleStyle.Render(room.Title)
	roomDesc := RoomDescStyle.Render(room.Description)

	roomBox := RoomBoxStyle.Width(min(72, m.width-4)).Render(roomTitle + "\n\n" + roomDesc)
	b.WriteString(roomBox)
	b.WriteString("\n")

	// Prompt section
	// Notice we use a Left-aligned style just for this prompt string block
	var choicesStr strings.Builder
	choicesStr.WriteString(PromptStyle.Render("What do you do?"))
	choicesStr.WriteString("\n\n")

	for i, choice := range room.Choices {
		cursor := "  "
		style := MenuItemStyle
		if i == m.cursor {
			cursor = MenuCursorStyle.Render("▸ ")
			style = MenuSelectedStyle
		}

		displayText := choice.Text
		if choice.RequiredItem != "" {
			hasItem := false
			for _, item := range m.player.Inventory {
				if item == choice.RequiredItem {
					hasItem = true
					break
				}
			}
			if hasItem {
				displayText += fmt.Sprintf(" (Use %s)", choice.RequiredItem)
			} else {
				displayText += fmt.Sprintf(" (Requires %s)", choice.RequiredItem)
				style = style.Foreground(ColorDim) // Dim the text if locked
			}
		}

		choicesStr.WriteString(cursor + style.Render(displayText) + "\n")
	}

	choicesStr.WriteString("\n")
	choicesStr.WriteString(HelpStyle.Render("↑/↓ navigate • enter select • esc/q back to menu"))

	// Create a left-aligned container for the choices
	choicesAligned := lipgloss.NewStyle().Align(lipgloss.Left).Render(choicesStr.String())
	b.WriteString(choicesAligned)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		b.String(),
	)
}

// ── Quit ───────────────────────────────────────────────────────────────

func (m Model) viewQuit() string {
	msg := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true).
		Render("Thanks for playing! See you next time, " + m.player.Name + " 👋")

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		msg,
	)
}
