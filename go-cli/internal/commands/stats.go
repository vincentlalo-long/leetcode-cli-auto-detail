package commands

import (
	"fmt"
	"os"
	"strings"

	"leetcli/internal/config"
	"leetcli/internal/template"
)

func Stats(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Problem Statistics ---\n")

	baseDir := cfg.BaseDir
	dataStructures := cfg.GetDataStructures()

	if baseDir == "" {
		ui.WriteOutput(MsgError, "Missing 'base_dir' in config")
		return
	}
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		ui.WriteOutput(MsgError, "Base directory does not exist: %s", baseDir)
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

	files := template.GetAllSolutionFiles(baseDir, exts)
	totalProblems := len(files)
	totalSolutions := 0

	byStructure := make(map[string]int)
	for name := range dataStructures {
		byStructure[name] = 0
	}
	unmatched := 0

	for _, f := range files {
		matched := false
		for name, folder := range dataStructures {
			token := strings.ToLower("/" + folder + "/")
			if strings.Contains(strings.ToLower(f), token) {
				byStructure[name]++
				matched = true
				break
			}
		}
		if !matched {
			unmatched++
		}

		if content, err := os.ReadFile(f); err == nil {
			totalSolutions += template.CountSolutions(string(content))
		}
	}

	ui.WriteOutput(MsgPlain, "\nOverview:")
	ui.WriteOutput(MsgInfo, "Base directory: %s", baseDir)
	ui.WriteOutput(MsgInfo, "Total problems: %d", totalProblems)
	ui.WriteOutput(MsgInfo, "Total solutions: %d", totalSolutions)

	ui.WriteOutput(MsgPlain, "\nBy Data Structure:")
	if len(dataStructures) == 0 {
		ui.WriteOutput(MsgInfo, "No data structures configured")
	} else {
		for _, name := range sortedKeys(byStructure) {
			ui.WriteOutput(MsgPlain, "  %-20s %d", name, byStructure[name])
		}
	}

	if unmatched > 0 {
		ui.WriteOutput(MsgInfo, "Unmatched problems: %d", unmatched)
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

var _ = fmt.Sprintf
