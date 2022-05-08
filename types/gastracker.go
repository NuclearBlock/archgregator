package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

const (
	EVENT_ContractRewardCalculationEvent = "archway.gastracker.v1.ContractRewardCalculationEvent"
	EVENT_RewardDistributionEvent        = "archway.gastracker.v1.RewardDistributionEvent"

	EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS  = EVENT_ContractRewardCalculationEvent + ".contract_address"
	EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS  = EVENT_ContractRewardCalculationEvent + ".contract_rewards"
	EVENT_ContractRewardCalculationEvent_GAS_CONSUMED      = EVENT_ContractRewardCalculationEvent + ".gas_consumed"
	EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS = EVENT_ContractRewardCalculationEvent + ".inflation_rewards"
	EVENT_ContractRewardCalculationEvent_METADATA          = EVENT_ContractRewardCalculationEvent + ".metadata"
	EVENT_RewardDistributionEvent_CONTRACT_REWARDS         = EVENT_RewardDistributionEvent + ".contract_rewards"
	EVENT_RewardDistributionEvent_LEFTOVER_REWARDS         = EVENT_RewardDistributionEvent + ".leftover_rewards"
	EVENT_RewardDistributionEvent_REWARD_ADDRESS           = EVENT_RewardDistributionEvent + ".reward_address"

	EVENT_FIELD_DENOM                      = "denom"
	EVENT_FIELD_AMOUNT                     = "amount"
	EVENT_FIELD_DEVELOPER_ADDRESS          = "developer_address"
	EVENT_FIELD_REWARD_ADDRESS             = "reward_address"
	EVENT_FIELD_GAS_REBATE_TO_USER         = "gas_rebate_to_user"
	EVENT_FIELD_COLLECT_PREMIUM            = "collect_premium"
	EVENT_FIELD_PREMIUM_PERCENTAGE_CHARGED = "premium_percentage_charged"
)

type GasTrackerReward struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type ContractReward struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string

	// ??? in DB we changed it to varchar(50) as postgresql does not support uint64
	GasConsumed uint64
	// ??? For sake of simplicity, we consider only one denom per record
	ContractRewards  GasTrackerReward
	InflationRewards GasTrackerReward
	LeftoverRewards  GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged uint64

	MetadataJson string
	Height       uint64
}

// NewContractReward allows to build a new smart contract reward instance from archway.gastracker event
func NewContractReward(evr *coretypes.ResultEvent) (*ContractReward, error) {
	var cr ContractReward

	if _, ok := evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]; !ok {
		// Nothing to process
		return nil, nil
	}

	b := evr.Data.(tmTypes.EventDataNewBlock)
	// The gastracking is processed in the next beginBlock
	cr.Height = uint64(b.Block.Height) - 1

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]) > 0 {
		cr.ContractAddress = strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS][0], "\"")
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA]) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]), &metadata); err != nil {
			return nil, err
		}

		cr.RewardAddress = metadata[EVENT_FIELD_REWARD_ADDRESS].(string)
		cr.DeveloperAddress = metadata[EVENT_FIELD_DEVELOPER_ADDRESS].(string)
		cr.GasRebateToUser = metadata[EVENT_FIELD_GAS_REBATE_TO_USER].(bool)
		cr.CollectPremium = metadata[EVENT_FIELD_COLLECT_PREMIUM].(bool)
		cr.MetadataJson = evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]

		intValue, err := strconv.ParseUint(metadata[EVENT_FIELD_PREMIUM_PERCENTAGE_CHARGED].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		cr.PremiumPercentageCharged = intValue
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED]) > 0 {

		intValue, err := strconv.ParseUint(strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED][0], "\""), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error in Unmarshaling '%s': %v", EVENT_ContractRewardCalculationEvent_GAS_CONSUMED, err)
		}
		cr.GasConsumed = intValue
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS]) > 0 {
		var err error
		cr.ContractRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS]) > 0 {
		var err error
		cr.InflationRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	if len(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS]) > 0 {
		var err error
		cr.LeftoverRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	return &cr, nil
}

// func (c *ContractRecord) getDBRow() database.RowType {
// 	return database.RowType{
// 		database.FIELD_CONTRACTS_CONTRACT_ADDRESS:           c.ContractAddress,
// 		database.FIELD_CONTRACTS_REWARD_ADDRESS:             c.RewardAddress,
// 		database.FIELD_CONTRACTS_DEVELOPER_ADDRESS:          c.DeveloperAddress,
// 		database.FIELD_CONTRACTS_BLOCK_HEIGHT:               c.BlockHeight,
// 		database.FIELD_CONTRACTS_GAS_CONSUMED:               fmt.Sprintf("%d", c.GasConsumed),
// 		database.FIELD_CONTRACTS_REWARDS_DENOM:              c.ContractRewards.Denom,
// 		database.FIELD_CONTRACTS_CONTRACT_REWARDS_AMOUNT:    c.ContractRewards.Amount,
// 		database.FIELD_CONTRACTS_INFLATION_REWARDS_AMOUNT:   c.InflationRewards.Amount,
// 		database.FIELD_CONTRACTS_LEFTOVER_REWARDS_AMOUNT:    c.LeftoverRewards.Amount,
// 		database.FIELD_CONTRACTS_COLLECT_PREMIUM:            c.CollectPremium,
// 		database.FIELD_CONTRACTS_GAS_REBATE_TO_USER:         c.GasRebateToUser,
// 		database.FIELD_CONTRACTS_PREMIUM_PERCENTAGE_CHARGED: c.PremiumPercentageCharged,
// 		database.FIELD_CONTRACTS_METADATA_JSON:              c.MetadataJson,
// 	}
// }

func getGasTrackerRewardFromString(str string) (GasTrackerReward, error) {

	// Let's make it an array if not, to keep compatibility
	if !strings.HasPrefix(str, "[") {
		str = "[" + str + "]"
	}

	var tmpMapArr []map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmpMapArr); err != nil {
		return GasTrackerReward{}, err
	}

	if len(tmpMapArr) == 0 {
		return GasTrackerReward{}, fmt.Errorf("no GasTrackerReward found")
	}
	tmpMap := tmpMapArr[0]

	numValue, err := strconv.ParseFloat(tmpMap[EVENT_FIELD_AMOUNT].(string), 64)
	if err != nil {
		return GasTrackerReward{}, err
	}

	return GasTrackerReward{
		Denom:  tmpMap[EVENT_FIELD_DENOM].(string),
		Amount: numValue,
	}, nil
}
