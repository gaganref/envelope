package main

import "github.com/charmbracelet/lipgloss"

// Palette defines the application's color scheme, inspired by Tailwind CSS.
type Palette struct {
	Primary struct {
		Base  lipgloss.TerminalColor
		Light lipgloss.TerminalColor
	}
	Grey struct {
		// Using a scale where a higher number is generally darker.
		// Note: Adaptive colors can invert this in light vs. dark mode.
		C300 lipgloss.TerminalColor // Lighter grey, for help text
		C400 lipgloss.TerminalColor // Mid-grey, for list items
		C500 lipgloss.TerminalColor // Brighter grey, for questions
		C100 lipgloss.TerminalColor // Adaptive white/black for high contrast text
	}
	Special struct {
		Success lipgloss.TerminalColor
		Error   lipgloss.TerminalColor
	}
}

// NewPalette creates the default color palette for the application.
func NewPalette() Palette {
	p := Palette{}

	p.Primary.Base = lipgloss.Color("#F26C0D")
	p.Primary.Light = lipgloss.Color("#FF8A42")

	// In dark mode, a higher hex value is brighter.
	// In light mode, a lower hex value is brighter (closer to black on white).
	// This creates the desired visual hierarchy in both themes.
	p.Grey.C300 = lipgloss.AdaptiveColor{Light: "#AFAFAF", Dark: "#626262"}
	p.Grey.C400 = lipgloss.AdaptiveColor{Light: "#7D7D7D", Dark: "#878787"}
	p.Grey.C500 = lipgloss.AdaptiveColor{Light: "#606060", Dark: "#B4B4B4"}
	p.Grey.C100 = lipgloss.AdaptiveColor{Dark: "#FFFFFF", Light: "#11181C"}

	p.Special.Success = lipgloss.Color("#40F284")
	p.Special.Error = lipgloss.Color("#FF4D4D")

	return p
}

// A struct to hold all our styles
type styles struct {
	App                 lipgloss.Style
	AppTitle            lipgloss.Style
	Title               lipgloss.Style
	Error               lipgloss.Style
	Success             lipgloss.Style
	Help                lipgloss.Style
	Spinner             lipgloss.Style
	LoadingText         lipgloss.Style
	TextInputPrompt     lipgloss.Style
	TextInputCursor     lipgloss.Style
	FinalMessage        lipgloss.Style
	FinalMessageHeader  lipgloss.Style
	FinalMessageContent lipgloss.Style
	ListItem            lipgloss.Style
	SelectedListItem    lipgloss.Style
	QuestionStyle       lipgloss.Style
	InfoTitleStyle      lipgloss.Style
	InfoValueStyle      lipgloss.Style
}

// DefaultStyles returns a struct with our default styles
func DefaultStyles() styles {
	colors := NewPalette()

	return styles{
		App: lipgloss.NewStyle().Margin(1, 2),
		AppTitle: lipgloss.NewStyle().
			Foreground(colors.Primary.Base).
			Bold(true).
			Underline(true),
		Title: lipgloss.NewStyle().
			Foreground(colors.Primary.Base).
			Bold(true).
			Underline(true).
			Padding(0, 1),
		Error: lipgloss.NewStyle().
			Foreground(colors.Special.Error).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(colors.Special.Success),
		Help: lipgloss.NewStyle().
			Foreground(colors.Grey.C300),
		Spinner: lipgloss.NewStyle().
			Foreground(colors.Primary.Base),
		LoadingText: lipgloss.NewStyle().
			Foreground(colors.Grey.C100).
			MarginLeft(1),
		TextInputPrompt: lipgloss.NewStyle().
			Foreground(colors.Primary.Base),
		TextInputCursor: lipgloss.NewStyle().
			Foreground(colors.Primary.Base),
		FinalMessage: lipgloss.NewStyle().
			Margin(1, 0).
			Padding(1, 2).
			Border(lipgloss.ASCIIBorder(), true).
			BorderForeground(colors.Primary.Base),
		FinalMessageHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Primary.Base).
			MarginBottom(1),
		FinalMessageContent: lipgloss.NewStyle().
			Foreground(colors.Grey.C100),
		ListItem: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(colors.Grey.C400), // Use mid-grey for list items
		SelectedListItem: lipgloss.NewStyle().
			Foreground(colors.Primary.Base).
			PaddingLeft(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(colors.Primary.Base),
		QuestionStyle:  lipgloss.NewStyle().Foreground(colors.Grey.C500),
		InfoTitleStyle: lipgloss.NewStyle().Foreground(colors.Grey.C300), // Use help color for titles
		InfoValueStyle: lipgloss.NewStyle().Foreground(colors.Primary.Light),
	}
}
