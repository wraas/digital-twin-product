package engine

import "testing"

func TestCalibrateSigh(t *testing.T) {
	tests := []struct {
		input    string
		expected SighLevel
	}{
		{"wip", SighExistential},
		{"", SighExistential},
		{"fix stuff", SighDeep},
		{"update the thing", SighModerate},
		{"Should I use Redis or a second database as a caching layer?", SighMild},
		{"Given the context of our microservices architecture, should we introduce a message queue between the auth service and the notification service?", SighSilent},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := CalibrateSigh(tt.input)
			if got != tt.expected {
				t.Errorf("CalibrateSigh(%q) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSighIntensity(t *testing.T) {
	if SighIntensity(SighSilent) >= SighIntensity(SighMild) {
		t.Error("SILENT should be less intense than MILD")
	}
	if SighIntensity(SighMild) >= SighIntensity(SighExistential) {
		t.Error("MILD should be less intense than EXISTENTIAL")
	}
}

func TestShouldEmit(t *testing.T) {
	if !ShouldEmit(SighExistential, SighMild) {
		t.Error("EXISTENTIAL should emit when threshold is MILD")
	}
	if ShouldEmit(SighSilent, SighMild) {
		t.Error("SILENT should not emit when threshold is MILD")
	}
	if !ShouldEmit(SighMild, SighMild) {
		t.Error("MILD should emit when threshold is MILD")
	}
}

func TestParseSighLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected SighLevel
	}{
		{"MILD", SighMild},
		{"mild", SighMild},
		{"EXISTENTIAL", SighExistential},
		{"none", SighSilent},
		{"SILENT", SighSilent},
		{"garbage", SighMild},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseSighLevel(tt.input)
			if got != tt.expected {
				t.Errorf("ParseSighLevel(%q) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}
