package main

import (
	"flag"
	"os"
	"time"

	"github.com/beck-8/filecoin-checker/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version = "unknown"
var CurrentCommit = "unknown"

var (
	configFile string
	logType    string
	logLevel   string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to config file")
	flag.StringVar(&logType, "log-type", "console", "Log type: console or json")
	flag.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, error")
	flag.Parse()
	// Configure logger
	if logType == "json" {
		zerolog.TimeFieldFormat = time.RFC3339
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	// Set log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level")
	}
	zerolog.SetGlobalLevel(level)

	err = config.LoadConfig(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("加载配置失败")
	}
}
