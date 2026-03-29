package engine

import "testing"

func TestGetDME(t *testing.T) {
	// All documented DME entries should be present
	expectedEntries := []string{
		"DME-0001", "DME-0042", "DME-0047", "DME-0088", "DME-0100",
		"DME-0143", "DME-0404", "DME-0408", "DME-0410", "DME-0418",
		"DME-0500", "DME-0508", "DME-0911", "DME-1337", "DME-1997",
		"DME-2001", "DME-3000",
	}

	for _, id := range expectedEntries {
		t.Run(id, func(t *testing.T) {
			entry, ok := GetDME(id)
			if !ok {
				t.Fatalf("DME entry %s not found", id)
			}
			if entry.ID != id {
				t.Errorf("entry.ID = %s, want %s", entry.ID, id)
			}
			if entry.Status != "Rejected" {
				t.Errorf("entry.Status = %s, want Rejected", entry.Status)
			}
			if entry.Title == "" {
				t.Error("entry.Title is empty")
			}
			if entry.RejectionRationale == "" {
				t.Error("entry.RejectionRationale is empty")
			}
		})
	}
}

func TestGetDMENotFound(t *testing.T) {
	_, ok := GetDME("DME-9999")
	if ok {
		t.Error("Expected DME-9999 to not be found")
	}
}

func TestRickrollEntries(t *testing.T) {
	// The commitment protocol series should map to Rick Astley lyrics
	rickrollEntries := map[string]string{
		"DME-0143": "Never gonna give you up",
		"DME-0404": "Never gonna let you down",
		"DME-0408": "Never gonna run around",
		"DME-0410": "and desert you",
		"DME-0418": "Never gonna make you cry",
		"DME-0500": "Never gonna say goodbye",
		"DME-0508": "Never gonna tell a lie",
		"DME-0911": "and hurt you",
	}

	for id, expectedMapsTo := range rickrollEntries {
		t.Run(id, func(t *testing.T) {
			entry, ok := GetDME(id)
			if !ok {
				t.Fatalf("DME entry %s not found", id)
			}
			if entry.MapsTo != expectedMapsTo {
				t.Errorf("entry.MapsTo = %q, want %q", entry.MapsTo, expectedMapsTo)
			}
		})
	}
}
