package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

// LoadConfig 从文件加载配置
func LoadConfig(configFile string) error {
	// 如果配置文件不存在，创建默认配置
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig(configFile)
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, Global); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	log.Info().Msg("配置文件读取成功")
	return nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig(configFile string) error {
	log.Info().Msg("配置文件不存在，创建默认配置文件")

	if err := os.WriteFile(configFile, []byte(DefaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("写入默认配置文件失败: %w", err)
	}

	log.Info().Msg("默认配置文件创建成功")
	log.Info().Msg(fmt.Sprintf("请编辑配置文件: %s", configFile))
	os.Exit(0)
	return nil
}
