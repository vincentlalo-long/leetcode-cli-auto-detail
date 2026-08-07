package commands

import (
	"leetcli/internal/config"
)

func ManageStructures(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Manage Data Structures ---\n")

	action := ui.PromptSelect("What would you like to do?", []string{
		"List all data structures",
		"Add new data structure",
		"Remove data structure",
	})

	switch action {
	case "List all data structures":
		listStructures(cfg, ui)
	case "Add new data structure":
		addNewStructure(cfg, ui)
	case "Remove data structure":
		removeStructure(cfg, ui)
	}
}

func listStructures(cfg *config.Config, ui UI) {
	structures := cfg.GetDataStructures()
	if len(structures) == 0 {
		ui.WriteOutput(MsgInfo, "No data structures found!")
		return
	}
	ui.WriteOutput(MsgPlain, "Available Data Structures:")
	for name, folder := range structures {
		ui.WriteOutput(MsgPlain, "  %-20s %s", name, folder)
	}
}

func addNewStructure(cfg *config.Config, ui UI) {
	name := ui.PromptText("Data structure name (e.g., tree)")
	if name == "" {
		ui.WriteOutput(MsgError, "Data structure name cannot be empty!")
		return
	}

	folder := ui.PromptText("Folder name (press Enter for same)")
	if folder == "" {
		folder = name
	}

	if cfg.AddDataStructure(name, folder) {
		cfg.Save()
		ui.WriteOutput(MsgSuccess, "Added data structure: %s -> %s", name, folder)
	} else {
		ui.WriteOutput(MsgError, "Data structure '%s' already exists!", name)
	}
}

func removeStructure(cfg *config.Config, ui UI) {
	structures := cfg.GetDataStructures()
	if len(structures) == 0 {
		ui.WriteOutput(MsgInfo, "No data structures to remove!")
		return
	}

	names := make([]string, 0, len(structures))
	for n := range structures {
		names = append(names, n)
	}

	name := ui.PromptSelect("Select data structure to remove", names)
	if name == "" {
		return
	}

	confirm := ui.PromptConfirm("Remove '" + name + "'?")
	if !confirm {
		ui.WriteOutput(MsgInfo, "Operation cancelled")
		return
	}

	if cfg.RemoveDataStructure(name) {
		cfg.Save()
		ui.WriteOutput(MsgSuccess, "Removed data structure: %s", name)
	} else {
		ui.WriteOutput(MsgError, "Data structure '%s' not found!", name)
	}
}
