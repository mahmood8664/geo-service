package cmd

import (
	"findhotel.net/geo-api/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "GeoApi",
	Short: "GeoApi",
}
var configFilePath string

func init() {
	cobra.OnInitialize(Configure)
	rootCMD.AddCommand(startCMD)
	rootCMD.PersistentFlags().StringVar(&configFilePath, "config", "resources/config.yml", "config file address")
}

func Configure() {
	config.Init(configFilePath)
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
