package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

func AddProblem(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Add New Problem ---\n")

	problemNum := ui.PromptText("Problem number")
	if problemNum == "" {
		return
	}

	suggestedName := ""
	ui.WriteOutput(MsgInfo, "Looking up problem details...")
	problemData, err := api.GetProblemByID(problemNum)
	if err == nil && problemData != nil {
		suggestedName = problemData.Title
		ui.WriteOutput(MsgSuccess, "Found: %s", suggestedName)
	} else {
		ui.WriteOutput(MsgError, "Could not find problem with ID %s", problemNum)
	}

	problemName := ui.PromptText(fmt.Sprintf("Problem name (suggested: %s)", suggestedName))
	if problemName == "" && suggestedName != "" {
		problemName = suggestedName
	}
	if problemName == "" {
		ui.WriteOutput(MsgError, "Problem name cannot be empty")
		return
	}

	dataStructures := cfg.GetDataStructures()
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

	slug := api.Slugify(problemName)
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

	ui.WriteOutput(MsgInfo, "Fetching problem details from LeetCode...")
	details, _ := api.GetProblemDetails(slug)

	content := ""
	if details != nil {
		tags := make([]string, len(details.TopicTags))
		for i, t := range details.TopicTags {
			tags[i] = t.Name
		}
		tagsStr := "None"
		if len(tags) > 0 {
			tagsStr = strings.Join(tags, ", ")
		}
		link := fmt.Sprintf("https://leetcode.com/problems/%s/", slug)
		snippet := details.GetCodeSnippet(langKey)
		content = template.BuildProblemTemplateWithSnippet(langKey, problemNum, details.Title, link,
			details.Difficulty, tagsStr, selected, snippet)
		ui.WriteOutput(MsgSuccess, "Found: %s (%s)", details.Title, details.Difficulty)
	} else {
		content = template.BuildProblemTemplate(langKey, problemNum, problemName, "",
			"Unknown", "None", selected)
		ui.WriteOutput(MsgError, "Could not fetch problem details. Using basic template.")
	}

	err = os.WriteFile(problemFile, []byte(content), 0644)
	if err != nil {
		ui.WriteOutput(MsgError, "Failed to write file: %v", err)
		return
	}

	if details != nil && details.Content != "" {
		readmePath := filepath.Join(problemDir, "README.md")
		mdContent := template.StripHTMLTags(details.Content)
		cleanNum := strings.TrimLeft(problemNum, "0")
		if cleanNum == "" {
			cleanNum = problemNum
		}
		link := fmt.Sprintf("https://leetcode.com/problems/%s/", slug)
		readmeContent := template.MakeReadmeContent(cleanNum, details.Title, link,
			details.Difficulty, strings.Join(func() []string {
				tags := make([]string, len(details.TopicTags))
				for i, t := range details.TopicTags {
					tags[i] = t.Name
				}
				return tags
			}(), ", "), mdContent)
		os.WriteFile(readmePath, []byte(readmeContent), 0644)
		ui.WriteOutput(MsgSuccess, "Saved problem description to README.md")
	}

	ui.WriteOutput(MsgSuccess, "Created problem directory and files")
	ui.WriteOutput(MsgInfo, "Path: %s", problemFile)
}
