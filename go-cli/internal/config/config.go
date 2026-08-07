package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type LangEntry struct {
	Label string `json:"label"`
	Ext   string `json:"ext"`
}

type Config struct {
	BaseDir         string                `json:"base_dir"`
	Languages       map[string]LangEntry  `json:"languages"`
	DefaultLanguage string                `json:"default_language"`
	DataStructures  map[string]string     `json:"data_structures"`
	Theme           string                `json:"theme"`
	Editor          string                `json:"editor"`
	LeetcodeSession string                `json:"leetcode_session,omitempty"`
	LeetcodeCsrf    string                `json:"leetcode_csrf,omitempty"`

	path string
}

func Load(path string) (*Config, error) {
	if path == "" {
		possible := []string{
			"config.json",
			filepath.Join("..", "config.json"),
			filepath.Join("..", "..", "config.json"),
		}
		for _, p := range possible {
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
	}
	if path == "" {
		return Default(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), nil
	}
	cfg.path, _ = filepath.Abs(path)
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		BaseDir:         ".",
		DefaultLanguage: "cpp",
		Theme:           "default",
		Editor:          "code",
		Languages: map[string]LangEntry{
			"cpp":    {Label: "C++", Ext: "cpp"},
			"c":      {Label: "C", Ext: "c"},
			"python": {Label: "Python", Ext: "py"},
			"go":     {Label: "Go", Ext: "go"},
		},
		DataStructures: map[string]string{
			"array":      "array",
			"string":     "string",
			"linkedlist": "linked_list",
			"stack":      "stack",
			"graph":      "graph",
			"queue":      "queue",
		},
	}
}

func (c *Config) GetPath() string {
	if c.path == "" {
		return "config.json"
	}
	return c.path
}

func (c *Config) SetPath(p string) {
	if p == "" {
		p = "config.json"
	}
	abs, err := filepath.Abs(p)
	if err == nil {
		c.path = abs
	} else {
		c.path = p
	}
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	p := c.GetPath()
	return os.WriteFile(p, data, 0644)
}

func (c *Config) GetDataStructures() map[string]string {
	if c.DataStructures == nil {
		return map[string]string{}
	}
	return c.DataStructures
}

func (c *Config) AddDataStructure(name, folder string) bool {
	if c.DataStructures == nil {
		c.DataStructures = make(map[string]string)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}
	if _, ok := c.DataStructures[name]; ok {
		return false
	}
	c.DataStructures[name] = folder
	return true
}

func (c *Config) RemoveDataStructure(name string) bool {
	if _, ok := c.DataStructures[name]; !ok {
		return false
	}
	delete(c.DataStructures, name)
	return true
}

func (c *Config) SetTheme(theme string) {
	c.Theme = theme
}

func (c *Config) GetTheme() string {
	if c.Theme == "" {
		return "default"
	}
	return c.Theme
}

func (c *Config) GetEditor() string {
	if c.Editor == "" {
		return "code"
	}
	return c.Editor
}

func (c *Config) SetEditor(editor string) {
	c.Editor = editor
}
