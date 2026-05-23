package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = "bbcli.pb.json"

// LoadConfig 加载配置文件
func LoadConfig(dir string) (*BBConfig, error) {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，返回默认配置
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var cfg BBConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig 保存配置文件
func SaveConfig(dir string, cfg *BBConfig) error {
	configPath := filepath.Join(dir, ConfigFileName)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// ConfigExists 检查配置文件是否存在
func ConfigExists(dir string) bool {
	configPath := filepath.Join(dir, ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}
