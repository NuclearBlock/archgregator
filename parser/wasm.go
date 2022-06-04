package parser

import (
	"fmt"
	"strconv"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	database "github.com/nuclearblock/archgregator/database"
	"github.com/nuclearblock/archgregator/node"
	types "github.com/nuclearblock/archgregator/types"
)

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
// The Store Code Event is to upload the contract code on the chain, where a Code ID is returned
func HandleMsgStoreCode(index int, tx *types.Tx, msg *wasmtypes.MsgStoreCode, node node.Node, db database.Database) error {

	var codeID uint64
	var creator, codeHash string
	var codeSize int

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

	codeID, err = strconv.ParseUint(codeIDKey, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing code id to uint64: %s", err)
	}

	creator = msg.Sender
	codeSize, codeHash, err = types.GetCodeData(msg.WASMByteCode)
	if err != nil {
		codeSize = 0
		codeHash = ""
	}

	// // Get the code info
	// codeInfo, err := node.GetCodeInfo(tx.Height, codeID)
	// if err != nil {
	// 	// For some reason sometimes we cant get Code Info via wasmClient query
	// 	// If there will bw an error - we'll calculate wasm code size and hash in runtime
	// 	creator = msg.Sender
	// 	codeSize, codeHash, err = types.GetCodeData(msg.WASMByteCode)
	// 	if err != nil {
	// 		codeSize = 0
	// 		codeHash = ""
	// 	}
	// } else {
	// 	creator = codeInfo.Creator,
	// 	codeHash = codeInfo.DataHash.String(),
	// 	codeID = codeInfo.CodeID,
	// 	codeSize = codeInfo.Size(),
	// }

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmCode(
		types.NewWasmCode(codeID, creator, codeSize, codeHash, tx.TxHash, timestamp, tx.Height),
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

	creator := msg.Sender

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmContract(
		types.NewWasmContract(msg, contractAddress, tx.TxHash, timestamp, creator, tx.Height),
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
