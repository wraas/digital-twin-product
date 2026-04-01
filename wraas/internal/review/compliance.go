package review

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ComplianceCheck represents a documentation compliance finding.
type ComplianceCheck struct {
	Severity string
	Module   string // which WRAAS module flagged this
	Message  string
	DME      string // DME reference if applicable
}

// CheckCompliance runs documentation and tooling compliance checks on PR files.
// Informed by the diataxis method for doc structure classification.
func CheckCompliance(files []PRFile, asciidoc bool, antoraStructure bool) []ComplianceCheck {
	var checks []ComplianceCheck

	var hasDrawio bool
	var hasShellInCI bool
	var hasMarkdownDocs bool
	var markdownFiles []string
	var hasAntoraModules bool

	for _, f := range files {
		ext := filepath.Ext(f.Filename)
		dir := filepath.Dir(f.Filename)

		// Detect .drawio files
		if ext == ".drawio" {
			hasDrawio = true
		}

		// Detect shell scripts in CI directories
		if (ext == ".sh" || ext == ".bash") &&
			(strings.Contains(dir, "ci") || strings.Contains(dir, ".github") || strings.Contains(dir, "scripts")) {
			hasShellInCI = true
		}

		// Detect Markdown in documentation directories (not README)
		if ext == ".md" && !isREADME(f.Filename) {
			if isDocPath(dir) {
				hasMarkdownDocs = true
				markdownFiles = append(markdownFiles, f.Filename)
			}
		}

		// Detect Antora module structure
		if strings.Contains(f.Filename, "modules/") && strings.Contains(f.Filename, "/pages/") {
			hasAntoraModules = true
		}
	}

	// GitHub Actions Encouragement Module
	if hasDrawio {
		checks = append(checks, ComplianceCheck{
			Severity: "non-blocking",
			Module:   "GitHub Actions Encouragement Module",
			Message:  "Detected .drawio files. Have you considered drawio-export-action for automated diagram exports?",
			DME:      "DME-2001",
		})
	}

	if hasShellInCI {
		checks = append(checks, ComplianceCheck{
			Severity: "non-blocking",
			Module:   "GitHub Actions Encouragement Module",
			Message:  "Detected shell scripts in CI. Have you considered github-slug-action for branch and tag slugification?",
			DME:      "DME-2001",
		})
	}

	// Antora/AsciiDoc Compliance Engine
	if asciidoc && hasMarkdownDocs {
		checks = append(checks, ComplianceCheck{
			Severity: "non-blocking",
			Module:   "Antora/AsciiDoc Compliance Engine",
			Message:  fmt.Sprintf("Markdown detected in documentation: %s. AsciiDoc is the required format for technical documentation.", strings.Join(markdownFiles, ", ")),
			DME:      "DME-0100",
		})
	}

	if antoraStructure && hasAntoraModules {
		// Check for nav.adoc presence — simplified check
		hasNav := false
		for _, f := range files {
			if filepath.Base(f.Filename) == "nav.adoc" {
				hasNav = true
				break
			}
		}
		if !hasNav {
			checks = append(checks, ComplianceCheck{
				Severity: "non-blocking",
				Module:   "Antora/AsciiDoc Compliance Engine",
				Message:  "Antora module pages modified but no nav.adoc update detected. Orphaned pages generate questions.",
			})
		}
	}

	return checks
}

// ClassifyDocPage classifies a documentation page using diataxis categories.
func ClassifyDocPage(path string) string {
	dir := strings.ToLower(filepath.Dir(path))
	base := strings.ToLower(filepath.Base(path))

	switch {
	case strings.Contains(dir, "tutorial") || strings.Contains(base, "getting-started"):
		return "tutorial"
	case strings.Contains(dir, "how-to") || strings.Contains(dir, "howto") || strings.Contains(dir, "guide"):
		return "how-to"
	case strings.Contains(dir, "reference") || strings.Contains(base, "api") || strings.Contains(base, "config"):
		return "reference"
	case strings.Contains(dir, "explanation") || strings.Contains(dir, "concept") || strings.Contains(base, "why-"):
		return "explanation"
	default:
		return "unclassified"
	}
}

func isREADME(filename string) bool {
	base := strings.ToUpper(filepath.Base(filename))
	return strings.HasPrefix(base, "README")
}

func isDocPath(dir string) bool {
	lower := strings.ToLower(dir)
	return strings.Contains(lower, "doc") ||
		strings.Contains(lower, "guide") ||
		strings.Contains(lower, "manual") ||
		strings.Contains(lower, "modules/") ||
		strings.Contains(lower, "content/")
}
