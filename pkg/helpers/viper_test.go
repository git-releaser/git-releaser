package helpers

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"testing"
)

func TestBindViperFlags(t *testing.T) {
	// Create a new cobra command
	cmd := &cobra.Command{
		Use:   "test",
		Short: "A test command",
	}

	// Add a flag to the command
	cmd.Flags().String("test-flag", "default", "A test flag")

	// Create a new viper instance
	v := viper.New()

	// Call the function
	BindViperFlags(cmd, v)

	// Check that the flag was bound correctly
	if v.GetString("test-flag") != "default" {
		t.Errorf("Unexpected flag value: got %v, want %v", v.GetString("test-flag"), "default")
	}
}
