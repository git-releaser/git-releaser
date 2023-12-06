package helpers

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindViperFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := v.BindPFlag(f.Name, f); err != nil {
			fmt.Printf("Error binding flag '%s': %v\n", f.Name, err)
		}
	})
}
