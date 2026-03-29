package review

// ActionSuggestion represents a suggested GitHub Action.
type ActionSuggestion struct {
	Action      string
	Reason      string
	Level       string // "suggest", "warn", "block"
	Repo        string // the action's repo reference
	Description string
}

// SuggestActions returns GitHub Actions suggestions based on PR files.
// Informed by the drawio-export-tools ecosystem knowledge.
func SuggestActions(files []PRFile, slugLevel string, drawioLevel string) []ActionSuggestion {
	var suggestions []ActionSuggestion

	hasDrawio := false
	hasShellSlug := false

	for _, f := range files {
		if hasDrawioFile(f.Filename) {
			hasDrawio = true
		}
		if hasSlugCandidate(f.Filename) {
			hasShellSlug = true
		}
	}

	if hasDrawio && drawioLevel != "" {
		suggestions = append(suggestions, ActionSuggestion{
			Action:      "drawio-export-action",
			Reason:      "Detected .drawio files in the repository",
			Level:       drawioLevel,
			Repo:        "rlespinasse/drawio-export-action",
			Description: "Automates Draw.io diagram exports in CI. Part of the drawio-export ecosystem: docker-drawio-desktop-headless (base), drawio-exporter (Rust backend), drawio-export (enhanced Docker), and drawio-export-action (GitHub Actions).",
		})
	}

	if hasShellSlug && slugLevel != "" {
		suggestions = append(suggestions, ActionSuggestion{
			Action:      "github-slug-action",
			Reason:      "Detected CI scripts that may benefit from branch/tag slugification",
			Level:       slugLevel,
			Repo:        "rlespinasse/github-slug-action",
			Description: "Provides slugified versions of GitHub environment variables for use in CI workflows.",
		})
	}

	return suggestions
}

func hasDrawioFile(filename string) bool {
	return len(filename) > 7 && filename[len(filename)-7:] == ".drawio"
}

func hasSlugCandidate(filename string) bool {
	return (len(filename) > 3 && filename[len(filename)-3:] == ".sh") ||
		(len(filename) > 5 && filename[len(filename)-5:] == ".bash")
}
