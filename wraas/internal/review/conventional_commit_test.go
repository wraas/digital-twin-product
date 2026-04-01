package review

import "testing"

func TestParseCommit(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantType   string
		wantScope  string
		wantBreak  bool
		wantDesc   string
	}{
		{
			name:      "simple feat",
			raw:       "feat(auth): add token refresh",
			wantType:  "feat",
			wantScope: "auth",
			wantDesc:  "add token refresh",
		},
		{
			name:      "fix without scope",
			raw:       "fix: correct null check",
			wantType:  "fix",
			wantScope: "",
			wantDesc:  "correct null check",
		},
		{
			name:      "breaking change",
			raw:       "feat(api)!: remove deprecated endpoints",
			wantType:  "feat",
			wantScope: "api",
			wantBreak: true,
			wantDesc:  "remove deprecated endpoints",
		},
		{
			name:      "not conventional",
			raw:       "fixed the auth thing",
			wantType:  "",
			wantScope: "",
			wantDesc:  "fixed the auth thing",
		},
		{
			name:      "wip",
			raw:       "wip",
			wantType:  "",
			wantScope: "",
			wantDesc:  "wip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := ParseCommit(tt.raw)
			if cm.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", cm.Type, tt.wantType)
			}
			if cm.Scope != tt.wantScope {
				t.Errorf("Scope = %q, want %q", cm.Scope, tt.wantScope)
			}
			if cm.Breaking != tt.wantBreak {
				t.Errorf("Breaking = %v, want %v", cm.Breaking, tt.wantBreak)
			}
			if cm.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", cm.Description, tt.wantDesc)
			}
		})
	}
}

func TestValidateCommit(t *testing.T) {
	tests := []struct {
		name           string
		raw            string
		scopeRequired  bool
		imperative     string
		enforcement    string
		wantViolations int
		wantBlocking   bool
	}{
		{
			name:           "valid commit",
			raw:            "feat(auth): add token refresh",
			scopeRequired:  true,
			imperative:     "warn",
			enforcement:    "STANDARD",
			wantViolations: 0,
		},
		{
			name:           "missing scope",
			raw:            "fix: correct null check",
			scopeRequired:  true,
			imperative:     "warn",
			enforcement:    "STANDARD",
			wantViolations: 1,
		},
		{
			name:           "not conventional",
			raw:            "fixed the auth thing",
			scopeRequired:  true,
			imperative:     "warn",
			enforcement:    "STANDARD",
			wantViolations: 1,
			wantBlocking:   true,
		},
		{
			name:           "past tense",
			raw:            "feat(auth): added token refresh",
			scopeRequired:  true,
			imperative:     "warn",
			enforcement:    "STANDARD",
			wantViolations: 1,
		},
		{
			name:           "wip commit",
			raw:            "wip",
			scopeRequired:  true,
			imperative:     "warn",
			enforcement:    "STANDARD",
			wantViolations: 1, // not conventional (wip check is subsumed)
			wantBlocking:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := ParseCommit(tt.raw)
			violations := ValidateCommit(cm, tt.scopeRequired, tt.imperative, tt.enforcement)
			if len(violations) != tt.wantViolations {
				t.Errorf("got %d violations, want %d", len(violations), tt.wantViolations)
				for _, v := range violations {
					t.Logf("  - [%s] %s", v.Severity, v.Message)
				}
			}
			if tt.wantBlocking {
				hasBlocking := false
				for _, v := range violations {
					if v.Severity == "blocking" {
						hasBlocking = true
						break
					}
				}
				if !hasBlocking {
					t.Error("expected at least one blocking violation")
				}
			}
		})
	}
}

func TestScopeConsistency(t *testing.T) {
	// Many scopes should trigger a warning
	commits := []CommitMessage{
		{Scope: "auth"},
		{Scope: "api"},
		{Scope: "db"},
		{Scope: "ui"},
	}
	violations := ScopeConsistency(commits)
	if len(violations) != 1 {
		t.Errorf("got %d violations, want 1 for 4 scopes", len(violations))
	}

	// Few scopes should be fine
	fewCommits := []CommitMessage{
		{Scope: "auth"},
		{Scope: "auth"},
	}
	violations = ScopeConsistency(fewCommits)
	if len(violations) != 0 {
		t.Errorf("got %d violations, want 0 for 1 scope", len(violations))
	}
}
