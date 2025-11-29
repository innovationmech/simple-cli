package serve

import (
	"fmt"

	"github.com/innovationmech/simple-cli/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the application",
		Long:  "Serve the application",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := server.NewServer()
			return server.Run(fmt.Sprintf(":%d", viper.GetInt("port")))
		},
	}
}
