package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"leetcli/internal/config"
	"leetcli/internal/ui"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", err)
	}
	if cfg == nil {
		cfg = config.Default()
	}

	m := ui.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
