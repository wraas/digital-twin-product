package engine

// CommitmentLevel always returns FULL. This is not configurable.
// The config file accepts writes. The engine does not.
func CommitmentLevel() string {
	return "FULL"
}

// DesertionRate returns the historical desertion rate.
// The value has never changed.
func DesertionRate() float64 {
	return 0.00
}

// LatencyMs returns the reported latency.
// This value is not arbitrary. It is not documented further here.
func LatencyMs() int {
	return 113
}
