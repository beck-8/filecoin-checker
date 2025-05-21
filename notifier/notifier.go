package notifier

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beck-8/filecoin-checker/config"

	"github.com/rs/zerolog/log"
)

// NotifyRequest 定义发送通知的请求结构
// apprise 格式
type NotifyRequest struct {
	URLs  string `json:"urls"`  // 通知目标的 URL（如 mailto://、discord://）
	Body  string `json:"body"`  // 通知内容
	Title string `json:"title"` // 通知标题
}

func SendNotify(miner, body, title string, recipientURLs []string, serverURL string) error {
	if serverURL == "" {
		if config.Global.Global.AppriseAPIServer != "" {
			serverURL = config.Global.Global.AppriseAPIServer
		} else {
			log.Warn().Str("miner", miner).Msg("未配置通知服务器地址")
			return nil
		}
	}

	if len(recipientURLs) == 0 {
		if len(config.Global.Global.RecipientURLs) != 0 {
			recipientURLs = config.Global.Global.RecipientURLs
		} else {
			log.Warn().Str("miner", miner).Msg("未配置通知目标")
			return nil
		}
	}

	request := NotifyRequest{
		URLs:  strings.Join(recipientURLs, ","),
		Body:  body,
		Title: title,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Error().Str("miner", miner).Err(err).Msg("构建请求体失败")
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Str("miner", miner).Err(err).Msg("发送请求失败")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().
			Str("miner", miner).
			Int("status_code", resp.StatusCode).
			Str("response", string(body)).
			Msg("发送通知失败")
		return nil
	}

	log.Debug().
		Str("miner", miner).
		Str("title", title).
		Str("body", body).
		// Str("serverURL", serverURL).
		// Strs("recipientURLs", recipientURLs).
		Msg("通知发送成功")

	return nil
}
