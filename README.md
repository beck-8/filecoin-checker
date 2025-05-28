# Filecoin Checker

[English](README.md) | [中文](README.zh.md)

## Project Introduction

Filecoin Checker is a tool for monitoring the WindowedPoSt status of Filecoin miners. It can check miners' WindowedPoSt submission status and faulty sector count in real-time, sending notifications through various channels when problems occur, helping miners discover and resolve issues promptly.

## Features

- **WindowedPoSt Monitoring**: Detects whether miners submit WindowedPoSt proofs on time, avoiding penalties due to late submissions
- **Faulty Sector Monitoring**: Monitors the number of faulty sectors for miners, triggering alerts when exceeding thresholds
- **Multi-Miner Support**: Can monitor multiple miner IDs simultaneously
- **Flexible Configuration**: Supports global configuration and miner-level customization
- **Multi-Channel Notifications**: Based on the Apprise notification system, supporting 100+ notification channels (such as Telegram, Discord, etc.)

## Installation

### Method 1: Using Docker

```bash
# Pull the image
docker pull ghcr.io/beck-8/filecoin-checker:latest

# Run the container
docker run -d --name filecoin-checker \
  ghcr.io/beck-8/filecoin-checker:latest
```

### Method 2: Building from Source

```bash
# Clone the repository
git clone https://github.com/beck-8/filecoin-checker.git
cd filecoin-checker

# Build
make build

# Run
./filecoin-checker
```

## Configuration

Before running, you need to create a `config.yaml` configuration file. You can copy `config/config.example.yaml` and modify it:

```bash
cp config/config.example.yaml config.yaml
```

Or running the program directly will generate a default configuration

### Configuration Parameters

```yaml
global:
    # Lotus RPC address, supports http, ws and other protocols
    lotus_api: "http://your-lotus-node:1234/rpc/v1"
    # Lotus RPC authentication token
    auth_token: ""
    # Check interval in seconds
    check_interval: 30

    # The following configurations can be overridden in miners configuration, allowing customization for each miner
    # If WindowedPoSt is not detected 10 minutes after deadline starts, consider it problematic
    timeout: 600
    # After 25 minutes from deadline start, stop checking because it's too late
    slient: 1500
    # After a WindowedPoSt issue, sleep for a while to prevent frequent notifications
    sleep_interval: 60
    # Only alert when faulty sectors exceed 100
    faults_sectors: 100
    # apprise_api_server address
    apprise_api_server: "https://your-apprise-server/notify"
    # Notification channels, supports 100+ types
    # For detailed usage, check the apprise documentation
    recipient_urls:
        - "telegram://bot_token:api_key/chat_id"
        # - "discord://webhook_id/webhook_token"

miners:
  - miner_id: f01234567
    # The following parameters are optional, global configuration will be used if not set
    # timeout: 600
    # slient: 1500
    # sleep_interval: 120
    # faults_sectors: 100
    # apprise_api_server: "http://localhost:8000/notify"
    # recipient_urls:
    #     - "telegram://bot_token:api_key/chat_id"
  - miner_id: f07654321
```

### Notification Configuration

This project uses [Apprise](https://github.com/caronc/apprise) as the notification system, supporting 100+ notification channels. You need to:
1. Set up your own Apprise API server ([Vercel deployment](https://github.com/beck-8/subs-check?tab=readme-ov-file#vercel-serverless-%E9%83%A8%E7%BD%B2), [docker deployment](https://github.com/beck-8/subs-check?tab=readme-ov-file#docker%E9%83%A8%E7%BD%B2), etc.)
2. Set `apprise_api_server` to your Apprise API server address
3. Configure notification target URLs in `recipient_urls`

Common notification channel examples:

- Telegram: `telegram://bot_token:api_key/chat_id`
- Discord: `discord://webhook_id/webhook_token`
- Email: `mailto://user:password@gmail.com`

For more notification channel configurations, please refer to the [Apprise documentation](https://github.com/caronc/apprise/wiki)

## Usage

1. Configure the `config.yaml` file
2. Start the program: `./filecoin-checker` or run with Docker
3. The program will automatically start monitoring the configured miner IDs
4. When it detects that WindowedPoSt is not submitted on time or the number of faulty sectors exceeds the threshold, it will send alerts through the configured notification channels

## Log Description

The program outputs logs during operation, including the following information:

- Startup information: version number, number of configured miners
- Monitoring information: WindowedPoSt status and faulty sector count for each miner
- Alert information: detailed alert content when problems are detected
- Notification status: information about successful or failed notification delivery

## License

[MIT License](LICENSE)