package game

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
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
			// Reset player state for a fresh start
			m.player.Progress = make(map[string]any)
			m.player.Inventory = make([]string, 0)
			m.player.GameVariables = make([]string, 0)
			m.state = stateGame
			m.cursor = 0
			log.Info("Player started a new game", "player", m.player.Name, "session", m.player.SessionID[:8])
			return m, tea.ClearScreen
		case 1: // Continue (placeholder)
			m.state = stateGame
			m.cursor = 0
			log.Info("Player continued game", "player", m.player.Name, "session", m.player.SessionID[:8])
			return m, tea.ClearScreen
		case 2: // Quit
			m.quitting = true
			log.Info("Player quit from menu", "player", m.player.Name, "session", m.player.SessionID[:8])
			return m, tea.Quit
		}
	case "q":
		m.quitting = true
		log.Info("Player quit via 'q'", "player", m.player.Name, "session", m.player.SessionID[:8])
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

	// Menu items — Width(12) applied here only so choices in viewGame are unaffected
	menuItemW := MenuItemStyle.Width(12)
	menuSelectedW := MenuSelectedStyle.Width(12)
	for i, item := range menuItems {
		cursor := "  "
		style := menuItemW
		if i == m.cursor {
			cursor = MenuCursorStyle.Render("▸ ")
			style = menuSelectedW
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
	Text          string
	NextID        string
	GiveItem      string
	RequiredItem  string
	GameVariable  string
	IfVariable    string
	IfNotVariable string
	AddCounter    string // Name of the counter to increment
	ReqCounter    string // Name of the counter to check
	ReqCountMin   int    // Choice appears if counter is >= this value
	ReqCountMax   int    // Choice appears if counter is <= this value (0 ignores the max limit)
}

type Room struct {
	ID          string
	Title       string
	Description string
	Choices     []Choice
}

var chambers = map[string]Room{
	"start": {
		ID:    "start",
		Title: "The beginning",
		Description: "You find yourself in a dimly lit bedroom.\n\n" +
			"You are currently sitting on a bed, you can see a door to your upper right, a chest with a lock close to the left wall where your bed is and a window to your right." +
			"\n\n" +
			"You start to recognize where you are, it is your first childhood bedroom.\n" +
			"There is this weird, lingering feeling, that something is off about it.\n\n" +
			"It is as if the room you always knew in your childhood had suddenly changed and became noticeably unfamiliar to you.\n",
		Choices: []Choice{
			{Text: "You don't move and observe", NextID: "observe"},
			{Text: "You get up and try to open the door", NextID: "door", IfNotVariable: "door_broken"},
			{Text: "You get up and look at the broken door", NextID: "door_broken", IfVariable: "door_broken"},
			{Text: "You get up and try to open the chest", NextID: "first-room-chest"},
			{Text: "You get up and try to look out the window", NextID: "window"},
			{Text: "You get up and look under the bed", NextID: "under_bed"},
		},
	},

	// Observing dialog
	"observe": {
		ID:          "observe",
		Title:       "You decided to sit and observe around",
		Description: "You sit on your bed and look around. Nothing happen.",
		Choices: []Choice{
			{Text: "Try harder", NextID: "observing_longer"},
		},
	},
	"observing_longer": {
		ID:          "observing",
		Title:       "You're still observing around.",
		Description: "Half an hour ago you decided to stay on the bed for longer, staring at the walls of the room aimlessly, zoning out, doing nothing productive as the time passes",
		Choices: []Choice{
			{Text: "Lay down", NextID: "observing_troll_loop"},
		},
	},
	"observing_troll_loop": {
		ID:          "observing_troll_loop",
		Title:       "Laying down",
		Description: "You decide to lay down on the bed and stare at the ceiling instead, slowly drifting to sleep",
		Choices: []Choice{
			{Text: "Relax", NextID: "observing_troll_loop", AddCounter: "relax_count", ReqCounter: "relax_count", ReqCountMax: 2},
			{Text: "Snap out of it and get up", NextID: "start", ReqCounter: "relax_count", ReqCountMin: 3},
		},
	},
	// Bedroom door dialog
	"door": {
		ID:          "door",
		Title:       "Bedroom door, locked from the outside",
		Description: "You grab the handle, but the door is locked tight.",
		Choices: []Choice{
			{Text: "Try to open the door", NextID: "door_try"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"door_try": {
		ID:          "door_try",
		Title:       "Bedroom door",
		Description: "You try to open the door, but it is locked tight.",
		Choices: []Choice{
			{Text: "Try to pull on the handle", NextID: "door_try_pull"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"door_try_pull": {
		ID:          "door_try_pull",
		Title:       "Bedroom door",
		Description: "You notice light on the other end when pulling hard on the handle. Maybe something will happen if you keep trying ?",
		Choices: []Choice{
			{Text: "Pull harder", NextID: "door_try_pull_harder"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"door_try_pull_harder": {
		ID:          "door_try_pull_harder",
		Title:       "Bedroom door",
		Description: "You pull harder on the handle, you feel the door creaking and the wood around the handle splintering. You can almost open it.",
		Choices: []Choice{
			{Text: "Pull even harder", NextID: "door_try_pull_harder_harder"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"door_try_pull_harder_harder": {
		ID:          "door_try_pull_harder_harder",
		Title:       "Bedroom door",
		Description: "You pull even harder on the handle, you hear the wood of the door cracking.",
		Choices: []Choice{
			{Text: "Pull with both hands with all your strength", NextID: "door_broken"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"door_broken": {
		ID:          "door_broken",
		Title:       "Bedroom door",
		Description: "You look at what remains of the door handle, what once was your only way of getting out of this corrupted room, now ruined by your own brutality..\n\n What a mess",
		Choices: []Choice{
			{Text: "Go back", NextID: "start", GameVariable: "door_broken"},
		},
	},

	// Chest dialog
	"first-room-chest": {
		ID:          "first-room-chest",
		Title:       "Heavy gold-ornamented chest",
		Description: "The chest is locked with a heavy padlock. You need a key.",
		Choices: []Choice{
			{Text: "Unlock the chest", NextID: "chest_try", RequiredItem: "Ornate Key"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"chest_try": {
		ID:          "chest_try",
		Title:       "Try to open the chest",
		Description: "The ornate key feels loose in the lock. You turn it, you can hear the metal of the key rattling inside but the heavy padlock doesn't budge, this is not the right key.",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
	"chest_open": {
		ID:          "chest_open",
		Title:       "The bedroom chest",
		Description: "You put the golden key inside and turn it. The heavy padlock opens with a satisfying *click*. \n\nYou open the chest and you find an old book with a wax sealed letter with the letters H.C written with a quill in red",
		Choices: []Choice{
			{Text: "Take the old book", NextID: "take_old_book", GiveItem: "Old book"},
			{Text: "Take the wax sealed letter", NextID: "take_wax_sealed_letter", GiveItem: "Wax sealed letter"},
			{Text: "Leave the key and go back", NextID: "start"},
		},
	},

	// Chest items dialogs
	"old_book": {
		ID:          "old_book",
		Title:       "An old book",
		Description: "The old book you found in the gold chest earlier, who know what it may contain",
		Choices: []Choice{
			{Text: "Read the book", NextID: "read_book"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"sealed_letter": {
		ID:          "sealed_letter",
		Title:       "A wax sealed letter",
		Description: "The letter that you found earlier in the golden chest, you can see a wax seal with the letters H.C written in red by a quill on the face of it",
		Choices: []Choice{
			{Text: "Read the letter", NextID: "read_letter"},
			{Text: "Go back", NextID: "start"},
		},
	},

	// Reading the book dialogs
	"read_book": {
		ID:          "read_book",
		Title:       "The book",
		Description: "The old book worn by time you managed to get out of the golden chest in one piece",
		Choices: []Choice{
			{Text: "Read the old book"},
			{Text: "Don't read the old book", NextID: "start"},
		},
	},

	// Window dialog
	"window": {
		ID:          "window",
		Title:       "The Window",
		Description: "You try to look out the window, but its opaque finishing prevents you from seeing anything outside.",
		Choices: []Choice{
			{Text: "Try to open the window", NextID: "window_try"},
			{Text: "Go back", NextID: "start"},
		},
	},
	"window_try": {
		ID:          "window_try",
		Title:       "The window",
		Description: "You try to open the window, you grab the handle and try to twist it. It doesn't budge, the only thing it does is cracking noises from the wood around the window frame. Maybe you could force it open with something ?",
		Choices: []Choice{
			{Text: "Go back", NextID: "start"},
		},
	},
	"under_bed": {
		ID:          "under_bed",
		Title:       "Under the Bed",
		Description: "You look under the bed and find a small, ornate key. Unfortunate that it doesn't look like it could fit the chest's lock though.",
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

func (m Model) hasVariable(variable string) bool {
	for _, v := range m.player.GameVariables {
		if v == variable {
			return true
		}
	}
	return false
}

func (m Model) getActiveChoices(room Room) []Choice {
	var active []Choice
	for _, c := range room.Choices {
		if c.IfVariable != "" && !m.hasVariable(c.IfVariable) {
			continue
		}
		if c.IfNotVariable != "" && m.hasVariable(c.IfNotVariable) {
			continue
		}

		if c.ReqCounter != "" {
			count := 0
			if val, ok := m.player.Progress[c.ReqCounter].(int); ok {
				count = val
			}
			if c.ReqCountMin > 0 && count < c.ReqCountMin {
				continue
			}
			if c.ReqCountMax > 0 && count > c.ReqCountMax {
				continue
			}
		}
		active = append(active, c)
	}
	return active
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

	activeChoices := m.getActiveChoices(room)

	switch key {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(activeChoices)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor >= 0 && m.cursor < len(activeChoices) {
			choice := activeChoices[m.cursor]

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

			log.Debug("Player actions/info debug:", "player", m.player.Name, "session", m.player.SessionID[:8], "room", currentRoomID, "choice", choice.Text, "next_room", choice.NextID)
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
			// handle game variables
			if choice.GameVariable != "" {
				if !m.hasVariable(choice.GameVariable) {
					m.player.GameVariables = append(m.player.GameVariables, choice.GameVariable)
				}
			}

			// handle counters
			if choice.AddCounter != "" {
				count := 0
				if val, ok := m.player.Progress[choice.AddCounter].(int); ok {
					count = val
				}
				m.player.Progress[choice.AddCounter] = count + 1
			}

			m.cursor = 0              // Reset cursor for the next room
			return m, tea.ClearScreen // Wipe old content before drawing new dialog box
		}
	}

	return m, nil
}

func (m Model) viewGame() string {
	var b strings.Builder

	// Get current room early to find active loop counters
	currentRoomID, _ := m.player.Progress["current_room"].(string)
	if currentRoomID == "" {
		currentRoomID = "start"
	}
	room := chambers[currentRoomID]

	// Find relevant loop counters for this room
	var loopCounters []string
	seenCounters := make(map[string]bool)
	for _, c := range room.Choices {
		counterName := c.ReqCounter
		if counterName == "" {
			counterName = c.AddCounter
		}
		if counterName != "" && !seenCounters[counterName] {
			count := 0
			if val, ok := m.player.Progress[counterName].(int); ok {
				count = val
			}
			loopCounters = append(loopCounters, fmt.Sprintf("%s: %d", counterName, count))
			seenCounters[counterName] = true
		}
	}
	loopStr := "-"
	if len(loopCounters) > 0 {
		loopStr = strings.Join(loopCounters, ", ")
	}

	// Player info bar
	invStr := strings.Join(m.player.Inventory, ", ")
	if invStr == "" {
		invStr = "Empty"
	}
	info := PlayerInfoStyle.Render(fmt.Sprintf("⚔ %s  │  Session: %s  │  Inventory: %s  │  Active variables: %v  │  Loop counters: %s", m.player.Name, m.player.SessionID[:8], invStr, m.player.GameVariables, loopStr))
	b.WriteString(info)
	b.WriteString("\n\n")

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

	activeChoices := m.getActiveChoices(room)

	for i, choice := range activeChoices {
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
	choicesStr.WriteString(HelpStyle.Render("↑/↓ j/k navigate • enter select • q back to menu"))

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
