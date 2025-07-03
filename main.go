package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// --- TEA MESSAGES ---

type itemsLoadedMsg struct{ items []OpListItem }
type itemSelectedMsg struct{ envContent string }
type fileWrittenMsg struct{ path string }
type errorMsg struct{ err error }

// --- TEA MODEL ---

type programState int

const (
	vaultInputState programState = iota
	loadingState
	itemListState
	fileInputState
	finishedState
)

type model struct {
	state         programState
	styles        styles
	vaultInput    textinput.Model
	fileInput     textinput.Model
	spinner       spinner.Model
	items         []OpListItem
	cursor        int
	selectedVault string
	selectedItem  OpListItem
	envContent    string
	finalPath     string
	err           error
	width         int
	height        int
}

func initialModel() model {
	styles := DefaultStyles()

	// Vault Name Input
	ti := textinput.New()
	ti.Placeholder = "Personal"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30
	ti.Prompt = "❯ "
	ti.PromptStyle = styles.TextInputPrompt
	ti.Cursor.Style = styles.TextInputCursor

	// File Name Input
	fi := textinput.New()
	fi.Placeholder = ".env"
	fi.CharLimit = 156
	fi.Width = 30
	fi.Prompt = "❯ "
	fi.PromptStyle = styles.TextInputPrompt
	fi.Cursor.Style = styles.TextInputCursor

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Spinner

	return model{
		state:      vaultInputState,
		styles:     styles,
		vaultInput: ti,
		fileInput:  fi,
		spinner:    s,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	// --- Custom Messages ---
	case itemsLoadedMsg:
		m.state = itemListState
		m.items = msg.items
		return m, nil

	case itemSelectedMsg:
		m.state = fileInputState
		m.envContent = msg.envContent
		m.fileInput.Focus()
		return m, textinput.Blink

	case fileWrittenMsg:
		m.state = finishedState
		m.finalPath = msg.path
		return m, tea.Quit

	case errorMsg:
		m.err = msg.err
		return m, tea.Quit

	case spinner.TickMsg:
		if m.state == loadingState {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	// --- Update Components ---
	switch m.state {
	case vaultInputState:
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			vaultName := m.vaultInput.Value()
			if vaultName != "" {
				m.selectedVault = vaultName
				m.state = loadingState
				cmds = append(cmds, m.spinner.Tick, fetchOpItems(vaultName))
			}
		} else {
			m.vaultInput, cmd = m.vaultInput.Update(msg)
			cmds = append(cmds, cmd)
		}

	case itemListState:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.items)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.items) > 0 {
					m.selectedItem = m.items[m.cursor]
					m.state = loadingState
					cmds = append(cmds, m.spinner.Tick, fetchItemDetails(m.items[m.cursor].ID))
				}
			}
		}

	case fileInputState:
		if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEnter {
			fileName := m.fileInput.Value()
			if fileName == "" {
				fileName = ".env" // Default to .env
			}
			// Create .env file in current directory
			currentDir, _ := os.Getwd()
			filePath := currentDir + "/" + fileName
			cmds = append(cmds, writeFile(filePath, m.envContent))
		} else {
			m.fileInput, cmd = m.fileInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return "" // Error is handled in main
	}
	if m.state == finishedState {
		return "" // Final message is handled in main
	}

	var s string
	if m.selectedVault != "" {
		s += m.styles.InfoTitleStyle.Render("Vault: ") + m.styles.InfoValueStyle.Render(m.selectedVault) + "\n"
	}
	if m.selectedItem.ID != "" {
		s += m.styles.InfoTitleStyle.Render("Item: ") + m.styles.InfoValueStyle.Render(m.selectedItem.ItemTitle) + "\n"
	}
	if s != "" {
		s += "\n"
	}

	switch m.state {
	case vaultInputState:
		s += m.styles.QuestionStyle.Render("Which 1Password vault should I search?") + "\n\n" +
			m.vaultInput.View() +
			"\n\n" + m.styles.Help.Render("enter: submit, ctrl+c: quit")
	case loadingState:
		s += fmt.Sprintf("%s%s", m.spinner.View(), m.styles.LoadingText.Render("Fetching from 1Password..."))
	case itemListState:
		list := ""
		for i, item := range m.items {
			row := fmt.Sprintf("%d) %s", i+1, item.ItemTitle)
			if m.cursor == i {
				list += m.styles.SelectedListItem.Render(row)
			} else {
				list += m.styles.ListItem.Render(row)
			}
			list += "\n"
		}
		s += m.styles.QuestionStyle.Render("Select an item from the vault:") + "\n\n" + list + "\n" + m.styles.Help.Render("↑/↓: navigate, enter: select, ctrl+c: quit")

	case fileInputState:
		s += m.styles.QuestionStyle.Render("Enter a filename for the .env file (default: .env):") + "\n\n" +
			m.fileInput.View() +
			"\n\n" + m.styles.Help.Render("enter: submit, ctrl+c: quit")
	}
	finalView := m.styles.AppTitle.Render("Envelope") + "\n\n" + s
	return m.styles.App.Render(finalView)
}

// Command to write the file
func writeFile(path, content string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return errorMsg{err}
		}
		return fileWrittenMsg{path: path}
	}
}

// --- MAIN FUNCTION ---

func main() {
	// Check for `op` CLI
	if _, err := exec.LookPath("op"); err != nil {
		fmt.Println("\n" + DefaultStyles().Error.Render("Error: 1Password CLI ('op') not found in your PATH."))
		fmt.Println(DefaultStyles().Help.Render("Please install it to continue: https://1password.com/downloads/cli/"))
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())

	m, err := p.Run()
	if err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}

	finalModel, _ := m.(model)

	if finalModel.err != nil {
		header := finalModel.styles.FinalMessageHeader.Render("✗ Oh no! An error occurred.")
		content := finalModel.styles.Error.Render(finalModel.err.Error())
		fmt.Println(finalModel.styles.FinalMessage.Render(header + "\n" + content))
		os.Exit(1)
	}

	if finalModel.state == finishedState {
		header := finalModel.styles.FinalMessageHeader.Render("✓ Success!")
		details := finalModel.styles.InfoTitleStyle.Render("Vault: ") + finalModel.styles.InfoValueStyle.Render(finalModel.selectedVault) + "\n" +
			finalModel.styles.InfoTitleStyle.Render("Item: ") + finalModel.styles.InfoValueStyle.Render(finalModel.selectedItem.ItemTitle)
		content := finalModel.styles.FinalMessageContent.Render(details) +
			"\n\n" +
			finalModel.styles.FinalMessageContent.Render("Your .env file was created at:") +
			"\n" + finalModel.styles.Success.Render(finalModel.finalPath)
		fmt.Println(finalModel.styles.FinalMessage.Render(header + "\n" + content))
	} else {
		fmt.Println("\nOperation cancelled.")
	}
}
