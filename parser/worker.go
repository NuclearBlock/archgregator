package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nuclearblock/archgregator/logging"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/nuclearblock/archgregator/database"
	"github.com/nuclearblock/archgregator/types/config"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	gastrackertypes "github.com/archway-network/archway/x/gastracker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/nuclearblock/archgregator/node"
	"github.com/nuclearblock/archgregator/types"
	"github.com/nuclearblock/archgregator/types/utils"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	index  int
	queue  types.HeightQueue
	codec  codec.Codec
	node   node.Node
	db     database.Database
	logger logging.Logger
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context, queue types.HeightQueue, index int) Worker {
	return Worker{
		index:  index,
		codec:  ctx.EncodingConfig.Marshaler,
		node:   ctx.Node,
		queue:  queue,
		db:     ctx.Database,
		logger: ctx.Logger,
	}
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w Worker) Start() {
	logging.WorkerCount.Inc()

	for i := range w.queue {
		if err := w.ProcessIfNotExists(i); err != nil {
			// re-enqueue any failed job
			// TODO: Implement exponential backoff or max retries for a block height.
			go func() {
				w.logger.Error("re-enqueueing failed block", "height", i, "err", err)
				w.queue <- i
			}()
		}

		logging.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.index)).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet. It returns an
// error if any export process fails.
func (w Worker) ProcessIfNotExists(height int64) error {
	exists, err := w.db.HasBlock(height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		w.logger.Debug("skipping already exported block", "height", height)
		return nil
	}

	return w.Process(height)
}

// Process fetches  a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (w Worker) Process(height int64) error {
	// process genesis if needed
	if height == 0 {
		cfg := config.Cfg.Parser
		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}
		return w.HandleGenesis(genesisDoc, genesisState)
	}

	w.logger.Debug("processing block", "height", height)

	block, err := w.node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := w.node.BlockResults(height)
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := w.node.Txs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	return w.ExportBlock(block, events, txs)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w Worker) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// TO-DO, Not yet implemented...
	return nil
}

// ExportBlock accepts a finalized block and a corresponding set of transactions
// and persists them to the database along with attributable metadata. An error
// is returned if the write fails.
func (w Worker) ExportBlock(b *tmctypes.ResultBlock, r *tmctypes.ResultBlockResults, txs []*types.Tx) error {
	// Save block to database
	err := w.db.SaveBlock(types.NewBlockFromTmBlock(b, sumGasTxs(txs)))
	if err != nil {
		return fmt.Errorf("failed to save block: %s", err)
	}

	// Get block date for gastracker rewards usage
	timeStampBlock := b.Block.Time.UTC()

	// Process block events to fing gastracker rewards
	err = w.ProcessEvents(r, timeStampBlock)
	if err != nil {
		return fmt.Errorf("failed to process events: %s", err)
	}

	// Process block transactions to find messages, that contains contracts data
	err = w.ProcessTransactions(txs)
	if err != nil {
		return fmt.Errorf("failed to process transactions: %s", err)
	}

	return nil
}

// ProcessEvents accepts a set of events of current BeginBlock
// Events will be processed to catch gastracker rewards
func (w Worker) ProcessEvents(r *tmctypes.ResultBlockResults, ts time.Time) error {
	for _, event := range r.BeginBlockEvents {
		// Only 'gastracker' events
		if strings.Contains(event.Type, gastrackertypes.ModuleName) {
			err := HandleGasTrackerRewards(&event, r.Height, ts, w.db)
			if err != nil {
				//return fmt.Errorf("error while handle gas tracker rewards: %s", err)
				w.logger.Error("error while handle gas tracker rewards: ", err)
			}
		}
	}
	return nil
}

// ProcessTransactions accepts a set of transactions of current block
// and process them to find wasm and gastracker messages.
func (w Worker) ProcessTransactions(txs []*types.Tx) error {
	// Handle all the transactions inside the block
	for _, tx := range txs {
		// Process only successful txs
		if !tx.Successful() {
			continue
		}
		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err := w.codec.UnpackAny(msg, &stdMsg)
			if err != nil {
				//  return fmt.Errorf("error while unpacking message: %s", err)
				w.logger.Error("error while unpacking message: ", err)
				continue
			}
			// Seek wasm and gastracker messages
			switch cosmosMsg := stdMsg.(type) {
			case *wasmtypes.MsgStoreCode:
				// Wasm code store
				err = HandleMsgStoreCode(i, tx, cosmosMsg, w.node, w.db)
				if err != nil {
					w.logger.MsgError(tx, cosmosMsg, err)
				}
			case *wasmtypes.MsgInstantiateContract:
				// Wasm contract instantiate
				err = HandleMsgInstantiateContract(i, tx, cosmosMsg, w.node, w.db)
				if err != nil {
					w.logger.MsgError(tx, cosmosMsg, err)
				}
			case *wasmtypes.MsgExecuteContract:
				// Wasm contract execute
				err = HandleMsgExecuteContract(i, tx, cosmosMsg, w.db)
				if err != nil {
					w.logger.MsgError(tx, cosmosMsg, err)
				}
			case *gastrackertypes.MsgSetContractMetadata:
				// Gastracker metadata set
				err = HandleMsgSetMetadata(i, tx, cosmosMsg, w.db)
				if err != nil {
					w.logger.MsgError(tx, cosmosMsg, err)
				}
			}
		}
	}
	return nil
}
