package parser

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	database "github.com/nuclearblock/archgregator/database"
	"github.com/nuclearblock/archgregator/node"
	types "github.com/nuclearblock/archgregator/types"
)

// HandleMsg implements modules.MessageModule
func HandleWasmMsg(index int, msg sdk.Msg, tx *types.Tx, node node.Node, db database.Database) error {
	fmt.Println(tx)
	fmt.Println(msg.String())

	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *wasmtypes.MsgStoreCode:
		return HandleMsgStoreCode(index, tx, cosmosMsg, db)
	case *wasmtypes.MsgInstantiateContract:
		return HandleMsgInstantiateContract(index, tx, cosmosMsg, node, db)
	case *wasmtypes.MsgExecuteContract:
		return HandleMsgExecuteContract(index, tx, cosmosMsg, db)
	case *wasmtypes.MsgMigrateContract:
		return HandleMsgMigrateContract(index, tx, cosmosMsg, db)
	case *wasmtypes.MsgUpdateAdmin:
		return HandleMsgUpdateAdmin(cosmosMsg, db)
	case *wasmtypes.MsgClearAdmin:
		return HandleMsgClearAdmin(cosmosMsg, db)
	}

	return nil
}

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
// The Store Code Event is to upload the contract code on the chain, where a Code ID is returned
func HandleMsgStoreCode(index int, tx *types.Tx, msg *wasmtypes.MsgStoreCode, db database.Database) error {
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

	codeID, err := strconv.ParseInt(codeIDKey, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing code id to uint64: %s", err)
	}

	return db.SaveWasmCode(
		types.NewWasmCode(msg, codeID, tx.Height),
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

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyResultDataHex: %s", err)
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	// Get the contract info
	contractInfo, err := node.GetContractInfo(tx.Height, contractAddress)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmContract(
		types.NewWasmContract(msg, contractAddress, string(resultDataBz), timestamp, contractInfo.Creator, contractInfo.Extension, tx.Height),
	)
}

// HandleMsgExecuteContract allows to properly handle a MsgExecuteContract
// Execute Event executes an instantiated contract
func HandleMsgExecuteContract(index int, tx *types.Tx, msg *wasmtypes.MsgExecuteContract, db database.Database) error {
	// Get Execute Contract event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeExecute)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeExecute: %s", err)
	}

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyResultDataHex: %s", err)
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveWasmExecuteContract(
		types.NewWasmExecuteContract(msg, string(resultDataBz), timestamp, tx.Height),
	)
}

// HandleMsgMigrateContract allows to properly handle a MsgMigrateContract
// Migrate Contract Event upgrade the contract by updating code ID generated from new Store Code Event
func HandleMsgMigrateContract(index int, tx *types.Tx, msg *wasmtypes.MsgMigrateContract, db database.Database) error {
	// Get Migrate Contract event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeMigrate)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeMigrate: %s", err)
	}

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyResultDataHex: %s", err)
	}
	resultDataBz, err := base64.StdEncoding.DecodeString(resultData)
	if err != nil {
		return fmt.Errorf("error while decoding result data: %s", err)
	}

	return db.UpdateContractWithMsgMigrateContract(msg.Sender, msg.Contract, msg.CodeID, msg.Msg, string(resultDataBz))
}

// HandleMsgUpdateAdmin allows to properly handle a MsgUpdateAdmin
// Update Admin Event updates the contract admin who can migrate the wasm contract
func HandleMsgUpdateAdmin(msg *wasmtypes.MsgUpdateAdmin, db database.Database) error {
	return db.UpdateContractAdmin(msg.Sender, msg.Contract, msg.NewAdmin)
}

// HandleMsgClearAdmin allows to properly handle a MsgClearAdmin
// Clear Admin Event clears the admin which make the contract no longer migratable
func HandleMsgClearAdmin(msg *wasmtypes.MsgClearAdmin, db database.Database) error {
	return db.UpdateContractAdmin(msg.Sender, msg.Contract, "")
}
