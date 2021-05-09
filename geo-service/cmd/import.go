package cmd

import (
	"findhotel.com/geo-service/csv"
	"findhotel.com/geo-service/di"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	importer csv.Importer
)

//importCMD is command to import csv into database
var importCMD = &cobra.Command{
	Use:   "import",
	Short: "start csv importer",
	Long:  "use this command to import csv data into database.",
	Run: func(cmd *cobra.Command, args []string) {
		importer = di.CreateImporter()
		err := importer.Import()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot import csv")
		}
	},
}
