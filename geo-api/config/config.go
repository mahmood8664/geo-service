package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var (
	envPrefix = "GEO_API"
	C         Config
)

type Config struct {
	HttpPort string  `yaml:"http_port"`
	Logging  Logging `yaml:"logging"`
	Cors     Cors    `yaml:"cors"`
	MongoDB  Mongodb `yaml:"mongodb"`
}

type Logging struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	FileName string `yaml:"file_name"`
}

type Cors struct {
	Domain string `yaml:"domain"`
}

type Mongodb struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init(filename string) {
	loadConfigs(filename)
	logConfigure()
}

//loadConfigs load configuration from config file and environment variables
func loadConfigs(filename string) {
	c := Config{}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if filename != "" {
		viper.SetConfigFile(filename)
		if err := viper.MergeInConfig(); err != nil {
			log.Fatal().Msgf("loading configs file [%s] failed: %s", filename, err)
		} else {
			log.Info().Msgf("configs file [%s] loaded successfully", filename)
		}
	} else {
		log.Fatal().Msg("Config file is not determined")
	}

	err := viper.Unmarshal(&c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	})
	if err != nil {
		log.Fatal().Msg("failed on configs unmarshal: " + err.Error())
	}

	C = c
}

func logConfigure() {
	level, err := zerolog.ParseLevel(C.Logging.Level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	//log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}) //Pretty output
	//log.Logger = log.Output(os.Stdout).With().Logger() //Json output
}
