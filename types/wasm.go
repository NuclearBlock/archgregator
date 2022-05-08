package types

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmCode represents the CosmWasm code in x/wasm module
type WasmCode struct {
	Sender                string
	WasmByteCode          []byte
	InstantiatePermission *wasmtypes.AccessConfig
	CodeID                int64
	Height                int64
}

// NewWasmCode allows to build a new x/wasm code instance from wasmtypes.MsgStoreCode
func NewWasmCode(msg *wasmtypes.MsgStoreCode, codeID int64, height int64) WasmCode {
	return WasmCode{
		Sender:                msg.Sender,
		WasmByteCode:          msg.WASMByteCode,
		InstantiatePermission: msg.InstantiatePermission,
		CodeID:                codeID,
		Height:                height,
	}
}

// WasmContract represents the CosmWasm contract in x/wasm module
type WasmContract struct {
	Sender                string
	Creator               string
	Admin                 string
	CodeID                uint64
	Label                 string
	RawContractMsg        []byte
	Funds                 sdk.Coins
	ContractAddress       string
	Data                  string
	InstantiatedAt        time.Time
	ContractInfoExtension wasmtypes.ContractInfoExtension
	Height                int64
}

// NewWasmCode allows to build a new x/wasm contract instance from wasmtypes.MsgStoreCode
func NewWasmContract(
	msg *wasmtypes.MsgInstantiateContract,
	contractAddress string,
	data string,
	instantiatedAt time.Time,
	creator string,
	contractInfoExtension wasmtypes.ContractInfoExtension,
	height int64,
) WasmContract {
	rawContractMsg, _ := msg.Msg.MarshalJSON()

	return WasmContract{
		Sender:                msg.Sender,
		Creator:               creator,
		Admin:                 msg.Admin,
		CodeID:                msg.CodeID,
		Label:                 msg.Label,
		RawContractMsg:        rawContractMsg,
		Funds:                 msg.Funds,
		ContractAddress:       contractAddress,
		Data:                  data,
		InstantiatedAt:        instantiatedAt,
		ContractInfoExtension: contractInfoExtension,
		Height:                height,
	}
}

// WasmExecuteContract represents the CosmWasm execute contract in x/wasm module
type WasmExecuteContract struct {
	Sender          string
	ContractAddress string
	RawContractMsg  []byte
	Funds           sdk.Coins
	Data            string
	ExecutedAt      time.Time
	Height          int64
}

// NewWasmExecuteContract allows to build a new x/wasm execute contract instance from wasmtypes.MsgExecuteContract
func NewWasmExecuteContract(
	msg *wasmtypes.MsgExecuteContract,
	data string,
	executedAt time.Time,
	height int64,
) WasmExecuteContract {
	rawContractMsg, _ := msg.Msg.MarshalJSON()

	return WasmExecuteContract{
		Sender:          msg.Sender,
		ContractAddress: msg.Contract,
		RawContractMsg:  rawContractMsg,
		Funds:           msg.Funds,
		Data:            data,
		ExecutedAt:      executedAt,
		Height:          height,
	}
}

// WasmMigrateContract represents the CosmWasm migrate contract in x/wasm module
type WasmMigrateContract struct {
	Sender          string
	ContractAddress string
	RawContractMsg  []byte
	Data            string
	MigratedAt      time.Time
	Height          int64
}

// NewWasmMigrateContract allows to build a new x/wasm migrate contract instance from wasmtypes.MsgMigrateContract
func NewWasmMigrateContract(
	msg *wasmtypes.MsgMigrateContract,
	data string,
	migratedAt time.Time,
	height int64,
) WasmMigrateContract {
	rawContractMsg, _ := msg.Msg.MarshalJSON()

	return WasmMigrateContract{
		Sender:          msg.Sender,
		ContractAddress: msg.Contract,
		RawContractMsg:  rawContractMsg,
		Data:            data,
		MigratedAt:      migratedAt,
		Height:          height,
	}
}

// WasmUpdateAdminContract represents the CosmWasm update admin contract in x/wasm module
type WasmUpdateAdminContract struct {
	Sender          string
	NewAdmin        string
	ContractAddress string
	UpdatedAt       time.Time
	Height          int64
}

// NewWasmUpdateAdminContract allows to build a new x/wasm update admin contract instance from wasmtypes.MsgUpdateAdmin
func NewWasmUpdateAdminContract(
	msg *wasmtypes.MsgUpdateAdmin,
	updatedAt time.Time,
	height int64,
) WasmUpdateAdminContract {

	return WasmUpdateAdminContract{
		Sender:    msg.Sender,
		NewAdmin:  msg.NewAdmin,
		UpdatedAt: updatedAt,
		Height:    height,
	}
}
