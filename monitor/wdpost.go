package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/beck-8/filecoin-checker/notifier"

	"github.com/beck-8/filecoin-checker/config"

	"github.com/beck-8/filecoin-checker/api"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/rs/zerolog/log"
)

func CheckWindowedPoSt(ctx context.Context, client *api.LotusClient, cfg *config.MinerConfig) error {
	// pre-check
	timeout := cfg.Timeout
	slient := cfg.Slient
	sleepInterval := cfg.SleepInterval
	if timeout == 0 {
		timeout = config.Global.Global.Timeout
	}
	if slient == 0 {
		slient = config.Global.Global.Slient
	}
	if sleepInterval == 0 {
		sleepInterval = config.Global.Global.SleepInterval
	}

	if timeout > 30*60 || slient > 30*60 || sleepInterval > 30*60 {
		return fmt.Errorf("Configuration error: timeout, slient, sleep_interval cannot be greater than 30")
	}

	addr, err := address.NewFromString(cfg.MinerID)
	if err != nil {
		return err
	}
	dlineInfo, err := client.StateMinerProvingDeadline(ctx, addr, types.EmptyTSK)
	if err != nil {
		return err
	}
	startDuration := uint64(dlineInfo.CurrentEpoch - dlineInfo.Open)
	if !dlineInfo.IsOpen() {
		log.Warn().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg("WindowedPoSt has not started yet")
		return nil
	}
	if startDuration < uint64(timeout/30) || startDuration > uint64(slient/30) {
		log.Debug().Str("miner", cfg.MinerID).
			Uint64("deadline", dlineInfo.Index).
			Int("timeout", timeout).
			Int("slient", slient).
			Uint64("startDuration", startDuration*30).
			Msg("Current time does not meet check conditions")
		return nil
	}
	deads, err := client.StateMinerDeadlines(ctx, addr, types.EmptyTSK)
	if err != nil {
		return err
	}
	parttitons, err := client.StateMinerPartitions(ctx, addr, dlineInfo.Index, types.EmptyTSK)
	if err != nil {
		return err
	}

	var total int64
	// Partitions with live sectors
	// This must be done this way because some partitions may contain only expired sectors
	for _, partiton := range parttitons {
		if liveCount, err := partiton.LiveSectors.Count(); err == nil && liveCount > 0 {
			total += 1
		} else if err != nil {
			return err
		}
	}
	if total == 0 {
		log.Info().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg(fmt.Sprintf("No partitions need to be submitted, waiting %vs before checking again", (dlineInfo.Close-dlineInfo.CurrentEpoch)*30))
		time.Sleep(time.Second * time.Duration((dlineInfo.Close-dlineInfo.CurrentEpoch)*30))
		return nil
	}
	submmited, err := deads[int(dlineInfo.Index)].PostSubmissions.Count()
	if err != nil {
		return err
	}
	if submmited < uint64(total) {
		body := fmt.Sprintf("deadline %v has submitted %v Partitions out of %v required", dlineInfo.Index, submmited, total)
		title := fmt.Sprintf("%s WindowedPoSt timeout %vmin not submitted", cfg.MinerID, startDuration*30/60)
		log.Error().Str("miner", cfg.MinerID).Str("title", title).Str("body", body).Msg("WindowedPoSt timeout not submitted")
		err := notifier.SendNotify(cfg.MinerID,
			body,
			title,
			cfg.RecipientURLs, cfg.AppriseAPIServer)
		if err != nil {
			return err
		}
		log.Debug().Str("miner", cfg.MinerID).Int("sleepInterval", sleepInterval).Msg("Waiting ⌛️ before continuing to check WindowedPoSt to prevent frequent alerts")
		time.Sleep(time.Second * time.Duration(sleepInterval))
		return nil
	}

	// All submissions completed, sleep until the next deadline
	remainingTime := (dlineInfo.Close - dlineInfo.CurrentEpoch) * 30
	log.Info().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg(fmt.Sprintf("WindowedPoSt fully submitted, waiting %vs before checking again", remainingTime))
	time.Sleep(time.Second * time.Duration(remainingTime))
	return nil
}

func CheckWDPost(ctx context.Context, client *api.LotusClient, c *config.MinerConfig) {
	for {
		err := CheckWindowedPoSt(ctx, client, c)
		if err != nil {
			log.Error().Str("miner", c.MinerID).Err(err).Msg("Failed to check WdPost")
		}
		time.Sleep(time.Second * time.Duration(config.Global.Global.CheckInterval))
	}
}
