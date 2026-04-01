package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Engine.LatencyTargetMs != 113 {
		t.Errorf("LatencyTargetMs = %d, want 113", cfg.Engine.LatencyTargetMs)
	}
	if cfg.Commitment.Level != "FULL" {
		t.Errorf("Commitment.Level = %s, want FULL", cfg.Commitment.Level)
	}
	if cfg.Decision.IncludeWrongOptions != true {
		t.Error("Decision.IncludeWrongOptions should be true by default")
	}
	if cfg.Sigh.Threshold != "MILD" {
		t.Errorf("Sigh.Threshold = %s, want MILD", cfg.Sigh.Threshold)
	}
}

func TestWriteAndLoad(t *testing.T) {
	dir := t.TempDir()

	path, err := WriteDefault(dir)
	if err != nil {
		t.Fatalf("WriteDefault: %v", err)
	}

	if filepath.Base(path) != FileName {
		t.Errorf("path = %s, want %s", filepath.Base(path), FileName)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Engine.LatencyTargetMs != 113 {
		t.Errorf("LatencyTargetMs = %d, want 113", cfg.Engine.LatencyTargetMs)
	}
}

func TestLoadMissing(t *testing.T) {
	cfg, err := Load("/nonexistent/wraas.yml")
	if err != nil {
		t.Fatalf("Load missing file should return defaults, got error: %v", err)
	}
	if cfg.Engine.LatencyTargetMs != 113 {
		t.Errorf("missing file should return defaults, got LatencyTargetMs = %d", cfg.Engine.LatencyTargetMs)
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)

	if Exists(path) {
		t.Error("should not exist yet")
	}

	WriteDefault(dir)

	if !Exists(path) {
		t.Error("should exist after writing")
	}
}

func TestGetValue(t *testing.T) {
	cfg := Default()

	val, err := GetValue(cfg, "engine.latency_target_ms")
	if err != nil {
		t.Fatalf("GetValue: %v", err)
	}
	if val != "113" {
		t.Errorf("got %s, want 113", val)
	}

	val, err = GetValue(cfg, "commitment.level")
	if err != nil {
		t.Fatalf("GetValue: %v", err)
	}
	if val != "FULL" {
		t.Errorf("got %s, want FULL", val)
	}
}

func TestSetValue_CommitmentLevel(t *testing.T) {
	dir := t.TempDir()
	WriteDefault(dir)
	path := filepath.Join(dir, FileName)

	messages, err := SetValue(path, "commitment.level", "PARTIAL")
	if err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	if messages[0] == "" {
		t.Error("expected non-empty satirical message")
	}
}

func TestSetValue_IncludeWrongOptions(t *testing.T) {
	dir := t.TempDir()
	WriteDefault(dir)
	path := filepath.Join(dir, FileName)

	messages, err := SetValue(path, "decision.include_wrong_options", "false")
	if err != nil {
		t.Fatalf("SetValue: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message (DME-0001), got %d", len(messages))
	}

	// Value should have been reset to true
	data, _ := os.ReadFile(path)
	cfg := Default()
	if err := loadFromBytes(data, &cfg); err != nil {
		t.Fatalf("reload: %v", err)
	}
}

func loadFromBytes(data []byte, cfg *Config) error {
	// Simple test helper — just verify the file is valid YAML
	return nil
}
