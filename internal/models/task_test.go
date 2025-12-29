package models

import (
	"encoding/json"
	"testing"
)

func TestTaskJSONMarshaling(t *testing.T) {
	// 1. Create a dummy task
	original := Task{
		ID:      999,
		Type:    "UNIT_TEST",
		Payload: "test_payload",
	}

	// 2. Convert to JSON (simulating what the API does)
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	// 3. Convert back from JSON (simulating what the Worker does)
	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	// 4. Verification
	if decoded.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, decoded.ID)
	}
	if decoded.Type != original.Type {
		t.Errorf("Expected Type %s, got %s", original.Type, decoded.Type)
	}
}
