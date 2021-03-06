package types

import (
	"github.com/archway-network/archway/app/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nuclearblock/archgregator/types/config"
)

// SdkConfigSetup represents a method that allows to customize the given sdk.Config.
// This should be used to set custom Bech32 addresses prefixes and other app-related configurations.
type SdkConfigSetup func(config config.Config, sdkConfig *sdk.Config)

// DefaultConfigSetup represents a handy implementation of SdkConfigSetup that simply setups the prefix
// inside the configuration
func DefaultConfigSetup(cfg config.Config, sdkConfig *sdk.Config) {
	prefix := cfg.Chain.Bech32Prefix
	sdkConfig.SetBech32PrefixForAccount(
		prefix,
		prefix+sdk.PrefixPublic,
	)
}

// -----------------------------------------------------------------

// EncodingConfigBuilder represents a function that is used to return the proper encoding config.
type EncodingConfigBuilder func() params.EncodingConfig
