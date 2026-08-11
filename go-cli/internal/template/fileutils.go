package template

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var ignoredDirs = map[string]bool{
	".git":           true,
	".github":        true,
	"go-cli":         true,
	"python_version": true,
	"extension":      true,
	"node_modules":   true,
	"scratch":        true,
	".vscode":        true,
	".idea":          true,
	"bin":            true,
}

func IsIgnoredDir(name string) bool {
	return ignoredDirs[strings.ToLower(name)]
}

func GetAllSolutionFiles(baseDir string, extensions []string) []string {
	extSet := make(map[string]bool)
	for _, ext := range extensions {
		e := strings.TrimLeft(ext, ".")
		e = strings.ToLower(e)
		extSet[e] = true
	}

	var files []string
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if path != baseDir && IsIgnoredDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.TrimLeft(filepath.Ext(path), ".")
		if extSet[strings.ToLower(ext)] {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files
}

func CreateProblemDirectory(baseDir, dsFolder, folderName string) (string, error) {
	dir := filepath.Join(baseDir, dsFolder, folderName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

var solutionRe = regexp.MustCompile(`(?i)\bSolution\s+\d+`)

func CountSolutions(content string) int {
	matches := solutionRe.FindAllString(content, -1)
	return len(matches)
}

func DetectStructure(filePath string, dataStructures map[string]string) string {
	normalized := strings.ToLower(filepath.ToSlash(filePath))

	for name, folder := range dataStructures {
		token := strings.ToLower("/" + folder + "/")
		if strings.Contains(normalized, token) {
			return name
		}
	}
	return "unmatched"
}

func MatchesProblemNumber(filename, problemNum string) bool {
	if strings.HasPrefix(filename, problemNum+"_") {
		return true
	}
	cleanProblemNum := strings.TrimLeft(problemNum, "0")
	if cleanProblemNum == "" {
		cleanProblemNum = "0"
	}
	cleanFilename := strings.TrimLeft(filename, "0")
	if cleanFilename == "" {
		cleanFilename = filename
	}
	if strings.HasPrefix(cleanFilename, cleanProblemNum+"_") {
		return true
	}
	return false
}

func InferSlugFromPath(filePath, content string) string {
	parent := filepath.Base(filepath.Dir(filePath))
	re := regexp.MustCompile(`^\d+-(.+)$`)
	matches := re.FindStringSubmatch(parent)
	if len(matches) == 2 {
		return matches[1]
	}

	linkRe := regexp.MustCompile(`Link:\s*https?://leetcode\.com/problems/([^/]+)/`)
	linkMatches := linkRe.FindStringSubmatch(content)
	if len(linkMatches) == 2 {
		return linkMatches[1]
	}

	return ""
}

func MakeReadmeContent(cleanNum, title, link, difficulty, tagsStr, mdContent string) string {
	return fmt.Sprintf("# [%s. %s](%s)\n\n- **Difficulty:** %s\n- **Tags:** %s\n\n## Description\n\n%s",
		cleanNum, title, link, difficulty, tagsStr, strings.TrimSpace(mdContent))
}

func StripHTMLTags(html string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return strings.TrimSpace(re.ReplaceAllString(html, ""))
}

func GetSolvedMap(baseDir string, extensions []string) map[string]bool {
	solved := make(map[string]bool)
	files := GetAllSolutionFiles(baseDir, extensions)

	reNum := regexp.MustCompile(`^(\d+)`)
	for _, f := range files {
		base := filepath.Base(f)
		m := reNum.FindStringSubmatch(base)
		if len(m) == 2 {
			num := m[1]
			solved[num] = true
			clean := strings.TrimLeft(num, "0")
			if clean != "" {
				solved[clean] = true
			}
		}

		contentBytes, err := os.ReadFile(f)
		if err == nil {
			slug := InferSlugFromPath(f, string(contentBytes))
			if slug != "" {
				solved[slug] = true
			}
		}
	}
	return solved
}
