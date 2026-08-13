package commands

import (
	"fmt"

	"leetcli/internal/config"
)

// Version is the semantic version of the Leet CLI.
const Version = "0.1.0"

// VersionCommand prints the CLI version and build info.
func VersionCommand(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "leet CLI v%s", Version)
	ui.WriteOutput(MsgInfo, "A terminal workspace manager for LeetCode problems.")
	ui.WriteOutput(MsgPlain, "Run 'leet' for the interactive UI, 'leet help' for commands.")
	_ = fmt.Sprintf
}