package commands

import (
	"strconv"
	"time"

	"leetcli/internal/config"
)

func Timer(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Practice Stopwatch / Timer ---\n")

	minutes := 0
	if len(args) > 0 {
		m, err := strconv.Atoi(args[0])
		if err == nil && m > 0 {
			minutes = m
		}
	}

	if minutes == 0 {
		choices := []string{
			"1. Easy Problem (20 minutes)",
			"2. Medium Problem (30 minutes)",
			"3. Hard Problem (45 minutes)",
			"4. Custom Duration",
		}
		selected := ui.PromptSelect("Select practice target duration", choices)
		switch {
		case stringsHasPrefix(selected, "1."):
			minutes = 20
		case stringsHasPrefix(selected, "2."):
			minutes = 30
		case stringsHasPrefix(selected, "3."):
			minutes = 45
		case stringsHasPrefix(selected, "4."):
			input := ui.PromptText("Enter duration in minutes (e.g. 15)")
			m, err := strconv.Atoi(input)
			if err == nil && m > 0 {
				minutes = m
			}
		}
	}

	if minutes <= 0 {
		ui.WriteOutput(MsgError, "Invalid duration specified.")
		return
	}

	ui.WriteOutput(MsgSuccess, "⏱️ Practice timer started for %d minutes!", minutes)
	ui.WriteOutput(MsgInfo, "Target end time: %s", time.Now().Add(time.Duration(minutes)*time.Minute).Format("15:04:05"))

	go func(mins int) {
		time.Sleep(time.Duration(mins) * time.Minute)
		ui.WriteOutput(MsgError, "\n⏰ ========================================")
		ui.WriteOutput(MsgError, "⏰ TIME IS UP! (%d minutes completed!)", mins)
		ui.WriteOutput(MsgError, "⏰ ========================================\n")
	}(minutes)
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
