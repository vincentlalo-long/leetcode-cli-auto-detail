package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"leetcli/internal/config"
	"leetcli/internal/template"
)

func OpenProblem(args []string, cfg *config.Config, ui UI) error {
	ui.WriteOutput(MsgPlain, "--- Open Problem ---\n")

	problemNum := ""
	if len(args) > 0 {
		problemNum = args[0]
	}
	if problemNum == "" {
		problemNum = ui.PromptText("Enter problem number to open")
	}
	if problemNum == "" {
		return fmt.Errorf("problem number is required")
	}

	baseDir := cfg.BaseDir
	if baseDir == "" {
		ui.WriteOutput(MsgError, "Invalid base directory in config")
		return fmt.Errorf("invalid base directory in config")
	}
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		ui.WriteOutput(MsgError, "Base directory does not exist: %s", baseDir)
		return fmt.Errorf("base directory does not exist: %s", baseDir)
	}

	ui.WriteOutput(MsgInfo, "Searching for problem...")
	languages := template.NormalizeLanguages(nil)
	for k, v := range cfg.Languages {
		languages[k] = template.LanguageInfo{Label: v.Label, Ext: v.Ext}
	}
	var exts []string
	for _, info := range languages {
		exts = append(exts, info.Ext)
	}

	allFiles := template.GetAllSolutionFiles(baseDir, exts)
	var matches []string
	for _, f := range allFiles {
		base := filepath.Base(f)
		if template.MatchesProblemNumber(base, problemNum) {
			matches = append(matches, f)
		}
	}

	if len(matches) == 0 {
		ui.WriteOutput(MsgError, "Could not find local file for problem %s.", problemNum)
		ui.WriteOutput(MsgInfo, "Try running 'add' or 'daily' first.")
		return fmt.Errorf("could not find local file for problem %s", problemNum)
	}

	targetFile := matches[0]
	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, f := range matches {
			names[i] = filepath.Base(f)
		}
		selected := ui.PromptSelect("Multiple matches. Select file to open", names)
		for i, n := range names {
			if n == selected {
				targetFile = matches[i]
				break
			}
		}
	}

	targetDir := filepath.Dir(targetFile)
	ui.WriteOutput(MsgSuccess, "Found problem: %s", filepath.Base(targetFile))

	readmePath := filepath.Join(targetDir, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		content, _ := os.ReadFile(readmePath)
		ui.WriteOutput(MsgPlain, "\n--- Problem Description ---\n")
		ui.WriteOutput(MsgPlain, "%s", string(content))
		ui.WriteOutput(MsgPlain, "---\n")
	} else {
		ui.WriteOutput(MsgInfo, "No README.md found for this problem.")
	}

	editor := cfg.GetEditor()
	ui.WriteOutput(MsgInfo, "Opening with %s...", editor)
	cmd := exec.Command(editor, targetFile)
	if err := cmd.Start(); err != nil {
		ui.WriteOutput(MsgError, "Failed to open editor '%s': %v", editor, err)
		ui.WriteOutput(MsgInfo, "You can manually open: %s", targetFile)
		return fmt.Errorf("failed to open editor '%s': %w", editor, err)
	}
	ui.WriteOutput(MsgSuccess, "Problem opened in editor!")
	return nil
}
