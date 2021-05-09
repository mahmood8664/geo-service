package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

var C Config

type Config struct {
	FileAddress     string
	MongodbUrl      string
	MongodbUsername string
	MongodbPassword string
	InsertBulkSize  int
	NumberOfWorkers int
}

func Init() {
	configLogger()
	C = Config{
		FileAddress:     "",
		MongodbUrl:      "",
		MongodbUsername: "",
		MongodbPassword: "",
		InsertBulkSize:  4000,
		NumberOfWorkers: 4,
	}
}

func configLogger() {
	level, err := zerolog.ParseLevel("info")
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(os.Stdout).With().Logger()
}
