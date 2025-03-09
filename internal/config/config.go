package config

import (
	"encoding/json"
	"os"
)

// Config 系统配置结构体
type Config struct {
	// 缓存配置
	Cache CacheConfig `json:"cache"`
	// HTTP服务器配置
	HTTP HTTPConfig `json:"http"`
	// 一致性哈希配置
	Hash HashConfig `json:"hash"`
}

// CacheConfig 缓存相关配置
type CacheConfig struct {
	// 最大缓存条目数
	MaxEntries int `json:"maxEntries"`
	// 最大缓存容量（字节）
	MaxBytes int64 `json:"maxBytes"`
}

// HTTPConfig HTTP服务器配置
type HTTPConfig struct {
	// 服务器地址
	Address string `json:"address"`
	// 服务器端口
	Port int `json:"port"`
}

// HashConfig 一致性哈希配置
type HashConfig struct {
	// 虚拟节点数
	Replicas int `json:"replicas"`
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Cache: CacheConfig{
		MaxEntries: 1000,
		MaxBytes:   1 << 20, // 默认1MB
	},
	HTTP: HTTPConfig{
		Address: "0.0.0.0",
		Port:    8001,
	},
	Hash: HashConfig{
		Replicas: 50,
	},
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	config := DefaultConfig

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析JSON配置
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
