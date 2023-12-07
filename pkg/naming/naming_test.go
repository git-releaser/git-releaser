package naming

import (
	"testing"
)

func TestCreatePrDescription(t *testing.T) {
	version := "1.0.0"
	changelog := "- Fixed bugs\n- Added features"

	expected := "This is a description for the new pull request for version 1.0.0.\n\n## Changelog\n\n- Fixed bugs\n- Added features"

	result := CreatePrDescription(version, changelog)

	if result != expected {
		t.Errorf("Unexpected result: got %v, want %v", result, expected)
	}
}
