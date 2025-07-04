package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// --- 1PASSWORD JSON STRUCTS ---

type OpListItem struct {
	ID        string   `json:"id"`
	ItemTitle string   `json:"title"`
	Tags      []string `json:"tags"`
}

func (i OpListItem) FilterValue() string { return i.ItemTitle }
func (i OpListItem) Title() string       { return i.ItemTitle }
func (i OpListItem) Description() string { return "" }

type OpItemDetail struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
	Fields   []Field   `json:"fields"`
}

type Section struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Field struct {
	ID      string  `json:"id"`
	Label   string  `json:"label"`
	Value   string  `json:"value"`
	Section Section `json:"section"`
}

// --- COMMANDS (I/O) ---

// createOpError cleans up error messages from the 'op' CLI.
func createOpError(output []byte, originalErr error) error {
	errMsg := strings.TrimSpace(string(output))
	if parts := strings.SplitN(errMsg, " ", 3); len(parts) == 3 && strings.HasPrefix(parts[0], "[ERROR]") {
		return fmt.Errorf("%s", parts[2])
	}
	return fmt.Errorf("op command failed: %w\nOutput: %s", originalErr, errMsg)
}

func fetchOpItems(vaultName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("op", "item", "list", "--vault", vaultName, "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errorMsg{createOpError(output, err)}
		}

		var opItems []OpListItem
		if err := json.Unmarshal(output, &opItems); err != nil {
			return errorMsg{fmt.Errorf("failed to parse JSON from 'op': %v", err)}
		}

		if len(opItems) == 0 {
			return errorMsg{fmt.Errorf("no items found in vault '%s'", vaultName)}
		}

		// Sort items by title
		sort.Slice(opItems, func(i, j int) bool {
			return opItems[i].ItemTitle < opItems[j].ItemTitle
		})

		return itemsLoadedMsg{items: opItems}
	}
}

func fetchItemDetails(itemID string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("op", "item", "get", itemID, "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errorMsg{createOpError(output, err)}
		}

		var detail OpItemDetail
		if err := json.Unmarshal(output, &detail); err != nil {
			return errorMsg{fmt.Errorf("failed to parse item detail JSON: %v", err)}
		}

		envContent := generateEnvContent(detail)
		return itemSelectedMsg{envContent: envContent}
	}
}

func generateEnvContent(detail OpItemDetail) string {
	var sb strings.Builder
	sectionLabels := make(map[string]string)
	for _, s := range detail.Sections {
		sectionLabels[s.ID] = s.Label
	}

	fieldsBySection := make(map[string][]Field)
	for _, f := range detail.Fields {
		if f.ID == "notesPlain" || f.Value == "" {
			continue
		}
		fieldsBySection[f.Section.ID] = append(fieldsBySection[f.Section.ID], f)
	}

	for sectionID, fields := range fieldsBySection {
		label, ok := sectionLabels[sectionID]
		if ok && label != "" {
			sb.WriteString(fmt.Sprintf("# %s\n", label))
		} else {
			sb.WriteString("# Other\n")
		}

		for _, field := range fields {
			sb.WriteString(fmt.Sprintf(`%s="%s"`+"\n", field.Label, field.Value))
		}
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String()) + "\n"
}
