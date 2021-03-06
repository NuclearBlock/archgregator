package postgresql

import (
	"database/sql"
	"strconv"

	"fmt"

	"github.com/nuclearblock/archgregator/logging"

	"github.com/archway-network/archway/app/params"
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
	INSERT INTO block (height, hash, num_txs, total_gas, timestamp)
	VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(sqlStatement,
		block.Height, block.Hash, block.TxNum, block.TotalGas, block.Timestamp,
	)
	return err
}

// SaveWasmCode allows to store the wasm code from MsgStoreCode
func (db *Database) SaveWasmCode(wasmCode types.WasmCode) error {
	stmt := `
	INSERT INTO wasm_code(creator, code_hash, code_id, size, tx_hash, saved_at, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(stmt,
		wasmCode.Creator, wasmCode.CodeHash,
		wasmCode.CodeID, wasmCode.Size, wasmCode.TxHash,
		wasmCode.SavedAt, wasmCode.Height,
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
	(sender, creator, admin, code_id, label, raw_contract_message, funds, contract_address, tx_hash, instantiated_at, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
	ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(stmt,
		wasmContract.Sender, wasmContract.Creator, wasmContract.Admin, wasmContract.CodeID, wasmContract.Label, string(wasmContract.RawContractMsg),
		pq.Array(dbtypes.NewDbCoins(wasmContract.Funds)), wasmContract.ContractAddress, wasmContract.TxHash,
		wasmContract.InstantiatedAt, wasmContract.Height,
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
	(sender, contract_address, raw_contract_message, funds, gas_used, fees_denom, fees_amount, tx_hash, executed_at, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	ON CONFLICT DO NOTHING`

	denom := "utorii"

	_, err := db.Sql.Exec(stmt,
		executeContract.Sender,
		executeContract.ContractAddress,
		executeContract.RawContractMsg,
		pq.Array(dbtypes.NewDbCoins(executeContract.Funds)),
		executeContract.GasUsed,
		denom,
		executeContract.Fees.AmountOf(denom).String(),
		executeContract.TxHash,
		executeContract.ExecutedAt,
		executeContract.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving wasm contract: %s", err)
	}

	return nil
}

func (db *Database) SaveContractRewardCalculation(contractRewardCalculation types.ContractRewardCalculation) error {

	stmt := `
	INSERT INTO contract_reward 
	(contract_address, reward_address, developer_address, gas_consumed, contract_rewards_denom, contract_rewards_amount, inflation_rewards_amount, gas_rebate_to_user, collect_premium, premium_percentage_charged, reward_date, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
	ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(
		stmt,
		contractRewardCalculation.ContractAddress,
		contractRewardCalculation.RewardAddress,
		contractRewardCalculation.DeveloperAddress,
		strconv.FormatUint(contractRewardCalculation.GasConsumed, 10),
		contractRewardCalculation.ContractRewards.Denom,
		contractRewardCalculation.ContractRewards.Amount.String(),
		contractRewardCalculation.InflationRewards.Amount.String(),
		contractRewardCalculation.GasRebateToUser,
		contractRewardCalculation.CollectPremium,
		contractRewardCalculation.PremiumPercentageCharged,
		contractRewardCalculation.RewardDate,
		contractRewardCalculation.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving contract reward into DB: %s, query=", err)
	}
	return nil
}

func (db *Database) SaveContractRewardDistribution(contractRewardDistribution types.ContractRewardDistribution) error {

	stmt := `UPDATE contract_reward SET 
	distributed_rewards_amount = $1, leftover_rewards_amount = $2 
	WHERE reward_address = $3 AND height = $4 `

	_, err := db.Sql.Exec(
		stmt,
		contractRewardDistribution.DistributedRewards.Amount.String(),
		contractRewardDistribution.LeftoverRewards.Amount.String(),
		contractRewardDistribution.RewardAddress,
		contractRewardDistribution.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving contract distribution rewards: %s", err)
	}
	return nil
}

func (db *Database) SaveGasTrackerContractMetadata(gastrackerContractMetadata types.GasTrackerContractMetadata) error {

	stmt := `INSERT INTO contract_metadata 
	(contract_address, reward_address, developer_address, collect_premium, gas_rebate_to_user, premium_percentage_charged, tx_hash, saved_at, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
	ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(
		stmt,
		gastrackerContractMetadata.ContractAddress,
		gastrackerContractMetadata.Metadata.RewardAddress,
		gastrackerContractMetadata.Metadata.DeveloperAddress,
		gastrackerContractMetadata.Metadata.CollectPremium,
		gastrackerContractMetadata.Metadata.GasRebateToUser,
		gastrackerContractMetadata.Metadata.PremiumPercentageCharged,
		gastrackerContractMetadata.TxHash,
		gastrackerContractMetadata.SavedAt,
		gastrackerContractMetadata.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving contract metadata: %s", err)
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
