package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type CommandDoc struct {
	Name        string
	Summary     string
	Usage       string
	Description string
	Examples    []string
}

var CommandDocs = map[string]CommandDoc{
	"add": {
		Name:        "add",
		Summary:     "Create a new problem workspace file & README",
		Usage:       "add [problem_number]",
		Description: "Fetches problem details from LeetCode API, prompts for data structure category, creates the problem folder, generates solution template code and README.md description.",
		Examples:    []string{"add 1", "add 15"},
	},
	"add-sol": {
		Name:        "add-sol",
		Summary:     "Add a new solution method to an existing problem",
		Usage:       "add-sol [problem_number]",
		Description: "Allows adding an alternative solution approach (e.g., Solution 2 Two Pointers) with method title, time & space complexities, and source code snippet into the problem file.",
		Examples:    []string{"add-sol 1"},
	},
	"list": {
		Name:        "list",
		Summary:     "List and filter local problems",
		Usage:       "list",
		Description: "Lists all problem files created in your workspace. Supports filtering by data structure category (array, string, tree, etc.) and showing solution counts.",
		Examples:    []string{"list"},
	},
	"search": {
		Name:        "search",
		Summary:     "Search problems by keyword or problem ID",
		Usage:       "search [query]",
		Description: "Searches for problem titles or IDs matching the query string in your local workspace and LeetCode problem library.",
		Examples:    []string{"search two sum", "search 121"},
	},
	"manage-structures": {
		Name:        "manage-structures",
		Summary:     "Add, remove, or list data structure categories",
		Usage:       "manage-structures",
		Description: "Manages data structure category mappings (e.g. array -> array, tree -> tree) used to organize problem folders in your workspace.",
		Examples:    []string{"manage-structures"},
	},
	"stats": {
		Name:        "stats",
		Summary:     "View workspace solving statistics",
		Usage:       "stats",
		Description: "Displays local workspace statistics including total problems created, breakdown by category, and solution count metrics.",
		Examples:    []string{"stats"},
	},
	"theme": {
		Name:        "theme",
		Summary:     "Change TUI color theme",
		Usage:       "theme",
		Description: "Changes the color palette of the terminal user interface.",
		Examples:    []string{"theme"},
	},
	"daily": {
		Name:        "daily",
		Summary:     "Get and setup LeetCode Daily Challenge",
		Usage:       "daily",
		Description: "Fetches today's active LeetCode Daily Challenge question, displays title & difficulty, and automatically prompts to create local problem files.",
		Examples:    []string{"daily"},
	},
	"random": {
		Name:        "random",
		Summary:     "Get a random LeetCode problem",
		Usage:       "random",
		Description: "Fetches a random problem from LeetCode based on difficulty or category, displaying details and offering to create local files.",
		Examples:    []string{"random"},
	},
	"hint": {
		Name:        "hint",
		Summary:     "View hints for a specific problem",
		Usage:       "hint [problem_number]",
		Description: "Fetches official hints from LeetCode for the specified problem to help you solve it without looking at full solutions.",
		Examples:    []string{"hint 1"},
	},
	"similar": {
		Name:        "similar",
		Summary:     "Find similar problems on LeetCode",
		Usage:       "similar [problem_number]",
		Description: "Displays a list of related / similar LeetCode questions for a given problem number.",
		Examples:    []string{"similar 1"},
	},
	"open": {
		Name:        "open",
		Summary:     "Open problem file in code editor",
		Usage:       "open [problem_number]",
		Description: "Opens the local problem solution file and README in your configured code editor (VS Code, Vim, etc.).",
		Examples:    []string{"open 1"},
	},
	"run": {
		Name:        "run",
		Summary:     "Compile or run code locally",
		Usage:       "run [problem_number]",
		Description: "Compiles and executes your local solution code using installed system compilers/interpreters (g++, gcc, python, go, javac, rustc, node, tsx, dotnet).",
		Examples:    []string{"run 1"},
	},
	"test": {
		Name:        "test",
		Summary:     "Run sample testcases on LeetCode API",
		Usage:       "test [problem_number]",
		Description: "Fetches official example testcases for the problem (or accepts custom testcase input), sends code to LeetCode API, and displays Expected vs Actual output.",
		Examples:    []string{"test 1"},
	},
	"submit": {
		Name:        "submit",
		Summary:     "Submit solution to LeetCode API",
		Usage:       "submit [problem_number]",
		Description: "Submits your solution code directly to LeetCode API, polls real-time status, and displays Accepted / Wrong Answer / Error details along with Runtime and Memory percentiles (% Beats).",
		Examples:    []string{"submit 1"},
	},
	"sync": {
		Name:        "sync",
		Summary:     "Sync workspace with Git remote repository",
		Usage:       "sync",
		Description: "Automates Git workflow: stages all workspace changes, creates a timestamped commit, pulls remote changes with rebase, and pushes to GitHub.",
		Examples:    []string{"sync"},
	},
	"profile": {
		Name:        "profile",
		Summary:     "View LeetCode user profile stats",
		Usage:       "profile [username]",
		Description: "Fetches user profile information, global ranking, reputation, and solved question counts from LeetCode GraphQL API.",
		Examples:    []string{"profile vincent"},
	},
	"contest": {
		Name:        "contest",
		Summary:     "View upcoming LeetCode contests",
		Usage:       "contest",
		Description: "Displays upcoming LeetCode Weekly and Biweekly contests along with start times and durations.",
		Examples:    []string{"contest"},
	},
	"config": {
		Name:        "config",
		Summary:     "View and edit CLI settings (base_dir, editor, session, etc.)",
		Usage:       "config [show|open]",
		Description: "Opens an interactive menu to view and modify CLI settings (base workspace directory, default language, text editor, LeetCode cookies, theme) or open config.json directly.",
		Examples:    []string{"config", "config show", "config open", "cfg"},
	},
	"web": {
		Name:        "web",
		Summary:     "Open problem page in default web browser",
		Usage:       "web [problem_number]",
		Description: "Opens the official LeetCode problem page (https://leetcode.com/problems/slug/) directly in your default web browser (Chrome, Edge, Safari...).",
		Examples:    []string{"web 1", "browser two-sum"},
	},
	"timer": {
		Name:        "timer",
		Summary:     "Start practice countdown stopwatch for interview practice",
		Usage:       "timer [minutes]",
		Description: "Starts a background countdown timer (e.g. 20, 30, 45 minutes) to practice solving problems under real interview time pressure.",
		Examples:    []string{"timer 20", "timer 30"},
	},
	"note": {
		Name:        "note",
		Summary:     "Add study notes and key takeaways to problem README",
		Usage:       "note [problem_number]",
		Description: "Appends personal study notes, key tricks, and complexity takeaways into the problem's README.md file.",
		Examples:    []string{"note 1"},
	},
	"clean": {
		Name:        "clean",
		Summary:     "Clean temporary build artifacts (.exe, .class, temp files)",
		Usage:       "clean",
		Description: "Scans workspace recursively and removes temporary build files generated during local execution (leet_run_*, .class, .tmp) to keep repository clean.",
		Examples:    []string{"clean"},
	},
	"clear": {
		Name:        "clear",
		Summary:     "Clear TUI terminal output",
		Usage:       "clear",
		Description: "Clears the scrollback history of the terminal interface.",
		Examples:    []string{"clear"},
	},
	"exit": {
		Name:        "exit",
		Summary:     "Exit the CLI application",
		Usage:       "exit",
		Description: "Saves current configuration and cleanly exits the application.",
		Examples:    []string{"exit", "quit"},
	},
}

func RenderCommandManual(cmdName string) string {
	doc, ok := CommandDocs[cmdName]
	if !ok {
		return ErrorStyle.Render(fmt.Sprintf("No detailed manual found for command '%s'.", cmdName)) + "\n"
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(Cyan)
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(Yellow)
	valStyle := lipgloss.NewStyle().Foreground(White)
	codeStyle := lipgloss.NewStyle().Foreground(Green)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render(fmt.Sprintf("📖 MANUAL: %s", strings.ToUpper(doc.Name))))
	sb.WriteString("\n")
	sb.WriteString(DimmedStyle.Render(strings.Repeat("─", 50)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Summary:    "), valStyle.Render(doc.Summary)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Usage:      "), codeStyle.Render(doc.Usage)))
	sb.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("Description:"), valStyle.Render(doc.Description)))
	if len(doc.Examples) > 0 {
		sb.WriteString(labelStyle.Render("Examples:   \n"))
		for _, ex := range doc.Examples {
			sb.WriteString(fmt.Sprintf("  • %s\n", codeStyle.Render(ex)))
		}
	}
	sb.WriteString(DimmedStyle.Render(strings.Repeat("─", 50)))
	sb.WriteString("\n")
	return sb.String()
}
