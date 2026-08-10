package commands

import (
	"os"
	"path/filepath"
	"strings"

	"leetcli/internal/config"
)

func CleanWorkspace(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- Workspace Build Cleanup ---\n")

	baseDir := cfg.BaseDir
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		ui.WriteOutput(MsgError, "Base directory '%s' does not exist.", baseDir)
		return
	}

	ui.WriteOutput(MsgInfo, "Scanning '%s' for temporary build artifacts...", baseDir)

	cleanedCount := 0
	var totalBytes int64 = 0

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		ext := strings.ToLower(filepath.Ext(path))

		isTemp := false
		if strings.HasPrefix(base, "leet_run_") ||
			ext == ".class" ||
			ext == ".tmp" ||
			base == "a.out" ||
			strings.HasSuffix(base, ".exe~") {
			isTemp = true
		}

		if isTemp {
			size := info.Size()
			if removeErr := os.Remove(path); removeErr == nil {
				cleanedCount++
				totalBytes += size
				ui.WriteOutput(MsgPlain, "  ✓ Removed: %s (%d bytes)", base, size)
			}
		}
		return nil
	})

	if err != nil {
		ui.WriteOutput(MsgError, "Error scanning workspace: %v", err)
		return
	}

	if cleanedCount == 0 {
		ui.WriteOutput(MsgSuccess, "Workspace is clean! No build artifacts found.")
	} else {
		mb := float64(totalBytes) / (1024 * 1024)
		ui.WriteOutput(MsgSuccess, "Cleaned %d temporary file(s) (%.2f MB reclaimed).", cleanedCount, mb)
	}
}
