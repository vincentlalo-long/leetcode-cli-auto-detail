package commands

import (
	"encoding/json"
	"strings"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

type SimilarProblem struct {
	Title       string `json:"title"`
	TitleSlug   string `json:"titleSlug"`
	Difficulty  string `json:"difficulty"`
}

func Similar(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Find Similar Problems ---\n")

	problemNum := ui.PromptText("Enter problem number")
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
	ui.WriteOutput(MsgInfo, "Fetching similar problems...")

	details, err := api.GetProblemDetails(problemData.Slug)
	if err != nil {
		ui.WriteOutput(MsgError, "Failed to fetch problem details from LeetCode.")
		return
	}

	if details.SimilarQuestions == "" {
		ui.WriteOutput(MsgInfo, "No similar problems available for this problem.")
		return
	}

	var similar []SimilarProblem
	if err := json.Unmarshal([]byte(details.SimilarQuestions), &similar); err != nil {
		ui.WriteOutput(MsgError, "Failed to parse similar problems data.")
		return
	}

	if len(similar) == 0 {
		ui.WriteOutput(MsgInfo, "No similar problems found for this problem.")
		return
	}

	languages := template.NormalizeLanguages(nil)
	for k, v := range cfg.Languages {
		languages[k] = template.LanguageInfo{Label: v.Label, Ext: v.Ext}
	}
	var exts []string
	for _, info := range languages {
		exts = append(exts, info.Ext)
	}
	solvedMap := template.GetSolvedMap(cfg.BaseDir, exts)

	ui.WriteOutput(MsgPlain, "\nSimilar Problems to %s (%d total):", problemData.Title, len(similar))

	for i, sp := range similar {
		difficulty := sp.Difficulty
		badge := "[○ Unsolved]"
		if solvedMap[sp.TitleSlug] || solvedMap[api.Slugify(sp.Title)] {
			badge = "[✔ Solved]"
		}
		ui.WriteOutput(MsgPlain, "  %d. %-35s (%s) %s", i+1, sp.Title, difficulty, badge)
	}

	ui.WriteOutput(MsgSuccess, "Finished showing similar problems.")
}

var _ = strings.Join
var _ = template.CountSolutions
