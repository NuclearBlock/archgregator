package database

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/nuclearblock/archgregator/logging"

	databaseconfig "github.com/nuclearblock/archgregator/database/config"

	"github.com/nuclearblock/archgregator/types"
)

// Database represents an abstract database that can be used to save data inside it
type Database interface {
	// HasBlock tells whether or not the database has already stored the block having the given height.
	// An error is returned if the operation fails.
	HasBlock(height int64) (bool, error)

	// SaveBlock will be called when a new block is parsed, passing the block itself
	// and the transactions contained inside that block.
	// An error is returned if the operation fails.
	SaveBlock(block *types.Block) error

	// SaveTx will be called to save each transaction contained inside a block.
	// An error is returned if the operation fails.
	SaveTx(tx *types.Tx) error

	// SaveMessage stores a single message.
	// An error is returned if the operation fails.
	SaveMessage(msg *types.Message) error

	SaveWasmCode(wasmCode types.WasmCode) error
	SaveWasmContract(wasmContract types.WasmContract) error
	SaveWasmExecuteContract(executeContract types.WasmExecuteContract) error
	UpdateContractWithMsgMigrateContract(sender string, contractAddress string, codeID uint64, rawContractMsg []byte, data string) error
	UpdateContractAdmin(sender string, contractAddress string, newAdmin string) error

	SaveContractRewardCalculation(contractRewardCalculation types.ContractRewardCalculation) error
	SaveContractRewardDistribution(contractRewardDistribution types.ContractRewardDistribution) error

	// Close closes the connection to the database
	Close()
}

// Context contains the data that might be used to build a Database instance
type Context struct {
	Cfg            databaseconfig.Config
	EncodingConfig *params.EncodingConfig
	Logger         logging.Logger
}

// NewContext allows to build a new Context instance
func NewContext(cfg databaseconfig.Config, encodingConfig *params.EncodingConfig, logger logging.Logger) *Context {
	return &Context{
		Cfg:            cfg,
		EncodingConfig: encodingConfig,
		Logger:         logger,
	}
}

// Builder represents a method that allows to build any database from a given codec and configuration
type Builder func(ctx *Context) (Database, error)
