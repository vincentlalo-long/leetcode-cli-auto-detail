package commands

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

func Random(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Random LeetCode Problem ---\n")

	difficultyChoices := []string{"Any", "Easy", "Medium", "Hard"}
	selectedDifficulty := ui.PromptSelect("Select difficulty level", difficultyChoices)

	ui.WriteOutput(MsgInfo, "Fetching problems...")
	data, err := api.GetAllProblems()
	if err != nil {
		ui.WriteOutput(MsgError, "Could not fetch problems. Please check your connection.")
		return
	}

	type candidate struct {
		FrontendID int
		Title      string
		Slug       string
		Level      int
		PaidOnly   bool
	}
	var candidates []candidate

	diffLevel := 0
	if selectedDifficulty != "Any" {
		diffLevel = template.DiffToLevel[selectedDifficulty]
	}

	for _, p := range data.StatStatusPairs {
		if p.PaidOnly {
			continue
		}
		if diffLevel > 0 && p.Difficulty.Level != diffLevel {
			continue
		}
		candidates = append(candidates, candidate{
			FrontendID: p.Stat.FrontendQuestionID,
			Title:      p.Stat.QuestionTitle,
			Slug:       p.Stat.QuestionTitleSlug,
			Level:      p.Difficulty.Level,
		})
	}

	if len(candidates) == 0 {
		ui.WriteOutput(MsgError, "No free %s problems found.", selectedDifficulty)
		return
	}

	prob := candidates[rand.Intn(len(candidates))]
	problemNum := fmt.Sprintf("%d", prob.FrontendID)
	problemName := prob.Title
	slug := prob.Slug
	difficulty := template.LevelToName[prob.Level]
	link := fmt.Sprintf("https://leetcode.com/problems/%s/", slug)

	ui.WriteOutput(MsgPlain, "\nYour Random Problem:")
	ui.WriteOutput(MsgPlain, "  %s. %s", problemNum, problemName)
	ui.WriteOutput(MsgInfo, "Difficulty: %s", difficulty)
	ui.WriteOutput(MsgInfo, "Link: %s", link)

	addIt := ui.PromptConfirm("Add this problem to your workspace?")
	if !addIt {
		return
	}

	dataStructures := cfg.GetDataStructures()
	if len(dataStructures) == 0 {
		ui.WriteOutput(MsgError, "No data structures found. Add one first!")
		return
	}

	dsChoices := []string{"[Uncategorized]"}
	for k := range dataStructures {
		dsChoices = append(dsChoices, k)
	}
	dsChoices = append(dsChoices, "Add new data structure")

	selected := ui.PromptSelect("Select data structure", dsChoices)
	if selected == "Add new data structure" {
		name := ui.PromptText("Data structure name (e.g., tree)")
		folder := name
		if folder == "" {
			return
		}
		folderInput := ui.PromptText(fmt.Sprintf("Folder name (press Enter for '%s')", name))
		if folderInput != "" {
			folder = folderInput
		}
		if cfg.AddDataStructure(name, folder) {
			cfg.Save()
			ui.WriteOutput(MsgSuccess, "Added data structure: %s -> %s", name, folder)
		}
		dataStructures = cfg.GetDataStructures()
		dsChoices = []string{"[Uncategorized]"}
		for k := range dataStructures {
			dsChoices = append(dsChoices, k)
		}
		selected = ui.PromptSelect("Select data structure", dsChoices)
	}

	dsFolder := "uncategorized"
	if selected != "[Uncategorized]" {
		if f, ok := dataStructures[selected]; ok {
			dsFolder = f
		}
	}

	languages := template.NormalizeLanguages(nil)
	for k, v := range cfg.Languages {
		languages[k] = template.LanguageInfo{Label: v.Label, Ext: v.Ext}
	}
	langChoices, langMapping := template.GetLanguageChoices(languages, cfg.DefaultLanguage)
	langChoice := ui.PromptSelect("Select language", langChoices)
	langKey := langMapping[langChoice]
	langExt := languages[langKey].Ext

	folderName := fmt.Sprintf("%s-%s", problemNum, slug)
	problemDir, err := template.CreateProblemDirectory(cfg.BaseDir, dsFolder, folderName)
	if err != nil {
		ui.WriteOutput(MsgError, "Failed to create directory: %v", err)
		return
	}

	problemFile := filepath.Join(problemDir, fmt.Sprintf("%s_%s.%s", problemNum, problemName, langExt))
	if _, err := os.Stat(problemFile); err == nil {
		ui.WriteOutput(MsgError, "Problem file already exists!")
		return
	}

	ui.WriteOutput(MsgInfo, "Fetching full problem details...")
	details, _ := api.GetProblemDetails(slug)

	tagsStr := "None"
	if details != nil {
		tags := make([]string, len(details.TopicTags))
		for i, t := range details.TopicTags {
			tags[i] = t.Name
		}
		if len(tags) > 0 {
			tagsStr = strings.Join(tags, ", ")
		}
	}

	content := template.BuildProblemTemplate(langKey, problemNum, problemName, link,
		difficulty, tagsStr, selected)
	os.WriteFile(problemFile, []byte(content), 0644)

	if details != nil && details.Content != "" {
		readmePath := filepath.Join(problemDir, "README.md")
		mdContent := template.StripHTMLTags(details.Content)
		cleanNum := strings.TrimLeft(problemNum, "0")
		if cleanNum == "" {
			cleanNum = problemNum
		}
		readmeContent := template.MakeReadmeContent(cleanNum, problemName, link,
			difficulty, tagsStr, mdContent)
		os.WriteFile(readmePath, []byte(readmeContent), 0644)
		ui.WriteOutput(MsgSuccess, "Saved problem description to README.md")
	}

	ui.WriteOutput(MsgSuccess, "Created problem directory and files")
	ui.WriteOutput(MsgInfo, "Path: %s", problemFile)
}
