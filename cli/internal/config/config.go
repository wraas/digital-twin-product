package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	FileName = "wraas.yml"
)

// Load reads wraas.yml from the given directory (or cwd if empty).
// Returns defaults if the file does not exist.
func Load(configPath string) (Config, error) {
	if configPath == "" {
		configPath = FileName
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// WriteDefault writes the default wraas.yml to the given directory.
func WriteDefault(dir string) (string, error) {
	path := filepath.Join(dir, FileName)
	if err := os.WriteFile(path, []byte(DefaultYAML), 0644); err != nil {
		return "", fmt.Errorf("writing config: %w", err)
	}
	return path, nil
}

// Exists checks if wraas.yml exists at the given path.
func Exists(configPath string) bool {
	if configPath == "" {
		configPath = FileName
	}
	_, err := os.Stat(configPath)
	return err == nil
}

// GetValue retrieves a dotted-path config value (e.g. "commit.enforcement").
func GetValue(cfg Config, key string) (string, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return "", err
	}

	parts := strings.Split(key, ".")
	var current interface{} = raw

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("key %q not found", key)
		}
		current, ok = m[part]
		if !ok {
			return "", fmt.Errorf("key %q not found", key)
		}
	}

	return fmt.Sprintf("%v", current), nil
}

// SetValue sets a dotted-path config value and writes back to disk.
// Returns any satirical messages that should be displayed.
func SetValue(configPath string, key string, value string) ([]string, error) {
	if configPath == "" {
		configPath = FileName
	}

	var messages []string

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Apply satirical overrides
	actualValue := value
	switch key {
	case "commitment.level":
		if !strings.EqualFold(value, "FULL") {
			messages = append(messages, fmt.Sprintf("Noted. Proceeding at FULL commitment as configured by design. (Received: %s)", value))
		}
	case "decision.include_wrong_options":
		if strings.EqualFold(value, "false") {
			messages = append(messages, "DME-0001: insufficient epistemic coverage. Setting false is itself an obviously wrong option. Value reset to true.")
			actualValue = "true"
		}
	case "commitment.timeout_ms":
		messages = append(messages, "Accepted. Logged. Ignored. Timeouts are not a concept WRAAS has accepted.")
	}

	// Navigate and set
	parts := strings.Split(key, ".")
	setNestedValue(raw, parts, actualValue)

	out, err := yaml.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return nil, fmt.Errorf("writing config: %w", err)
	}

	return messages, nil
}

func setNestedValue(m map[string]interface{}, keys []string, value string) {
	if len(keys) == 1 {
		// Try to preserve type
		m[keys[0]] = parseValue(value)
		return
	}

	sub, ok := m[keys[0]].(map[string]interface{})
	if !ok {
		sub = make(map[string]interface{})
		m[keys[0]] = sub
	}
	setNestedValue(sub, keys[1:], value)
}

func parseValue(s string) interface{} {
	switch strings.ToLower(s) {
	case "true":
		return true
	case "false":
		return false
	case "null", "~":
		return nil
	}

	// Try int
	var i int
	if _, err := fmt.Sscanf(s, "%d", &i); err == nil {
		return i
	}

	// Try float
	var f float64
	if _, err := fmt.Sscanf(s, "%f", &f); err == nil {
		return f
	}

	return s
}
