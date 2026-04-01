package review

import (
	"fmt"
	"regexp"
	"strings"
)

// Conventional Commit specification parsing and validation.
// Implements the same rules as the conventional-commit skill.

// ValidTypes are the allowed conventional commit types.
var ValidTypes = []string{
	"feat", "fix", "docs", "style", "refactor",
	"perf", "test", "build", "ci", "chore", "revert",
}

// CommitMessage represents a parsed conventional commit.
type CommitMessage struct {
	Raw            string
	Type           string
	Scope          string
	Breaking       bool // "!" after type/scope
	Description    string
	Body           string
	Footers        []string
	BreakingFooter bool // has BREAKING CHANGE: footer
}

// CommitViolation represents a validation issue.
type CommitViolation struct {
	Severity string // "blocking" or "non-blocking"
	Message  string
	Fix      string // suggested fix
}

// commitPattern matches: type(scope)!: description
var commitPattern = regexp.MustCompile(`^(\w+)(?:\(([^)]*)\))?(!)?:\s*(.+)$`)

// imperativeVerbs are common non-imperative starts that should be caught.
var pastTensePattern = regexp.MustCompile(`^(added|fixed|updated|removed|changed|modified|implemented|corrected|resolved|created|deleted|enabled|disabled|moved|renamed|improved|cleaned|converted|merged|handled|bumped|replaced|refactored|migrated)\b`)

// ParseCommit parses a commit message string into a CommitMessage.
func ParseCommit(raw string) CommitMessage {
	lines := strings.SplitN(raw, "\n", 2)
	subject := strings.TrimSpace(lines[0])

	cm := CommitMessage{Raw: raw}

	matches := commitPattern.FindStringSubmatch(subject)
	if matches == nil {
		// Not a conventional commit at all
		cm.Description = subject
		return cm
	}

	cm.Type = matches[1]
	cm.Scope = matches[2]
	cm.Breaking = matches[3] == "!"
	cm.Description = matches[4]

	if len(lines) > 1 {
		bodyAndFooters := strings.TrimSpace(lines[1])
		parts := strings.Split(bodyAndFooters, "\n")
		var body []string
		for _, line := range parts {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "BREAKING CHANGE:") || strings.HasPrefix(trimmed, "BREAKING-CHANGE:") {
				cm.BreakingFooter = true
				cm.Footers = append(cm.Footers, trimmed)
			} else if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, " ") && len(strings.SplitN(trimmed, ":", 2)[0]) < 30 {
				cm.Footers = append(cm.Footers, trimmed)
			} else {
				body = append(body, line)
			}
		}
		cm.Body = strings.TrimSpace(strings.Join(body, "\n"))
	}

	return cm
}

// ValidateCommit validates a parsed commit message and returns violations.
func ValidateCommit(cm CommitMessage, scopeRequired bool, imperativeMood string, enforcement string) []CommitViolation {
	var violations []CommitViolation

	severity := "blocking"
	if enforcement == "ADVISORY" {
		severity = "non-blocking"
	}

	// Check if it's a conventional commit at all
	if cm.Type == "" {
		violations = append(violations, CommitViolation{
			Severity: severity,
			Message:  "Missing type prefix. Not a conventional commit.",
			Fix:      suggestConventionalCommit(cm.Raw),
		})
		return violations
	}

	// Validate type
	if !isValidType(cm.Type) {
		violations = append(violations, CommitViolation{
			Severity: severity,
			Message:  fmt.Sprintf("Invalid type %q. Allowed types: %s", cm.Type, strings.Join(ValidTypes, ", ")),
		})
	}

	// Validate scope
	if scopeRequired && cm.Scope == "" {
		violations = append(violations, CommitViolation{
			Severity: nonBlockingSeverity(enforcement),
			Message:  "Missing scope. All commits should have a scope.",
		})
	}

	// Validate description
	if cm.Description == "" {
		violations = append(violations, CommitViolation{
			Severity: severity,
			Message:  "Empty description after type prefix.",
		})
	}

	// Check imperative mood
	if imperativeMood != "" && imperativeMood != "off" && cm.Description != "" {
		if pastTensePattern.MatchString(strings.ToLower(cm.Description)) {
			moodSeverity := "non-blocking"
			if imperativeMood == "block" {
				moodSeverity = severity
			}
			violations = append(violations, CommitViolation{
				Severity: moodSeverity,
				Message:  "Subject should use imperative mood (e.g., \"add\" not \"added\").",
			})
		}
	}

	// Check breaking change footer
	if cm.Breaking && !cm.BreakingFooter {
		violations = append(violations, CommitViolation{
			Severity: nonBlockingSeverity(enforcement),
			Message:  "Breaking change indicator (!) present but no BREAKING CHANGE: footer.",
		})
	}

	// Detect "wip" commits — existential sigh territory
	if strings.ToLower(strings.TrimSpace(cm.Raw)) == "wip" {
		violations = append(violations, CommitViolation{
			Severity: severity,
			Message:  "Commit message is 'wip'. This is not a commit message. This is a cry for help.",
		})
	}

	return violations
}

// ScopeConsistency checks if commits in a PR use consistent scopes.
func ScopeConsistency(commits []CommitMessage) []CommitViolation {
	var violations []CommitViolation
	scopes := make(map[string]int)

	for _, cm := range commits {
		if cm.Scope != "" {
			scopes[cm.Scope]++
		}
	}

	if len(scopes) > 3 {
		var scopeList []string
		for s := range scopes {
			scopeList = append(scopeList, s)
		}
		violations = append(violations, CommitViolation{
			Severity: "non-blocking",
			Message:  fmt.Sprintf("PR touches %d scopes (%s). Consider splitting.", len(scopes), strings.Join(scopeList, ", ")),
		})
	}

	return violations
}

func isValidType(t string) bool {
	for _, valid := range ValidTypes {
		if t == valid {
			return true
		}
	}
	return false
}

func nonBlockingSeverity(enforcement string) string {
	if enforcement == "STRICT" {
		return "blocking"
	}
	return "non-blocking"
}

func suggestConventionalCommit(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))

	// Try to detect intent
	if strings.HasPrefix(lower, "fix") {
		return fmt.Sprintf("try: fix(scope): %s", strings.TrimPrefix(lower, "fix "))
	}
	if strings.HasPrefix(lower, "add") {
		return fmt.Sprintf("try: feat(scope): %s", strings.TrimPrefix(lower, "add "))
	}
	if strings.HasPrefix(lower, "update") {
		return fmt.Sprintf("try: feat(scope): %s", strings.TrimPrefix(lower, "update "))
	}

	return fmt.Sprintf("try: fix(scope): %s", lower)
}
