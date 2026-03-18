package game

import "charm.land/lipgloss/v2"

// ── Color Palette ──────────────────────────────────────────────────────

var (
	ColorPrimary   = lipgloss.Color("#7C3AED") // violet
	ColorSecondary = lipgloss.Color("#06B6D4") // cyan
	ColorAccent    = lipgloss.Color("#F59E0B") // amber
	ColorMuted     = lipgloss.Color("#6B7280") // gray
	ColorSuccess   = lipgloss.Color("#10B981") // emerald
	ColorDanger    = lipgloss.Color("#EF4444") // red
	ColorText      = lipgloss.Color("#E5E7EB") // light gray
	ColorDim       = lipgloss.Color("#4B5563") // dim gray
)

// ── Title ──────────────────────────────────────────────────────────────

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary).
	MarginBottom(1)

// ── Menu ───────────────────────────────────────────────────────────────

var (
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(ColorText)

	MenuSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(ColorAccent).
				Bold(true)

	MenuCursorStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)
)

// ── Game View ──────────────────────────────────────────────────────────

var (
	RoomBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginBottom(1)

	RoomTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			MarginBottom(1)

	RoomDescStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	PromptStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	PlayerInfoStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(1)
)
