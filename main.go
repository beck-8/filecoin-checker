package main

import (
	"context"
	"fmt"

	"github.com/beck-8/filecoin-checker/api"

	"github.com/beck-8/filecoin-checker/monitor"

	"github.com/rs/zerolog/log"

	"github.com/beck-8/filecoin-checker/config"
)

func main() {
	ctx := context.TODO()

	log.Info().Msg(fmt.Sprintf("Starting Filecoin Checker, Version: %s", fmt.Sprintf("%s-%s", Version, CurrentCommit)))
	log.Info().Msg(fmt.Sprintf("Total configured miners: %v", len(config.Global.Miners)))

	client, err := api.NewLotusClient(ctx, config.Global.Global.LotusAPI, config.Global.Global.AuthToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Lotus node")
	}
	defer client.Close()

	for _, miner := range config.Global.Miners {
		go monitor.CheckWDPost(ctx, client, miner)
		go monitor.CheckFault(ctx, client, miner)
	}
	select {}
}
