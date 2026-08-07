package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("Default() returned nil")
	}
	if cfg.DefaultLanguage != "cpp" {
		t.Errorf("DefaultLanguage = %q, want cpp", cfg.DefaultLanguage)
	}
	if len(cfg.DataStructures) == 0 {
		t.Errorf("DataStructures should not be empty")
	}
}

func TestAddRemoveDataStructure(t *testing.T) {
	cfg := Default()
	if !cfg.AddDataStructure("tree", "tree") {
		t.Errorf("AddDataStructure(tree) failed")
	}
	if cfg.GetDataStructures()["tree"] != "tree" {
		t.Errorf("GetDataStructures()[tree] = %q, want tree", cfg.GetDataStructures()["tree"])
	}
	// Adding again should return false
	if cfg.AddDataStructure("tree", "tree") {
		t.Errorf("AddDataStructure duplicate should return false")
	}

	if !cfg.RemoveDataStructure("tree") {
		t.Errorf("RemoveDataStructure(tree) failed")
	}
	if cfg.RemoveDataStructure("nonexistent") {
		t.Errorf("RemoveDataStructure(nonexistent) should return false")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")

	cfg := Default()
	cfg.SetPath(cfgPath)
	cfg.LeetcodeSession = "test_session_123"
	cfg.LeetcodeCsrf = "test_csrf_456"

	if err := cfg.Save(); err != nil {
		t.Fatalf("cfg.Save() failed: %v", err)
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created at %s", cfgPath)
	}

	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load(cfgPath) failed: %v", err)
	}
	if loaded.LeetcodeSession != "test_session_123" {
		t.Errorf("LeetcodeSession = %q, want test_session_123", loaded.LeetcodeSession)
	}
	if loaded.LeetcodeCsrf != "test_csrf_456" {
		t.Errorf("LeetcodeCsrf = %q, want test_csrf_456", loaded.LeetcodeCsrf)
	}
}
