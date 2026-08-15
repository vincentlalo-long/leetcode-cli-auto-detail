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
			// Support `leet <cmd> --help` / `leet <cmd> -h`.
			if wantsHelp(args) && cmdName != "help" && cmdName != "man" && cmdName != "--help" && cmdName != "-h" {
				if err := commands.Registry["help"]([]string{cmdName}, cfg, commands.HeadlessUI{}); err != nil {
					os.Exit(1)
				}
				return
			}
			if err := handler(args, cfg, commands.HeadlessUI{}); err != nil {
				os.Exit(1)
			}
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

func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}
