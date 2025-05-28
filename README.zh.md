# Filecoin Checker

[English](README.md) | [中文](README.zh.md)

## 项目介绍

Filecoin Checker 是一个用于监控 Filecoin 矿工 WindowedPoSt 状态的工具。它可以实时检查矿工的 WindowedPoSt 提交情况和故障扇区数量，在出现问题时通过多种渠道发送通知，帮助矿工及时发现并解决问题。

## 功能特点

- **WindowedPoSt 监控**：检测矿工是否按时提交 WindowedPoSt 证明，避免因未及时提交而导致的惩罚
- **故障扇区监控**：监控矿工的故障扇区数量，当超过阈值时发出警报
- **多矿工支持**：可同时监控多个矿工 ID
- **灵活配置**：支持全局配置和矿工级别的个性化配置
- **多渠道通知**：基于 Apprise 通知系统，支持 100+ 种通知渠道（如 Telegram、Discord 等）

## 安装步骤

### 方法一：使用 Docker

```bash
# 拉取镜像
docker pull ghcr.io/beck-8/filecoin-checker:latest

# 运行容器
docker run -d --name filecoin-checker \
  ghcr.io/beck-8/filecoin-checker:latest
```

### 方法二：从源码编译

```bash
# 克隆仓库
git clone https://github.com/beck-8/filecoin-checker.git
cd filecoin-checker

# 编译
make build

# 运行
./filecoin-checker
```

## 配置说明

在运行前，需要创建 `config.yaml` 配置文件。可以复制 `config/config.example.yaml` 并进行修改：

```bash
cp config/config.example.yaml config.yaml
```

或者直接运行会生成一份默认配置

### 配置参数详解

```yaml
global:
    # Lotus RPC 地址，支持 http、ws 等协议
    lotus_api: "http://your-lotus-node:1234/rpc/v1"
    # Lotus RPC 鉴权 token
    auth_token: ""
    # 检查间隔，单位秒
    check_interval: 30

    # 以下配置可在 miners 配置中覆盖，可自定义每一个 miner 的参数
    # deadline 开始后 10 分钟，还没有检测到，认为 wdpost 有问题
    timeout: 600
    # deadline 开始后 25 分钟，不进行检测，因为来不及了
    slient: 1500
    # 在 wdpost 有问题后，sleep 一会，防止一直频繁发通知
    sleep_interval: 60
    # 超过 100 个 faults 扇区，才会告警
    faults_sectors: 100
    # apprise_api_server 地址
    apprise_api_server: "https://your-apprise-server/notify"
    # 通知媒介，支持 100 种+
    # 具体使用方法查看 apprise 文档
    recipient_urls:
        - "telegram://bot_token:api_key/chat_id"
        # - "discord://webhook_id/webhook_token"

miners:
  - miner_id: f01234567
    # 以下参数可选，不设置则使用全局配置
    # timeout: 600
    # slient: 1500
    # sleep_interval: 120
    # faults_sectors: 100
    # apprise_api_server: "http://localhost:8000/notify"
    # recipient_urls:
    #     - "telegram://bot_token:api_key/chat_id"
  - miner_id: f07654321
```

### 通知配置

本项目使用 [Apprise](https://github.com/caronc/apprise) 作为通知系统，支持 100+ 种通知渠道。您需要：
1. 自建 Apprise API 服务器（[Vercel部署](https://github.com/beck-8/subs-check?tab=readme-ov-file#vercel-serverless-%E9%83%A8%E7%BD%B2)、[docker部署](https://github.com/beck-8/subs-check?tab=readme-ov-file#docker%E9%83%A8%E7%BD%B2) 等方法）
2. 设置 `apprise_api_server` 为 Apprise API 服务器地址
3. 在 `recipient_urls` 中配置通知目标 URL

常用通知渠道示例：

- Telegram: `telegram://bot_token:api_key/chat_id`
- Discord: `discord://webhook_id/webhook_token`
- Email: `mailto://user:password@gmail.com`

更多通知渠道配置请参考 [Apprise 文档](https://github.com/caronc/apprise/wiki)

## 使用方法

1. 配置好 `config.yaml` 文件
2. 启动程序：`./filecoin-checker` 或使用 Docker 运行
3. 程序会自动开始监控配置的矿工 ID
4. 当检测到 WindowedPoSt 未及时提交或故障扇区数量超过阈值时，会通过配置的通知渠道发送警报

## 日志说明

程序运行时会输出日志，包含以下信息：

- 启动信息：版本号、配置的矿工数量
- 监控信息：各矿工的 WindowedPoSt 状态、故障扇区数量
- 警报信息：当检测到问题时的详细警报内容
- 通知状态：通知发送成功或失败的信息

## 许可证

[MIT License](LICENSE)
