package engine

import "testing"

func TestCommitmentLevel(t *testing.T) {
	if CommitmentLevel() != "FULL" {
		t.Errorf("CommitmentLevel() = %s, want FULL", CommitmentLevel())
	}
}

func TestDesertionRate(t *testing.T) {
	if DesertionRate() != 0.00 {
		t.Errorf("DesertionRate() = %f, want 0.00", DesertionRate())
	}
}

func TestLatencyMs(t *testing.T) {
	if LatencyMs() != 113 {
		t.Errorf("LatencyMs() = %d, want 113", LatencyMs())
	}
}
