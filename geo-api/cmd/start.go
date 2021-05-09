package cmd

import (
	"findhotel.net/geo-api/api"
	"github.com/spf13/cobra"
)

var startCMD = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Long:  "start command start http server to provide geolocation services",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartHttpServer()
	},
}
