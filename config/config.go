package config

import (
	_ "embed"
)

// GlobalConfig stores global configuration
type GlobalConfig struct {
	LotusAPI         string   `yaml:"lotus_api"`
	AuthToken        string   `yaml:"auth_token"`
	CheckInterval    int      `yaml:"check_interval"`
	Timeout          int      `yaml:"timeout"`
	Slient           int      `yaml:"slient"`
	SleepInterval    int      `yaml:"sleep_interval"`
	FaultsSectors    int      `yaml:"faults_sectors"`
	AppriseAPIServer string   `yaml:"apprise_api_server"`
	RecipientURLs    []string `yaml:"recipient_urls"`
}

// MinerConfig stores miner-specific configuration
type MinerConfig struct {
	MinerID          string   `yaml:"miner_id"`
	Timeout          int      `yaml:"timeout"`
	Slient           int      `yaml:"slient"`
	SleepInterval    int      `yaml:"sleep_interval"`
	FaultsSectors    int      `yaml:"faults_sectors"`
	AppriseAPIServer string   `yaml:"apprise_api_server"`
	RecipientURLs    []string `yaml:"recipient_urls"`
}

// Config overall configuration structure
type Config struct {
	Global GlobalConfig   `yaml:"global"`
	Miners []*MinerConfig `yaml:"miners"`
}

//go:embed config.example.yaml
var DefaultConfigTemplate []byte

var Global = &Config{}
