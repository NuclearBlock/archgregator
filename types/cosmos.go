package types

import (
	"fmt"
	"time"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

// Block contains the data of a single chain block
type Block struct {
	Height    int64
	Hash      string
	TxNum     int
	TotalGas  uint64
	Timestamp time.Time
}

// NewBlock allows to build a new Block instance
func NewBlock(
	height int64, hash string, txNum int, totalGas uint64, timestamp time.Time,
) *Block {
	return &Block{
		Height:    height,
		Hash:      hash,
		TxNum:     txNum,
		TotalGas:  totalGas,
		Timestamp: timestamp,
	}
}

// NewBlockFromTmBlock builds a new Block instance from a given ResultBlock object
func NewBlockFromTmBlock(blk *tmctypes.ResultBlock, totalGas uint64) *Block {
	return NewBlock(
		blk.Block.Height,
		blk.Block.Hash().String(),
		len(blk.Block.Txs),
		totalGas,
		blk.Block.Time,
	)
}

// Tx represents an already existing blockchain transaction
type Tx struct {
	*tx.Tx
	*sdk.TxResponse
}

// NewTx allows to create a new Tx instance from the given txResponse
func NewTx(txResponse *sdk.TxResponse, tx *tx.Tx) (*Tx, error) {
	return &Tx{
		Tx:         tx,
		TxResponse: txResponse,
	}, nil
}

// FindEventByType searches inside the given tx events for the message having the specified index, in order
// to find the event having the given type, and returns it.
// If no such event is found, returns an error instead.
func (tx Tx) FindEventByType(index int, eventType string) (sdk.StringEvent, error) {
	if len(tx.Logs) > index {
		for _, ev := range tx.Logs[index].Events {
			if ev.Type == eventType {
				return ev, nil
			}
		}
	}
	return sdk.StringEvent{}, fmt.Errorf("no %s event found inside tx with hash %s", eventType, tx.TxHash)
}

// FindAttributeByKey searches inside the specified event of the given tx to find the attribute having the given key.
// If the specified event does not contain a such attribute, returns an error instead.
func (tx Tx) FindAttributeByKey(event sdk.StringEvent, attrKey string) (string, error) {
	for _, attr := range event.Attributes {
		if attr.Key == attrKey {
			return attr.Value, nil
		}
	}

	return "", fmt.Errorf("no event with attribute %s found inside tx with hash %s", attrKey, tx.TxHash)
}

// Successful tells whether this tx is successful or not
func (tx Tx) Successful() bool {
	return tx.TxResponse.Code == 0
}
