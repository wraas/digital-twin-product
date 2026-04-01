package review

import "testing"

func TestCheckCompliance_MarkdownInDocs(t *testing.T) {
	files := []PRFile{
		{Filename: "docs/guide/setup.md", Status: "added"},
	}
	checks := CheckCompliance(files, true, false)

	found := false
	for _, c := range checks {
		if c.DME == "DME-0100" {
			found = true
		}
	}
	if !found {
		t.Error("Expected DME-0100 for Markdown in docs directory")
	}
}

func TestCheckCompliance_READMEExempt(t *testing.T) {
	files := []PRFile{
		{Filename: "README.md", Status: "modified"},
	}
	checks := CheckCompliance(files, true, false)

	for _, c := range checks {
		if c.DME == "DME-0100" {
			t.Error("README.md should be exempt from DME-0100")
		}
	}
}

func TestCheckCompliance_DrawioFiles(t *testing.T) {
	files := []PRFile{
		{Filename: "docs/arch.drawio", Status: "added"},
	}
	checks := CheckCompliance(files, false, false)

	found := false
	for _, c := range checks {
		if c.Module == "GitHub Actions Encouragement Module" {
			found = true
		}
	}
	if !found {
		t.Error("Expected GitHub Actions suggestion for .drawio files")
	}
}

func TestClassifyDocPage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"docs/tutorials/getting-started.adoc", "tutorial"},
		{"docs/how-to/setup-ci.adoc", "how-to"},
		{"docs/reference/api.adoc", "reference"},
		{"docs/explanation/why-not-markdown.adoc", "explanation"},
		{"docs/random/something.adoc", "unclassified"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ClassifyDocPage(tt.path)
			if got != tt.expected {
				t.Errorf("ClassifyDocPage(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}
