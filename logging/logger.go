package logging

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/nuclearblock/archgregator/types"
)

const (
	LogKeyHeight  = "height"
	LogKeyTxHash  = "tx_hash"
	LogKeyMsgType = "msg_type"
)

// Logger defines a function that takes an error and logs it.
type Logger interface {
	SetLogLevel(level string) error
	SetLogFormat(format string) error

	Info(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})

	GenesisError(err error)
	BlockError(block *tmctypes.ResultBlock, err error)
	EventsError(results *tmctypes.ResultBlock, err error)
	TxError(tx *types.Tx, err error)
	MsgError(tx *types.Tx, msg sdk.Msg, err error)
}
