package database

import (
	//"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/archway-network/archway/app/params"

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

	// SaveWasmCode stores a single WASM Code.
	// An error is returned if the operation fails.
	SaveWasmCode(wasmCode types.WasmCode) error

	// SaveWasmContract stores an contract instance of WASM Code.
	// An error is returned if the operation fails.
	SaveWasmContract(wasmContract types.WasmContract) error

	// SaveWasmExecuteContract stores each contract execution.
	// An error is returned if the operation fails.
	SaveWasmExecuteContract(executeContract types.WasmExecuteContract) error

	// SaveSaveContractRewardCalculation helps add to db a gastracker reward data.
	// When Calculation event will be processed - we can add to db initial rewards data
	// An error is returned if the operation fails.
	SaveContractRewardCalculation(contractRewardCalculation types.ContractRewardCalculation) error

	// SaveContractRewardDistribution helps add to db row a gastracker reward data.
	// When Distribution event will be processed - we can update db with rewards data
	// An error is returned if the operation fails.
	SaveContractRewardDistribution(contractRewardDistribution types.ContractRewardDistribution) error

	// SaveGasTrackerContractMetadata stores each Gastracker Metadata set attempt.
	// An error is returned if the operation fails.
	SaveGasTrackerContractMetadata(gastrackerContractMetadata types.GasTrackerContractMetadata) error

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
