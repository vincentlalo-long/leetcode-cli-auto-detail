package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

func OpenWebProblem(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Open LeetCode in Browser ---\n")

	problemNum := ""
	if len(args) > 0 {
		problemNum = args[0]
	} else {
		problemNum = ui.PromptText("Enter problem number (or title slug)")
	}
	if problemNum == "" {
		return
	}

	slug := api.Slugify(problemNum)

	languages := template.NormalizeLanguages(nil)
	for k, v := range cfg.Languages {
		languages[k] = template.LanguageInfo{Label: v.Label, Ext: v.Ext}
	}
	var exts []string
	for _, info := range languages {
		exts = append(exts, info.Ext)
	}

	allFiles := template.GetAllSolutionFiles(cfg.BaseDir, exts)
	for _, f := range allFiles {
		base := filepath.Base(f)
		if template.MatchesProblemNumber(base, problemNum) {
			content, err := os.ReadFile(f)
			if err == nil {
				inferred := template.InferSlugFromPath(f, string(content))
				if inferred != "" {
					slug = inferred
					break
				}
			}
		}
	}

	if problemData, err := api.GetProblemByID(problemNum); err == nil && problemData != nil {
		slug = problemData.Slug
	}

	url := fmt.Sprintf("https://leetcode.com/problems/%s/", slug)
	ui.WriteOutput(MsgInfo, "Opening in default browser: %s", url)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		ui.WriteOutput(MsgError, "Failed to open browser: %v", err)
		ui.WriteOutput(MsgInfo, "URL: %s", url)
		return
	}

	ui.WriteOutput(MsgSuccess, "Opened problem page in browser!")
}
