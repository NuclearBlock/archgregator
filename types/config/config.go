package config

import (
	databaseconfig "github.com/nuclearblock/archgregator/database/config"
	loggingconfig "github.com/nuclearblock/archgregator/logging/config"
	nodeconfig "github.com/nuclearblock/archgregator/node/config"
	parserconfig "github.com/nuclearblock/archgregator/parser/config"
)

var (
	// Cfg represents the configuration to be used during the execution
	Cfg Config
)

// Config defines all necessary archgregator configuration parameters.
type Config struct {
	bytes []byte

	Chain    ChainConfig           `yaml:"chain"`
	Node     nodeconfig.Config     `yaml:"node"`
	Parser   parserconfig.Config   `yaml:"parsing"`
	Database databaseconfig.Config `yaml:"database"`
	Logging  loggingconfig.Config  `yaml:"logging"`
}

// NewConfig builds a new Config instance
func NewConfig(
	nodeCfg nodeconfig.Config,
	chainCfg ChainConfig,
	dbConfig databaseconfig.Config,
	parserConfig parserconfig.Config,
	loggingConfig loggingconfig.Config,
) Config {
	return Config{
		Node:     nodeCfg,
		Chain:    chainCfg,
		Database: dbConfig,
		Parser:   parserConfig,
		Logging:  loggingConfig,
	}
}

func DefaultConfig() Config {
	return NewConfig(
		nodeconfig.DefaultConfig(),
		DefaultChainConfig(),
		databaseconfig.DefaultDatabaseConfig(),
		parserconfig.DefaultParsingConfig(),
		loggingconfig.DefaultLoggingConfig(),
	)
}

func (c Config) GetBytes() ([]byte, error) {
	return c.bytes, nil
}

type ChainConfig struct {
	Bech32Prefix string `yaml:"bech32_prefix"`
}

// NewChainConfig returns a new ChainConfig instance
func NewChainConfig(bech32Prefix string) ChainConfig {
	return ChainConfig{
		Bech32Prefix: bech32Prefix,
	}
}

// DefaultChainConfig returns the default instance of ChainConfig
func DefaultChainConfig() ChainConfig {
	return NewChainConfig("archway")
}
