package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"leetcli/internal/commands"
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
		cfg.ResolveBaseDir()
	}

	// Non-interactive mode: leet <command> [args...]
	if len(os.Args) > 1 {
		cmdName := os.Args[1]
		args := os.Args[2:]

		if handler, ok := commands.Registry[cmdName]; ok {
			handler(args, cfg, commands.HeadlessUI{})
			return
		}
		fmt.Fprintf(os.Stderr, "Unknown command: '%s'. Available commands:\n", cmdName)
		for name := range commands.Registry {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		os.Exit(1)
	}

	m := ui.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
