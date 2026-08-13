package commands

import (
	"fmt"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

func Hint(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Get Problem Hints ---\n")

	problemNum := ""
	if len(args) > 0 {
		problemNum = args[0]
	}
	if problemNum == "" {
		problemNum = ui.PromptText("Enter problem number")
	}
	if problemNum == "" {
		return
	}

	ui.WriteOutput(MsgInfo, "Looking up problem...")
	problemData, err := api.GetProblemByID(problemNum)
	if err != nil {
		ui.WriteOutput(MsgError, "Could not find problem with ID %s", problemNum)
		return
	}

	ui.WriteOutput(MsgSuccess, "Found: %s", problemData.Title)
	ui.WriteOutput(MsgInfo, "Fetching hints...")

	details, err := api.GetProblemDetails(problemData.Slug)
	if err != nil {
		ui.WriteOutput(MsgError, "Failed to fetch problem details from LeetCode.")
		return
	}

	if len(details.Hints) == 0 {
		ui.WriteOutput(MsgInfo, "No official hints available for this problem.")
		return
	}

	ui.WriteOutput(MsgPlain, "\nHints for %s (%d total):", problemData.Title, len(details.Hints))

	for i, hint := range details.Hints {
		hintText := template.DecodeHTMLEntities(template.StripHTMLTags(hint))
		ui.WriteOutput(MsgPlain, "  Hint %d: %s", i+1, hintText)

		if i < len(details.Hints)-1 {
			cont := ui.PromptConfirm("Next hint?")
			if !cont {
				break
			}
		}
	}

	ui.WriteOutput(MsgSuccess, "Finished showing hints.")
}

var _ = fmt.Sprintf
