CREATE TYPE COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TYPE ACCESS_CONFIG AS
(
    permission  INT,
    address     TEXT
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
    sender                  TEXT            NOT NULL,
    byte_code               TEXT            NOT NULL,
    instantiate_permission  ACCESS_CONFIG   NULL,
    code_id                 BIGINT          NOT NULL UNIQUE,
    height                  BIGINT          NOT NULL REFERENCES block (height)
);
CREATE INDEX wasm_code_height_index ON wasm_code (height);


CREATE TABLE wasm_contract
(
    sender                  TEXT            NOT NULL,
    creator                 TEXT            NOT NULL,
    admin                   TEXT            NOT NULL DEFAULT '',
    code_id                 BIGINT          NOT NULL REFERENCES wasm_code (code_id),
    label                   TEXT            NULL,
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    contract_address        TEXT            NOT NULL UNIQUE,
    data                    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    instantiated_at         TIMESTAMP       NOT NULL,
    contract_info_extension JSONB           NOT NULL DEFAULT '{}'::JSONB,
    height                  BIGINT          NOT NULL REFERENCES block (height)
);
CREATE INDEX wasm_contract_height_index ON wasm_contract (height);
CREATE INDEX wasm_contract_creator ON wasm_contract (creator);
CREATE INDEX wasm_contract_contract_address ON wasm_contract (contract_address);


CREATE TABLE wasm_execute_contract
(
    sender                  TEXT            NOT NULL,
    contract_address        TEXT            NOT NULL REFERENCES wasm_contract (contract_address),
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    data                    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    executed_at             TIMESTAMP       NOT NULL,
    height                  BIGINT          NOT NULL REFERENCES block (height)
);
CREATE INDEX execute_contract_height_index ON wasm_execute_contract (height);
CREATE INDEX execute_contract_contract_address ON wasm_execute_contract (contract_address);


CREATE TABLE contract_reward
(
    contract_address           TEXT    NOT NULL REFERENCES wasm_contract (contract_address),
    reward_address             TEXT    NOT NULL,
    developer_address          TEXT    NOT NULL,
    contract_rewards_amount    COIN[]  NOT NULL DEFAULT '{}',
    inflation_rewardsAmount    COIN[]  NOT NULL DEFAULT '{}',
    leftover_rewards_amount    COIN[]  NOT NULL DEFAULT '{}',
    collect_premium            BOOLEAN,
    gas_rebate_to_user         BOOLEAN,
    premium_percentage_charged BIGINT,
    gas_consumed               BIGINT  DEFAULT 0
    dataCalculation            JSONB   NOT NULL DEFAULT '{}'::JSONB,
    dataDistribution           JSONB   NOT NULL DEFAULT '{}'::JSONB,
    height                     BIGINT  NOT NULL REFERENCES block (height),
);
CREATE INDEX contract_reward_contract_address_index ON contract_reward (contract_address);
CREATE INDEX contract_reward_developer_address_index ON contract_reward (developer_address);
CREATE INDEX contract_reward_reward_address_index ON contract_reward (reward_address);