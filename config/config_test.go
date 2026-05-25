package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Name    string `yaml:"name" json:"name"`
	Version int    `yaml:"version" json:"version"`
}

func TestLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	content := []byte("name: planx\nversion: 4\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var cfg testConfig
	if err := LoadYAML(path, &cfg); err != nil {
		t.Fatalf("LoadYAML: %v", err)
	}
	if cfg.Name != "planx" || cfg.Version != 4 {
		t.Fatalf("got %+v", cfg)
	}
}

func TestLoadYAML_FileNotFound(t *testing.T) {
	var cfg testConfig
	err := LoadYAML("/nonexistent/path.yaml", &cfg)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadYAML_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("::not{valid"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var cfg testConfig
	err := LoadYAML(path, &cfg)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	content := []byte(`{"name":"planx","version":4}`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var cfg testConfig
	if err := LoadJSON(path, &cfg); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	if cfg.Name != "planx" || cfg.Version != 4 {
		t.Fatalf("got %+v", cfg)
	}
}

func TestLoadJSON_FileNotFound(t *testing.T) {
	var cfg testConfig
	err := LoadJSON("/nonexistent/path.json", &cfg)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadJSON_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{invalid}"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var cfg testConfig
	err := LoadJSON(path, &cfg)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte("name: planx\nversion: 4\n")
	var cfg testConfig
	if err := ParseYAML(data, &cfg); err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if cfg.Name != "planx" || cfg.Version != 4 {
		t.Fatalf("got %+v", cfg)
	}
}

func TestParseYAML_Invalid(t *testing.T) {
	var cfg testConfig
	err := ParseYAML([]byte("::not{valid"), &cfg)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseJSON(t *testing.T) {
	data := []byte(`{"name":"planx","version":4}`)
	var cfg testConfig
	if err := ParseJSON(data, &cfg); err != nil {
		t.Fatalf("ParseJSON: %v", err)
	}
	if cfg.Name != "planx" || cfg.Version != 4 {
		t.Fatalf("got %+v", cfg)
	}
}

func TestParseJSON_Invalid(t *testing.T) {
	var cfg testConfig
	err := ParseJSON([]byte("{bad}"), &cfg)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadJSON_Intomap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "m.json")
	content := []byte(`{"key":"val"}`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var m map[string]string
	if err := LoadJSON(path, &m); err != nil {
		t.Fatalf("LoadJSON into map: %v", err)
	}
	if m["key"] != "val" {
		t.Fatalf("got %v", m)
	}
}

func TestParseYAML_Intoslice(t *testing.T) {
	data := []byte("- a\n- b\n- c\n")
	var s []string
	if err := ParseYAML(data, &s); err != nil {
		t.Fatalf("ParseYAML into slice: %v", err)
	}
	if len(s) != 3 || s[0] != "a" || s[2] != "c" {
		t.Fatalf("got %v", s)
	}
}

func TestParseJSON_RawMessage(t *testing.T) {
	data := []byte(`{"nested":{"k":1}}`)
	var raw json.RawMessage
	if err := ParseJSON(data, &raw); err != nil {
		t.Fatalf("ParseJSON RawMessage: %v", err)
	}
	if string(raw) != `{"nested":{"k":1}}` {
		t.Fatalf("got %s", raw)
	}
}
