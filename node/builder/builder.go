package builder

import (
	"fmt"

	//"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/archway-network/archway/app/params"

	"github.com/nuclearblock/archgregator/node"
	nodeconfig "github.com/nuclearblock/archgregator/node/config"
	"github.com/nuclearblock/archgregator/node/local"
	"github.com/nuclearblock/archgregator/node/remote"
)

func BuildNode(cfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (node.Node, error) {
	switch cfg.Type {
	case nodeconfig.TypeRemote:
		return remote.NewNode(cfg.Details.(*remote.Details), encodingConfig.Marshaler)
	case nodeconfig.TypeLocal:
		return local.NewNode(cfg.Details.(*local.Details), encodingConfig.TxConfig, encodingConfig.Marshaler)
	case nodeconfig.TypeNone:
		return nil, nil

	default:
		return nil, fmt.Errorf("invalid node type: %s", cfg.Type)
	}
}
