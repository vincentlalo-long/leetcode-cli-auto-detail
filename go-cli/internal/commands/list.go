package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"leetcli/internal/config"
	"leetcli/internal/template"
)

type ProblemRecord struct {
	FilePath   string
	FileName   string
	Structure  string
	Solutions  int
}

func ListProblems(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- List Problems ---\n")

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

	records := collectProblems(baseDir, dataStructures, exts)
	if len(records) == 0 {
		ui.WriteOutput(MsgInfo, "No problem files found")
		return
	}

	structSet := make(map[string]bool)
	for _, r := range records {
		structSet[r.Structure] = true
	}
	availableStructs := make([]string, 0, len(structSet))
	for s := range structSet {
		availableStructs = append(availableStructs, s)
	}
	sort.Strings(availableStructs)

	filterChoices := append([]string{"All"}, availableStructs...)
	selectedStructure := ui.PromptSelect("Filter by data structure", filterChoices)
	unsolvedOnly := ui.PromptConfirm("Show only unsolved problems?")

	filtered := filterProblems(records, selectedStructure, unsolvedOnly)

	ui.WriteOutput(MsgPlain, "Overview")
	ui.WriteOutput(MsgInfo, "Total problems: %d", len(records))
	ui.WriteOutput(MsgInfo, "After filters: %d", len(filtered))
	if selectedStructure != "All" {
		ui.WriteOutput(MsgInfo, "Structure: %s", selectedStructure)
	}
	if unsolvedOnly {
		ui.WriteOutput(MsgInfo, "Mode: Unsolved only")
	} else {
		ui.WriteOutput(MsgInfo, "Mode: All")
	}

	if len(filtered) == 0 {
		ui.WriteOutput(MsgInfo, "No problems match selected filters")
		return
	}

	ui.WriteOutput(MsgPlain, "\nProblems:")
	for i, r := range filtered {
		status := "unsolved"
		if r.Solutions > 0 {
			status = fmt.Sprintf("%d solution(s)", r.Solutions)
		}
		ui.WriteOutput(MsgPlain, "  %d. %s (%s, %s)", i+1, r.FileName, r.Structure, status)
		ui.WriteOutput(MsgPlain, "     %s", r.FilePath)
	}
}

func collectProblems(baseDir string, dataStructures map[string]string, exts []string) []ProblemRecord {
	files := template.GetAllSolutionFiles(baseDir, exts)
	records := make([]ProblemRecord, 0, len(files))

	for _, f := range files {
		solutions := 0
		if content, err := os.ReadFile(f); err == nil {
			solutions = template.CountSolutions(string(content))
		}

		records = append(records, ProblemRecord{
			FilePath:  f,
			FileName:  filepath.Base(f),
			Structure: template.DetectStructure(f, dataStructures),
			Solutions: solutions,
		})
	}

	sort.Slice(records, func(i, j int) bool {
		re := regexp.MustCompile(`(\d+)`)
		mi := re.FindStringSubmatch(records[i].FileName)
		mj := re.FindStringSubmatch(records[j].FileName)
		ni := 999999
		nj := 999999
		if len(mi) > 1 {
			fmt.Sscanf(mi[1], "%d", &ni)
		}
		if len(mj) > 1 {
			fmt.Sscanf(mj[1], "%d", &nj)
		}
		if ni != nj {
			return ni < nj
		}
		return strings.ToLower(records[i].FileName) < strings.ToLower(records[j].FileName)
	})

	return records
}

func filterProblems(records []ProblemRecord, selectedStructure string, unsolvedOnly bool) []ProblemRecord {
	var filtered []ProblemRecord
	for _, r := range records {
		if selectedStructure != "All" && r.Structure != selectedStructure {
			continue
		}
		if unsolvedOnly && r.Solutions > 0 {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}
