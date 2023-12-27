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

func TestGeneratePrTitle(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GeneratePrTitle",
			args: args{
				version: "0.0.1",
			},
			want: "Release 0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePrTitle(tt.args.version); got != tt.want {
				t.Errorf("GeneratePrTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
