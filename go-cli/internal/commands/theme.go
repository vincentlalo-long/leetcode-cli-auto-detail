package commands

import (
	"leetcli/internal/config"
)

var AvailableThemes = []string{"default", "dracula", "hacker", "sunset"}

func Theme(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Select UI Theme ---\n")

	currentTheme := cfg.GetTheme()
	choices := make([]string, len(AvailableThemes))
	for i, t := range AvailableThemes {
		label := t
		if t == currentTheme {
			label += " (Current)"
		}
		choices[i] = label
	}

	selected := ui.PromptSelect("Choose a theme", choices)
	if selected == "" {
		return
	}

	themeName := selected
	for _, t := range AvailableThemes {
		if selected == t+" (Current)" {
			themeName = t
			break
		}
	}

	if themeName == currentTheme {
		ui.WriteOutput(MsgSuccess, "Theme is already set to %s", themeName)
		return
	}

	cfg.SetTheme(themeName)
	cfg.Save()
	ui.WriteOutput(MsgSuccess, "Theme successfully changed to %s!", themeName)
	ui.WriteOutput(MsgInfo, "Restart the CLI to apply the new theme.")
}
