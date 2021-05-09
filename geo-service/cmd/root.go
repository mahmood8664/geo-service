package cmd

import "C"
import (
	"findhotel.com/geo-service/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//rootCMD is root command of project, see importCMD for import command
var rootCMD = &cobra.Command{
	Use:   "csv-importer",
	Short: "csv-importer",
	Long:  "CSV importer is a go module to import csv data into MongoDB database and also provide model to access data from database",
}

func init() {
	config.Init()
	//Config importCMD
	rootCMD.AddCommand(importCMD)
	importCMD.Flags().StringVarP(&config.C.FileAddress, "file", "f", "", "CSV file address")
	importCMD.Flags().StringVarP(&config.C.MongodbUrl, "db-url", "d", "mongodb://localhost:27017", "MongoDb URL")
	importCMD.Flags().StringVarP(&config.C.MongodbUsername, "username", "u", "", "MongoDB username")
	importCMD.Flags().StringVarP(&config.C.MongodbPassword, "password", "p", "", "MongoDB password")
	err := importCMD.MarkFlagRequired("file")
	if err != nil {
		log.Fatal().Err(err)
	}
}

//Execute runs command processor
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
