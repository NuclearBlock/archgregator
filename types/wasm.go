package types

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmCode represents the CosmWasm code in x/wasm module
type WasmCode struct {
	Creator  string
	CodeHash string
	CodeID   uint64
	Size     int
	TxHash   string
	Height   int64
}

// NewWasmCode allows to build a new x/wasm code instance from wasmtypes.MsgStoreCode
func NewWasmCode(codeInfo *wasmtypes.QueryCodeResponse, txHash string, txHeight int64) WasmCode {
	return WasmCode{
		Creator:  codeInfo.Creator,
		CodeHash: codeInfo.DataHash.String(),
		CodeID:   codeInfo.CodeID,
		Size:     codeInfo.Size(),
		TxHash:   txHash,
		Height:   txHeight,
	}
}

// WasmContract represents the CosmWasm contract in x/wasm module
type WasmContract struct {
	Sender          string
	Creator         string
	Admin           string
	CodeID          uint64
	Label           string
	RawContractMsg  []byte
	Funds           sdk.Coins
	ContractAddress string
	TxHash          string
	InstantiatedAt  time.Time
	Height          int64
}

// NewWasmCode allows to build a new x/wasm contract instance from wasmtypes.MsgStoreCode
func NewWasmContract(
	msg *wasmtypes.MsgInstantiateContract,
	contractAddress string,
	txHash string,
	instantiatedAt time.Time,
	creator string,
	height int64,
) WasmContract {
	rawContractMsg, _ := msg.Msg.MarshalJSON()

	return WasmContract{
		Sender:          msg.Sender,
		Creator:         creator,
		Admin:           msg.Admin,
		CodeID:          msg.CodeID,
		Label:           msg.Label,
		RawContractMsg:  rawContractMsg,
		Funds:           msg.Funds,
		ContractAddress: contractAddress,
		TxHash:          txHash,
		InstantiatedAt:  instantiatedAt,
		Height:          height,
	}
}

// WasmExecuteContract represents the CosmWasm execute contract in x/wasm module
type WasmExecuteContract struct {
	Sender          string
	ContractAddress string
	RawContractMsg  []byte
	Funds           sdk.Coins
	GasUsed         int64
	Fees            sdk.Coins
	TxHash          string
	ExecutedAt      time.Time
	Height          int64
}

// NewWasmExecuteContract allows to build a new x/wasm execute contract instance
// from wasmtypes.MsgExecuteContract
func NewWasmExecuteContract(
	msg *wasmtypes.MsgExecuteContract,
	tx *Tx,
	executedAt time.Time,
) WasmExecuteContract {
	rawContractMsg, _ := msg.Msg.MarshalJSON()

	return WasmExecuteContract{
		Sender:          msg.Sender,
		ContractAddress: msg.Contract,
		RawContractMsg:  rawContractMsg,
		Funds:           msg.Funds,
		GasUsed:         tx.GasUsed,
		Fees:            tx.GetFee(),
		TxHash:          tx.TxHash,
		ExecutedAt:      executedAt,
		Height:          tx.Height,
	}
}
