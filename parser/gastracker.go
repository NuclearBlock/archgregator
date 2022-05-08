package parser

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"github.com/nuclearblock/archgregator/types/gastracker"

// 	coretypes "github.com/tendermint/tendermint/rpc/core/types"
// 	tmTypes "github.com/tendermint/tendermint/types"
// )

// // // NewContractReward allows to build a new smart contract reward instance from archway.gastracker event
// // func NewContractReward(evr *coretypes.ResultEvent) ContractReward (*ContractReward, error) {
// // 	var cr ContractReward

// // 	if _, ok := evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]; !ok {
// // 		// Nothing to process
// // 		return nil, nil
// // 	}

// // 	b := evr.Data.(tmTypes.EventDataNewBlock)
// // 	// The gastracking is processed in the next beginBlock
// // 	cr.Height = uint64(b.Block.Height) - 1

// // 	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]) > 0 {
// // 		cr.ContractAddress = strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS][0], "\"")
// // 	}

// // 	if len(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA]) > 0 {
// // 		var metadata map[string]interface{}
// // 		if err := json.Unmarshal([]byte(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]), &metadata); err != nil {
// // 			return nil, err
// // 		}

// // 		cr.RewardAddress = metadata[EVENT_FIELD_REWARD_ADDRESS].(string)
// // 		cr.DeveloperAddress = metadata[EVENT_FIELD_DEVELOPER_ADDRESS].(string)
// // 		cr.GasRebateToUser = metadata[EVENT_FIELD_GAS_REBATE_TO_USER].(bool)
// // 		cr.CollectPremium = metadata[EVENT_FIELD_COLLECT_PREMIUM].(bool)
// // 		cr.MetadataJson = evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]

// // 		intValue, err := strconv.ParseUint(metadata[EVENT_FIELD_PREMIUM_PERCENTAGE_CHARGED].(string), 10, 64)
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 		cr.PremiumPercentageCharged = intValue
// // 	}

// // 	if len(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED]) > 0 {

// // 		intValue, err := strconv.ParseUint(strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED][0], "\""), 10, 64)
// // 		if err != nil {
// // 			return nil, fmt.Errorf("error in Unmarshaling '%s': %v", EVENT_ContractRewardCalculationEvent_GAS_CONSUMED, err)
// // 		}
// // 		cr.GasConsumed = intValue
// // 	}

// // 	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS]) > 0 {
// // 		var err error
// // 		cr.ContractRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS][0])
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 	}

// // 	if len(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS]) > 0 {
// // 		var err error
// // 		cr.InflationRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS][0])
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 	}

// // 	if len(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS]) > 0 {
// // 		var err error
// // 		cr.LeftoverRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS][0])
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 	}

// // 	return &cr, nil
// // }

// // func getGasTrackerRewardFromString(str string) (GasTrackerReward, error) {

// // 	// Let's make it an array if not, to keep compatibility
// // 	if !strings.HasPrefix(str, "[") {
// // 		str = "[" + str + "]"
// // 	}

// // 	var tmpMapArr []map[string]interface{}
// // 	if err := json.Unmarshal([]byte(str), &tmpMapArr); err != nil {
// // 		return GasTrackerReward{}, err
// // 	}

// // 	if len(tmpMapArr) == 0 {
// // 		return GasTrackerReward{}, fmt.Errorf("no GasTrackerReward found")
// // 	}
// // 	tmpMap := tmpMapArr[0]

// // 	numValue, err := strconv.ParseFloat(tmpMap[EVENT_FIELD_AMOUNT].(string), 64)
// // 	if err != nil {
// // 		return GasTrackerReward{}, err
// // 	}

// // 	return GasTrackerReward{
// // 		Denom:  tmpMap[EVENT_FIELD_DENOM].(string),
// // 		Amount: numValue,
// // 	}, nil
// // }