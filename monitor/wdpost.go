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
		return fmt.Errorf("配置错误: timeout、slient、sleep_interval不能大于30")
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
		log.Warn().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg("当前未开始WindowedPoSt")
		return nil
	}
	if startDuration < uint64(timeout/30) || startDuration > uint64(slient/30) {
		log.Info().Str("miner", cfg.MinerID).
			Uint64("deadline", dlineInfo.Index).
			Int("timeout", timeout).
			Int("slient", slient).
			Uint64("startDuration", startDuration*30).
			Msg("当前未满足检查时间")
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
	// 有live扇区的分区
	// 这里必须这样做，因为有的part可能全是过期扇区
	for _, partiton := range parttitons {
		if liveCount, err := partiton.LiveSectors.Count(); err == nil && liveCount > 0 {
			total += 1
		} else if err != nil {
			return err
		}
	}
	if total == 0 {
		log.Info().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg("当前没有需要提交的分区")
		return nil
	}
	submmited, err := deads[int(dlineInfo.Index)].PostSubmissions.Count()
	if err != nil {
		return err
	}
	if submmited < uint64(total) {
		body := fmt.Sprintf("deadline %v 当前已提交%v个Partitions，总共需要提交%v个", dlineInfo.Index, submmited, total)
		title := fmt.Sprintf("%s WindowedPoSt超时%vmin未提交", cfg.MinerID, startDuration*30/60)
		log.Error().Str("miner", cfg.MinerID).Str("title", title).Str("body", body).Msg("WindowedPoSt超时未提交")
		err := notifier.SendNotify(cfg.MinerID,
			body,
			title,
			cfg.RecipientURLs, cfg.AppriseAPIServer)
		if err != nil {
			return err
		}
		log.Debug().Str("miner", cfg.MinerID).Int("sleepInterval", sleepInterval).Msg("等待⌛️一会继续检查WindowedPoSt,防止重复频繁发送告警")
		time.Sleep(time.Second * time.Duration(sleepInterval))
		return nil
	}

	// 都提交完成了，睡到下一个deadline
	remainingTime := (dlineInfo.Close - dlineInfo.CurrentEpoch) * 30
	log.Info().Str("miner", cfg.MinerID).Uint64("deadline", dlineInfo.Index).Msg(fmt.Sprintf("WindowedPoSt已全部提交,等到 %vs 后继续检查", remainingTime))
	time.Sleep(time.Second * time.Duration(remainingTime))
	return nil
}

func CheckWDPost(ctx context.Context, client *api.LotusClient, c *config.MinerConfig) {
	for {
		err := CheckWindowedPoSt(ctx, client, c)
		if err != nil {
			log.Error().Str("miner", c.MinerID).Err(err).Msg("检查WdPost失败")
		}
		time.Sleep(time.Second * time.Duration(config.Global.Global.CheckInterval))
	}
}
