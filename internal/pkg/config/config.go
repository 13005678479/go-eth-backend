package config

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
	"os"
)

// Config 结构体定义配置文件结构
type Config struct {
	Ethereum EthereumConfig `yaml:"ethereum"`
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type EthereumConfig struct {
	Accounts AccountsConfig `yaml:"accounts"`
	Networks NetworksConfig `yaml:"networks"`
}

type AccountsConfig struct {
	TestPrivateKey string `yaml:"test_private_key"`
}

type NetworksConfig struct {
	Mainnet NetworkConfig `yaml:"mainnet"`
	Sepolia NetworkConfig `yaml:"sepolia"`
}

type NetworkConfig struct {
	RPCURL  string `yaml:"rpc_url"`
	ChainID int64  `yaml:"chain_id"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	FilePath string `yaml:"file_path"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// GetTestPrivateKey 获取测试网络私钥
func (c *Config) GetTestPrivateKey() string {
	return c.Ethereum.Accounts.TestPrivateKey
}

// GetSepoliaConfig 获取Sepolia网络配置
func (c *Config) GetSepoliaConfig() NetworkConfig {
	return c.Ethereum.Networks.Sepolia
}

// LoadConfigOrExit 加载配置，如果失败则退出程序
func LoadConfigOrExit(configPath string) *Config {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}
	
	// 验证测试私钥是否存在
	if config.GetTestPrivateKey() == "" {
		log.Fatal("❌ 配置文件中未找到测试私钥")
	}
	
	return config
}