package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"leetcli/internal/config"
	"leetcli/internal/template"
)

func Sync(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Sync LeetCode Workspace ---\n")

	baseDir := cfg.BaseDir
	if baseDir == "" {
		ui.WriteOutput(MsgError, "Base directory not configured.")
		return
	}
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		ui.WriteOutput(MsgError, "Base directory '%s' does not exist.", baseDir)
		return
	}

	gitDir := filepath.Join(baseDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		ui.WriteOutput(MsgError, "'%s' is not a Git repository.", baseDir)
		ui.WriteOutput(MsgInfo, "Please initialize git and set up a remote repository first:")
		ui.WriteOutput(MsgPlain, "  cd %s", baseDir)
		ui.WriteOutput(MsgPlain, "  git init")
		ui.WriteOutput(MsgPlain, "  git remote add origin <your-repo-url>")
		return
	}

	ui.WriteOutput(MsgInfo, "Syncing workspace: %s", baseDir)

	ui.WriteOutput(MsgInfo, "Updating root README.md index...")
	if _, count, err := template.GenerateRootReadme(baseDir, cfg); err == nil {
		ui.WriteOutput(MsgSuccess, "Updated README.md index with %d problem(s).", count)
	} else {
		ui.WriteOutput(MsgError, "Warning: Failed to update root README.md: %v", err)
	}

	ui.WriteOutput(MsgInfo, "Staging changes...")
	if out, err := runGit(baseDir, "add", "."); err != nil {
		ui.WriteOutput(MsgError, "Failed to stage changes: %v", err)
		if out != "" {
			ui.WriteOutput(MsgPlain, "%s", out)
		}
		return
	}

	statusOut, _ := runGit(baseDir, "status", "--porcelain")
	if strings.TrimSpace(statusOut) == "" {
		ui.WriteOutput(MsgInfo, "No changes to sync. Workspace is up to date.")
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	commitMsg := fmt.Sprintf("LeetCode Sync: %s", timestamp)
	ui.WriteOutput(MsgInfo, "Committing changes...")
	if out, err := runGit(baseDir, "commit", "-m", commitMsg); err != nil {
		ui.WriteOutput(MsgError, "Failed to commit changes: %v", err)
		if out != "" {
			ui.WriteOutput(MsgPlain, "%s", out)
		}
		return
	}

	ui.WriteOutput(MsgInfo, "Pulling latest changes from remote (if any)...")
	runGit(baseDir, "pull", "--rebase")

	ui.WriteOutput(MsgInfo, "Pushing to remote...")
	if out, err := runGit(baseDir, "push"); err != nil {
		ui.WriteOutput(MsgError, "Failed to push changes.")
		ui.WriteOutput(MsgInfo, "You might need to set upstream branch manually:")
		ui.WriteOutput(MsgPlain, "  cd %s", baseDir)
		ui.WriteOutput(MsgPlain, "  git push -u origin main")
		if out != "" {
			ui.WriteOutput(MsgPlain, "%s", out)
		}
		return
	}

	ui.WriteOutput(MsgSuccess, "Successfully synced LeetCode workspace!")
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}


