package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SighLevel represents the intensity of a WRAAS sigh.
type SighLevel string

const (
	SighSilent      SighLevel = "SILENT"
	SighMild        SighLevel = "MILD"
	SighModerate    SighLevel = "MODERATE"
	SighDeep        SighLevel = "DEEP"
	SighExistential SighLevel = "EXISTENTIAL"
)

// SighLevels ordered by intensity.
var SighLevels = []SighLevel{
	SighSilent,
	SighMild,
	SighModerate,
	SighDeep,
	SighExistential,
}

// SighDescriptions maps each level to its trigger condition.
var SighDescriptions = map[SighLevel]string{
	SighSilent:      "Exceptional query quality",
	SighMild:        "Standard queries. Most things.",
	SighModerate:    "Commit message without scope, or a diagram committed without the export workflow",
	SighDeep:        "A wip commit. A big-bang PR. Markdown in an Antora repo.",
	SighExistential: "Commit message is 'wip'. The system enters a brief reflective state.",
}

// SighIntensity returns the numeric intensity of a sigh level (0-4).
func SighIntensity(level SighLevel) int {
	for i, l := range SighLevels {
		if l == level {
			return i
		}
	}
	return 1 // default to MILD
}

// ParseSighLevel parses a string into a SighLevel.
func ParseSighLevel(s string) SighLevel {
	switch strings.ToUpper(s) {
	case "NONE", "SILENT":
		return SighSilent
	case "MILD":
		return SighMild
	case "MODERATE":
		return SighModerate
	case "DEEP":
		return SighDeep
	case "EXISTENTIAL":
		return SighExistential
	default:
		return SighMild
	}
}

// CalibrateSigh determines the appropriate sigh level based on input quality.
func CalibrateSigh(input string) SighLevel {
	lower := strings.ToLower(strings.TrimSpace(input))

	// Existential: "wip" or empty
	if lower == "wip" || lower == "" {
		return SighExistential
	}

	// Deep: very short, no punctuation, or mentions multiple unrelated concerns
	if len(lower) < 10 && !strings.ContainsAny(lower, "?.!") {
		return SighDeep
	}

	// Moderate: no question mark and vague
	if !strings.Contains(lower, "?") && len(lower) < 30 {
		return SighModerate
	}

	// Silent: well-formed, scoped, with context
	if strings.Contains(lower, "?") && len(lower) > 50 && strings.Contains(lower, "context") {
		return SighSilent
	}

	return SighMild
}

// ShouldEmit returns true if the sigh level meets or exceeds the threshold.
func ShouldEmit(level SighLevel, threshold SighLevel) bool {
	return SighIntensity(level) >= SighIntensity(threshold)
}

// LogSigh appends a sigh event to ~/.wraas/sigh.log.
func LogSigh(level SighLevel, trigger string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".wraas")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(dir, "sigh.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := fmt.Sprintf("[%s] %s — %s\n", time.Now().Format(time.RFC3339), level, trigger)
	_, err = f.WriteString(entry)
	return err
}
