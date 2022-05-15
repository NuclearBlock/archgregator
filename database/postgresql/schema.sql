CREATE TYPE COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE block
(
    height           BIGINT UNIQUE PRIMARY KEY,
    hash             TEXT NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        BIGINT  DEFAULT 0,
    proposer_address TEXT,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX block_height_index ON block (height);
CREATE INDEX block_hash_index ON block (hash);
CREATE INDEX block_proposer_address_index ON block (proposer_address);


CREATE TABLE wasm_code
(
    creator                 TEXT            NOT NULL,
    code_hash               TEXT            NOT NULL,
    code_id                 BIGINT          NOT NULL UNIQUE,
    size                    INT             NOT NULL,
    tx_hash                 TEXT            NOT NULL,
    height                  BIGINT          NOT NULL
);
CREATE INDEX wasm_code_height_index ON wasm_code (height);


CREATE TABLE wasm_contract
(
    sender                  TEXT            NOT NULL,
    creator                 TEXT            NOT NULL,
    admin                   TEXT            NOT NULL DEFAULT '',
    code_id                 BIGINT          NOT NULL,
    label                   TEXT            NULL,
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    contract_address        TEXT            NOT NULL UNIQUE,
    tx_hash                 TEXT            NOT NULL,
    instantiated_at         TIMESTAMP       NOT NULL,
    height                  BIGINT          NOT NULL
);
CREATE INDEX wasm_contract_height_index ON wasm_contract (height);
CREATE INDEX wasm_contract_creator ON wasm_contract (creator);
CREATE INDEX wasm_contract_contract_address ON wasm_contract (contract_address);


CREATE TABLE wasm_execute_contract
(
    sender                  TEXT            NOT NULL,
    contract_address        TEXT            NOT NULL,
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    tx_hash                 TEXT            NOT NULL,
    executed_at             TIMESTAMP       NOT NULL,
    height                  BIGINT          NOT NULL
);
CREATE INDEX execute_contract_height_index ON wasm_execute_contract (height);
CREATE INDEX execute_contract_contract_address ON wasm_execute_contract (contract_address);


CREATE TABLE contract_metadata
(
    contract_address           TEXT    NOT NULL,
    reward_address             TEXT    NOT NULL,
    developer_address          TEXT    NOT NULL,
    collect_premium            BOOLEAN,
    gas_rebate_to_user         BOOLEAN,
    premium_percentage_charged BIGINT,
    metadata_json              JSONB   NOT NULL DEFAULT '{}'::JSONB,
    tx_hash                    TEXT    NOT NULL,              
    saved_at                   TIMESTAMP  NOT NULL,
    height                     BIGINT  NOT NULL
);
CREATE INDEX contract_metadata_height_index ON contract_metadata (height);
CREATE INDEX contract_metadata_contract_address_index ON contract_metadata (contract_address);
CREATE INDEX contract_metadata_developer_address_index ON contract_metadata (developer_address);
CREATE INDEX contract_metadata_reward_address_index ON contract_metadata (reward_address);


CREATE TABLE contract_reward
(
    contract_address           TEXT    NOT NULL,
    reward_address             TEXT    NOT NULL,
    developer_address          TEXT    NOT NULL,
    gas_consumed               TEXT    DEFAULT 0,
    contract_rewards           COIN[]  NOT NULL DEFAULT '{}',
    inflation_rewards          COIN[]  NOT NULL DEFAULT '{}',
    distributed_rewards        COIN[]  NOT NULL DEFAULT '{}',
    leftover_rewards           COIN[]  NOT NULL DEFAULT '{}',
    metadata_json              JSONB   NOT NULL DEFAULT '{}'::JSONB,
    height                     BIGINT  NOT NULL
);
CREATE INDEX contract_reward_height_index ON contract_reward (height);
CREATE INDEX contract_reward_contract_address_index ON contract_reward (contract_address);
