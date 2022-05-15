package parser

import (
	"fmt"
	"strconv"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	database "github.com/nuclearblock/archgregator/database"
	"github.com/nuclearblock/archgregator/node"
	types "github.com/nuclearblock/archgregator/types"
)

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
// The Store Code Event is to upload the contract code on the chain, where a Code ID is returned
func HandleMsgStoreCode(index int, tx *types.Tx, node node.Node, db database.Database) error {

	// Get store code event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeStoreCode)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeInstantiate: %s", err)
	}

	// Get code ID from store code event
	codeIDKey, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyCodeID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyContractAddr: %s", err)
	}

	codeID, err := strconv.ParseUint(codeIDKey, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing code id to uint64: %s", err)
	}

	// Get the code info
	codeInfo, err := node.GetCodeInfo(tx.Height, codeID)
	if err != nil {
		return fmt.Errorf("error while getting contract info: %s", err)
	}

	return db.SaveWasmCode(
		types.NewWasmCode(codeInfo, tx.TxHash, tx.Height),
	)
}

// HandleMsgInstantiateContract allows to properly handle a MsgInstantiateContract
// Instantiate Contract Event instantiates an executable contract with the code previously stored with Store Code Event
func HandleMsgInstantiateContract(index int, tx *types.Tx, msg *wasmtypes.MsgInstantiateContract, node node.Node, db database.Database) error {
	// Get instantiate contract event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeInstantiate)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeInstantiate: %s", err)
	}

	// Get contract address
	contractAddress, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyContractAddr)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyContractAddr: %s", err)
	}

	// Get the contract info
	contractInfo, err := node.GetContractInfo(tx.Height, contractAddress)
	if err != nil {
		return fmt.Errorf("error while getting contract info: %s", err)
	}

	// Get creator address
	creator, err := sdk.AccAddressFromBech32(contractInfo.Creator)
	if err != nil {
		return fmt.Errorf("error while parsing contract creator: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmContract(
		types.NewWasmContract(msg, contractAddress, tx.TxHash, timestamp, creator.String(), tx.Height),
	)
}

// HandleMsgExecuteContract allows to properly handle a MsgExecuteContract
// Execute Event executes an instantiated contract
func HandleMsgExecuteContract(index int, tx *types.Tx, msg *wasmtypes.MsgExecuteContract, db database.Database) error {

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmExecuteContract(
		types.NewWasmExecuteContract(msg, tx, timestamp),
	)
}
