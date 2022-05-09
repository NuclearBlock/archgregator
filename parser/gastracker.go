package parser

import (
	"encoding/json"
	"fmt"

	"strconv"
	"strings"

	//"google.golang.org/grpc/internal/metadata"
	// coretypes "github.com/tendermint/tendermint/rpc/core/types"
	// tmTypes "github.com/tendermint/tendermint/types"
	database "github.com/nuclearblock/archgregator/database"
	types "github.com/nuclearblock/archgregator/types"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
)

// NewContractReward allows to build a new smart contract reward instance from archway.gastracker event
func HandleReward(event *tmabcitypes.Event, height uint64, db database.Database) error {

	// We have to check if current event is Gastracker module reward event

	// Firstly, try to catch reward calculation event
	if strings.Contains(event.Type, "archway.gastracker.v1.ContractRewardCalculationEvent") {
		fmt.Printf("event = %s\n", event.Type)

		var metadataReward types.MetadataReward
		// Get all event attributes
		eventAttributes := event.GetAttributes()
		// Handle all the atribution inside the event type
		for _, attribute := range eventAttributes {
			switch string(attribute.Key) {
			case "contract_address":
				contractAddress := HandleAddress(attribute.Value)
			case "gas_consumed":
				gasConsumed, err := HandleGas(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing gas consumed (calculation event): %s", err)
				}
			case "contract_rewards":
				contractRewards, err := HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing contract rewards (calculation event): %s", err)
				}
			case "inflation_rewards":
				inflationRewards, err := HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing inflation rewards (calculation event): %s", err)
				}
			case "metadata":
				metadataReward, err := HandleMetadata(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing metafata (calculation event): %s", err)
				}
			}
		}

		// eventJson is strng JSON of ContractRewardCalculationEvent
		eventJson, err := GetEventJson(event)
		if err != nil {
			return fmt.Errorf("error while parsing event JSON (calculation event): %s", err)
		}

		// We have to decrement target block height,
		// because reward is always processed in the next beginBlock
		rewardCaculationHeight := height - 1

		return db.SaveContractRewardCalculation(
			types.NewContractRewardCalculation(
				contractAddress,
				metadataReward.RewardAddress,
				metadataReward.DeveloperAddress,
				gasConsumed,
				contractRewards,
				inflationRewards,
				metadataReward.CollectPremium,
				metadataReward.GasRebateToUser,
				metadataReward.PremiumPercentageCharged,
				eventJson,
				rewardCaculationHeight,
			),
		)

	}

	// Now try to catch reward distribution event
	// Because rewards distriburion happens in another event, we have to collect this data
	// and update db row previously added with 'ContractRewardCalculationEvent'
	if strings.Contains(event.Type, "archway.gastracker.v1.RewardDistributionEvent") {
		fmt.Printf("event = %s\n", event.Type)
		// Get all event attributes
		eventAttributes := event.GetAttributes()
		// Handle all the atribution inside the event type
		for _, attribute := range eventAttributes {
			switch string(attribute.Key) {
			case "reward_address":
				contractAddress := HandleAddress(attribute.Value)
			case "contract_rewards":
				contractDistributionRewards, err := HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing contract rewards (distribution event): %s", err)
				}
			case "leftover_rewards":
				leftoverRewards, err := HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing leftover rewards (distribution event): %s", err)
				}
			}
		}

		// eventJson is strng JSON of ContractRewardCalculationEvent
		eventJson, err := GetEventJson(event)
		if err != nil {
			return fmt.Errorf("error while parsing event JSON (calculation event): %s", err)
		}

		// We have to decrement target block height,
		// because reward is always processed in the next beginBlock
		rewardDistributionHeight := height - 1

		return db.SaveContractRewardDistribution(
			types.NewContractRewardDistribution(
				contractAddress,
				contractDistributionRewards,
				leftoverRewards,
				eventJson,
				rewardDistributionHeight,
			),
		)

	}

	return nil
}

func HandleAddress(value []byte) string {
	return string(value)
}

func HandleGas(value []byte) (uint64, error) {
	return strconv.ParseUint(string(value), 10, 64)
}

func HandleRewards(value []byte) (types.GasTrackerReward, error) {
	return getGasTrackerRewardFromString(string(value))
}

func HandleMetadata(value []byte) (types.MetadataReward, error) {
	var metadata types.MetadataReward
	json.Unmarshal([]byte(value), &metadata)
}

func GetEventJson(event *tmabcitypes.Event) (string, error) {
	return "", nil
}

func getGasTrackerRewardFromString(str string) (types.GasTrackerReward, error) {
	// Let's make it an array if not, to keep compatibility
	if !strings.HasPrefix(str, "[") {
		str = "[" + str + "]"
	}

	var tmpMapArr []map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmpMapArr); err != nil {
		return types.GasTrackerReward{}, err
	}

	if len(tmpMapArr) == 0 {
		return types.GasTrackerReward{}, fmt.Errorf("no GasTrackerReward found")
	}
	tmpMap := tmpMapArr[0]

	numValue, err := strconv.ParseFloat(tmpMap["amount"].(string), 64)
	if err != nil {
		return types.GasTrackerReward{}, err
	}

	return types.GasTrackerReward{
		Denom:  tmpMap["denom"].(string),
		Amount: numValue,
	}, nil
}