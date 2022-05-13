package postgresql

import (
	"database/sql"
	// "encoding/base64"
	"fmt"

	"github.com/nuclearblock/archgregator/logging"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/lib/pq"

	_ "github.com/lib/pq" // nolint

	"github.com/nuclearblock/archgregator/database"
	dbtypes "github.com/nuclearblock/archgregator/database/types"
	"github.com/nuclearblock/archgregator/types"
)

// Builder creates a database connection with the given database connection info
// from config. It returns a database connection handle or an error if the
// connection fails.
func Builder(ctx *database.Context) (database.Database, error) {
	sslMode := "disable"
	if ctx.Cfg.SSLMode != "" {
		sslMode = ctx.Cfg.SSLMode
	}

	schema := "public"
	if ctx.Cfg.Schema != "" {
		schema = ctx.Cfg.Schema
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s sslmode=%s search_path=%s",
		ctx.Cfg.Host, ctx.Cfg.Port, ctx.Cfg.Name, ctx.Cfg.User, sslMode, schema,
	)

	if ctx.Cfg.Password != "" {
		connStr += fmt.Sprintf(" password=%s", ctx.Cfg.Password)
	}

	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set max open connections
	postgresDb.SetMaxOpenConns(ctx.Cfg.MaxOpenConnections)
	postgresDb.SetMaxIdleConns(ctx.Cfg.MaxIdleConnections)

	return &Database{
		Sql:            postgresDb,
		EncodingConfig: ctx.EncodingConfig,
		Logger:         ctx.Logger,
	}, nil
}

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	Sql            *sql.DB
	EncodingConfig *params.EncodingConfig
	Logger         logging.Logger
}

// HasBlock implements database.Database
func (db *Database) HasBlock(height int64) (bool, error) {
	var res bool
	err := db.Sql.QueryRow(`SELECT EXISTS(SELECT 1 FROM block WHERE height = $1);`, height).Scan(&res)
	return res, err
}

// SaveBlock implements database.Database
func (db *Database) SaveBlock(block *types.Block) error {
	sqlStatement := `
INSERT INTO block (height, hash, num_txs, total_gas, proposer_address, timestamp)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`

	proposerAddress := sql.NullString{Valid: len(block.ProposerAddress) != 0, String: block.ProposerAddress}
	_, err := db.Sql.Exec(sqlStatement,
		block.Height, block.Hash, block.TxNum, block.TotalGas, proposerAddress, block.Timestamp,
	)
	return err
}

// SaveTx implements database.Database
func (db *Database) SaveTx(tx *types.Tx) error {
	//TO-DO
	return nil
}

// SaveMessage implements database.Database
func (db *Database) SaveMessage(msg *types.Message) error {
	//TO-DO
	return nil
}

// SaveWasmCode allows to store the wasm code from MsgStoreCode
func (db *Database) SaveWasmCode(wasmCode types.WasmCode) error {

	stmt := `
INSERT INTO wasm_code(sender, byte_code, instantiate_permission, code_id, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT DO NOTHING`

	// TO-DO: check if string(wasmCode.WasmByteCode) saved as string in DB

	_, err := db.Sql.Exec(stmt,
		wasmCode.Sender, string(wasmCode.WasmByteCode),
		pq.Array(dbtypes.NewDbAccessConfig(wasmCode.InstantiatePermission)),
		wasmCode.CodeID, wasmCode.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving wasm code: %s", err)
	}

	return nil
}

// SaveWasmContract allows to store the wasm contract from MsgInstantiateContract
func (db *Database) SaveWasmContract(wasmContract types.WasmContract) error {

	stmt := `
INSERT INTO wasm_contract 
(sender, admin, code_id, label, raw_contract_message, funds, contract_address, data, instantiated_at, contract_info_extension, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
ON CONFLICT DO NOTHING`

	ExtensionBz, err := db.EncodingConfig.Marshaler.MarshalJSON(wasmContract.ContractInfoExtension)
	if err != nil {
		return fmt.Errorf("error while marshaling contract info extension: %s", err)
	}

	// TO-DO: check if the below is stored as Json in DB:
	// - Data
	// - ContractInfoExtension
	// - RawContractMsg

	_, err = db.Sql.Exec(stmt,
		wasmContract.Sender, wasmContract.Admin, wasmContract.CodeID, wasmContract.Label, string(wasmContract.RawContractMsg),
		pq.Array(dbtypes.NewDbCoins(wasmContract.Funds)), wasmContract.ContractAddress, wasmContract.Data,
		wasmContract.InstantiatedAt, string(ExtensionBz), wasmContract.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving wasm contract: %s", err)
	}

	return nil
}

// SaveWasmExecuteContract allows to store the wasm contract from MsgExecuteeContract
func (db *Database) SaveWasmExecuteContract(executeContract types.WasmExecuteContract) error {

	stmt := `
INSERT INTO wasm_execute_contract 
(sender, contract_address, raw_contract_message, funds, tx_hash, executed_at, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT DO NOTHING`

	// TO-DO: check if the below is stored as Json in DB:
	// - Data

	_, err := db.Sql.Exec(stmt,
		executeContract.Sender,
		executeContract.ContractAddress,
		executeContract.RawContractMsg,
		pq.Array(dbtypes.NewDbCoins(executeContract.Funds)),
		executeContract.TxHash,
		executeContract.ExecutedAt,
		executeContract.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving wasm contract: %s", err)
	}

	return nil
}

func (db *Database) UpdateContractWithMsgMigrateContract(
	sender string,
	contractAddress string,
	codeID uint64,
	rawContractMsg []byte,
	data string,
) error {

	stmt := `UPDATE wasm_contract SET 
sender = $1, code_id = $2, raw_contract_message = $3, data = $4 
WHERE contract_address = $5 `

	// TO-DO: check if the below is stored as Json in DB:
	// - rawContractMsg
	// - Data

	_, err := db.Sql.Exec(stmt,
		sender, codeID, string(rawContractMsg), data,
		contractAddress,
	)
	if err != nil {
		return fmt.Errorf("error while updating wasm contract from contract migration: %s", err)

	}
	return nil
}

func (db *Database) UpdateContractAdmin(sender string, contractAddress string, newAdmin string) error {

	stmt := `UPDATE wasm_contract SET 
sender = $1, admin = $2 WHERE contract_address = $2 `

	_, err := db.Sql.Exec(stmt, sender, newAdmin, contractAddress)
	if err != nil {
		return fmt.Errorf("error while updating wsm contract admin: %s", err)
	}
	return nil
}

func (db *Database) SaveContractRewardCalculation(contractRewardCalculation types.ContractRewardCalculation) error {

	fmt.Printf("contractRewardCalculation: %+v\n", contractRewardCalculation)

	stmt := `
INSERT INTO contract_reward 
(contract_address, reward_address, developer_address, contract_rewards, inflation_rewards, collect_premium, gas_rebate_to_user, premium_percentage_charged, gas_consumed, metadata_json, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
ON CONFLICT DO NOTHING`

	// TO-DO: check if the below is stored as Json in DB:
	// - Data

	_, err := db.Sql.Exec(
		stmt,
		contractRewardCalculation.ContractAddress,
		contractRewardCalculation.RewardAddress,
		contractRewardCalculation.DeveloperAddress,
		pq.Array(dbtypes.NewDbCoinsFromGastracker(contractRewardCalculation.ContractRewards)),
		pq.Array(dbtypes.NewDbCoinsFromGastracker(contractRewardCalculation.InflationRewards)),
		contractRewardCalculation.CollectPremium,
		contractRewardCalculation.GasRebateToUser,
		contractRewardCalculation.PremiumPercentageCharged,
		string(contractRewardCalculation.GasConsumed),
		contractRewardCalculation.MetadataJson,
		contractRewardCalculation.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving contract reward into DB: %s, query=", err)
	}
	return nil
}

func (db *Database) SaveContractRewardDistribution(contractRewardDistribution types.ContractRewardDistribution) error {

	fmt.Printf("contractRewardDistribution: %+v\n", contractRewardDistribution)

	stmt := `UPDATE contract_reward SET 
	distributed_rewards = $1, leftover_rewards = $2 
	WHERE reward_address = $3 AND height = $4 `

	_, err := db.Sql.Exec(
		stmt,
		pq.Array(dbtypes.NewDbCoinsFromGastracker(contractRewardDistribution.DistributedRewards)),
		pq.Array(dbtypes.NewDbCoinsFromGastracker(contractRewardDistribution.LeftoverRewards)),
		contractRewardDistribution.RewardAddress,
		contractRewardDistribution.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving contract distribution rewards: %s", err)
	}
	return nil
}

// Close implements database.Database
func (db *Database) Close() {
	err := db.Sql.Close()
	if err != nil {
		db.Logger.Error("error while closing connection", "err", err)
	}
}
