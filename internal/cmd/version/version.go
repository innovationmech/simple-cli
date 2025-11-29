package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

func GetVersion() string {
	return version
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersion())
		},
	}
}
