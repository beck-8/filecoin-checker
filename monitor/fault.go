package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/beck-8/filecoin-checker/notifier"

	"github.com/beck-8/filecoin-checker/api"
	"github.com/beck-8/filecoin-checker/config"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/rs/zerolog/log"
)

func CheckFaultSectors(ctx context.Context, client *api.LotusClient, cfg *config.MinerConfig) (err error) {
	// pre-check
	faultsSectors := cfg.FaultsSectors
	if faultsSectors == 0 {
		faultsSectors = config.Global.Global.FaultsSectors
	}
	addr, err := address.NewFromString(cfg.MinerID)
	if err != nil {
		return err
	}
	dlineInfo, err := client.StateMinerProvingDeadline(ctx, addr, types.EmptyTSK)
	if err != nil {
		return err
	}
	if !dlineInfo.IsOpen() {
		log.Warn().Str("miner", cfg.MinerID).Msg("当前未开始WindowedPoSt,跳过Fault检查")
		return nil
	}

	lastOpen := dlineInfo.Open - 60
	lastClose := dlineInfo.Close - 60

	lastOpenTSK, err := client.ChainGetTipSetByHeight(ctx, lastOpen, types.EmptyTSK)
	if err != nil {
		return err
	}
	lastCloseTSK, err := client.ChainGetTipSetByHeight(ctx, lastClose, types.EmptyTSK)
	if err != nil {
		return err
	}

	lastOpenFault, err := client.StateMinerFaults(ctx, addr, lastOpenTSK.Key())
	if err != nil {
		return err
	}
	lastCloseFault, err := client.StateMinerFaults(ctx, addr, lastCloseTSK.Key())
	if err != nil {
		return err
	}
	diff, err := bitfield.SubtractBitField(lastCloseFault, lastOpenFault)
	if err != nil {
		return err
	}
	if count, err := diff.Count(); err != nil {
		return err
	} else if count > 0 {
		log.Error().Str("miner", cfg.MinerID).Uint64("count", count).Msg("检测到新的Faults扇区")

		if count > uint64(faultsSectors) {
			err := notifier.SendNotify(cfg.MinerID,
				fmt.Sprintf("新增%v个faults扇区", count),
				fmt.Sprintf("%v 检测到新的Faults扇区", cfg.MinerID),
				cfg.RecipientURLs, cfg.AppriseAPIServer)

			if err != nil {
				return err
			}
		}
	}

	// fault检测逻辑每半小时只需要一次
	remainingTime := (dlineInfo.Close - dlineInfo.CurrentEpoch) * 30
	log.Info().Str("miner", cfg.MinerID).Msg(fmt.Sprintf("没有fault扇区,等待 %vs 后继续检查", remainingTime))
	time.Sleep(time.Second * time.Duration(remainingTime))
	return nil
}

func CheckFault(ctx context.Context, client *api.LotusClient, c *config.MinerConfig) {
	for {
		err := CheckFaultSectors(ctx, client, c)
		if err != nil {
			log.Error().Str("miner", c.MinerID).Err(err).Msg("检查Faults扇区失败")
		}
		// 里边sleep了
		// time.Sleep(time.Second * time.Duration(config.Global.Global.CheckInterval))
	}
}
