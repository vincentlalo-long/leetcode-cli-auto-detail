package commands

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"leetcli/internal/api"
	"leetcli/internal/config"
	"leetcli/internal/template"
)

func TestProblem(args []string, cfg *config.Config, ui UI) {
	ui.WriteOutput(MsgPlain, "--- LeetCode Test ---\n")

	problemNum := ""
	if len(args) > 0 {
		problemNum = args[0]
	}
	if problemNum == "" {
		problemNum = ui.PromptText("Enter problem number to test")
	}
	if problemNum == "" {
		return
	}

	baseDir := cfg.BaseDir
	if baseDir == "" {
		ui.WriteOutput(MsgError, "Invalid base directory in config")
		return
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
		return
	}

	targetFile := matches[0]
	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, f := range matches {
			names[i] = filepath.Base(f)
		}
		selected := ui.PromptSelect("Multiple matches. Select file to test", names)
		for i, n := range names {
			if n == selected {
				targetFile = matches[i]
				break
			}
		}
	}

	contentBytes, err := os.ReadFile(targetFile)
	if err != nil {
		ui.WriteOutput(MsgError, "Failed to read file: %v", err)
		return
	}
	content := string(contentBytes)

	ext := filepath.Ext(targetFile)
	langKey := template.GetLanguageByExtension(languages, ext)
	if langKey == "" {
		ui.WriteOutput(MsgError, "Unsupported file extension for testing.")
		return
	}

	lcLang, ok := template.LeetCodeLangMap[langKey]
	if !ok {
		ui.WriteOutput(MsgError, "Language '%s' is not supported by LeetCode API.", langKey)
		return
	}

	slug := template.InferSlugFromPath(targetFile, content)
	if slug == "" {
		slug = ui.PromptText("Enter LeetCode problem slug (e.g., two-sum)")
		if slug == "" {
			return
		}
	}

	details, err := api.GetProblemDetails(slug)
	if err != nil {
		ui.WriteOutput(MsgError, "Could not fetch problem details from LeetCode: %v", err)
		return
	}

	session := cfg.GetLeetcodeSession()
	csrf := cfg.GetLeetcodeCsrf()
	if session == "" || csrf == "" {
		ui.WriteOutput(MsgError, "LeetCode test requires leetcode_session and leetcode_csrf.")
		ui.WriteOutput(MsgInfo, "Set them in config.local.json or use LEETCODE_SESSION / LEETCODE_CSRF env vars.")
		return
	}

	ui.WriteOutput(MsgInfo, "Fetching example testcases for slug: %s...", slug)
	testcases, err := api.GetProblemTestcases(slug)
	if err != nil || testcases == "" {
		ui.WriteOutput(MsgInfo, "Could not fetch sample testcases automatically.")
		testcases = ui.PromptText("Enter custom input testcases")
		if testcases == "" {
			return
		}
	} else {
		ui.WriteOutput(MsgSuccess, "Sample testcases fetched.")
	}

	ui.WriteOutput(MsgInfo, "Sending test run to LeetCode (%s)...", lcLang)
	interpretID, err := api.InterpretSolution(session, csrf, slug, details.QuestionID, lcLang, content, testcases)
	if err != nil {
		ui.WriteOutput(MsgError, "Test run request failed: %v", err)
		return
	}

	ui.WriteOutput(MsgInfo, "Test Run ID: %s. Polling for results...", interpretID)
	for i := 0; i < 15; i++ {
		time.Sleep(2 * time.Second)
		res, err := api.CheckSubmissionStatus(session, csrf, interpretID)
		if err != nil {
			ui.WriteOutput(MsgError, "Error checking status: %v", err)
			continue
		}
		if res.State == "PENDING" || res.State == "STARTED" {
			ui.WriteOutput(MsgInfo, "Status: %s...", res.State)
			continue
		}

		if res.StatusMsg == "Accepted" {
			ui.WriteOutput(MsgSuccess, "✔ TEST PASSED!")
			if res.StatusRuntime != "" {
				ui.WriteOutput(MsgPlain, "  Runtime: %s", res.StatusRuntime)
			}
			if res.StatusMemory != "" {
				ui.WriteOutput(MsgPlain, "  Memory:  %s", res.StatusMemory)
			}
			if len(res.CodeOutput) > 0 {
				ui.WriteOutput(MsgPlain, "  Your Output:     %s", strings.Join(res.CodeOutput, ", "))
			}
			if len(res.ExpectedOutput) > 0 {
				ui.WriteOutput(MsgPlain, "  Expected Output: %s", strings.Join(res.ExpectedOutput, ", "))
			}
			if len(res.StdOutput) > 0 {
				ui.WriteOutput(MsgInfo, "  Stdout Output:\n%s", strings.Join(res.StdOutput, "\n"))
			}
		} else {
			ui.WriteOutput(MsgError, "✘ TEST FAILED: %s", res.StatusMsg)
			if len(res.CodeOutput) > 0 {
				ui.WriteOutput(MsgPlain, "  Your Output:     %s", strings.Join(res.CodeOutput, ", "))
			}
			if len(res.ExpectedOutput) > 0 {
				ui.WriteOutput(MsgPlain, "  Expected Output: %s", strings.Join(res.ExpectedOutput, ", "))
			}
			if res.CompileError != "" || res.FullCompileError != "" {
				errMsg := res.CompileError
				if errMsg == "" {
					errMsg = res.FullCompileError
				}
				ui.WriteOutput(MsgError, "Compile Error:\n%s", errMsg)
			}
			if res.RuntimeError != "" || res.FullRuntimeError != "" {
				errMsg := res.RuntimeError
				if errMsg == "" {
					errMsg = res.FullRuntimeError
				}
				ui.WriteOutput(MsgError, "Runtime Error:\n%s", errMsg)
			}
		}
		return
	}

	ui.WriteOutput(MsgError, "Timed out waiting for test run result.")
}
