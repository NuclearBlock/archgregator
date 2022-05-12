package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	database "github.com/nuclearblock/archgregator/database"
	types "github.com/nuclearblock/archgregator/types"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
)

// NewContractReward allows to build a new smart contract reward instance from archway.gastracker event
func HandleReward(event *tmabcitypes.Event, height int64, db database.Database) error {

	// We have to check if the current event is Gastracker module reward event

	// Firstly, try to catch reward calculation event
	if strings.Contains(event.Type, "archway.gastracker.v1.ContractRewardCalculationEvent") {

		var contractAddress string
		var gasConsumed string
		var contractRewards, inflationRewards sdk.Coins
		var metadataCalculationReward *types.MetadataReward
		//var metadata map[string]interface{}
		var err error

		// Get all event attributes
		eventAttributes := event.GetAttributes()
		// Handle all the atribution inside the event type
		for _, attribute := range eventAttributes {
			switch string(attribute.Key) {
			case "contract_address":
				contractAddress = HandleAddress(attribute.Value)
			case "gas_consumed":
				gasConsumed, err = HandleGas(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing gas consumed (calculation event): %s", err)
				}
			case "contract_rewards":
				contractRewards, err = HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing contract rewards (calculation event): %s", err)
				}
			case "inflation_rewards":
				inflationRewards, err = HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing inflation rewards (calculation event): %s", err)
				}
			case "metadata":
				metadataCalculationReward, err = HandleMetadata(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while handle metadata (calculation event): %s", err)
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
		rewardCalculationHeight := height - 1

		return db.SaveContractRewardCalculation(
			types.NewContractRewardCalculation(
				contractAddress,
				metadataCalculationReward.RewardAddress,
				metadataCalculationReward.DeveloperAddress,
				gasConsumed,
				contractRewards,
				inflationRewards,
				metadataCalculationReward.CollectPremium,
				metadataCalculationReward.GasRebateToUser,
				metadataCalculationReward.PremiumPercentageCharged,
				eventJson,
				rewardCalculationHeight,
			),
		)
	}

	// Now try to catch reward distribution event
	// Because rewards distriburion happens in another event, we have to collect this data
	// and update db row previously added with 'ContractRewardCalculationEvent'
	if strings.Contains(event.Type, "archway.gastracker.v1.RewardDistributionEvent") {

		var contractDistributionAddress string
		var contractDistributionRewards, leftoverRewards types.GasTrackerReward
		var err error

		// Get all event attributes
		eventAttributes := event.GetAttributes()
		// Handle all the atribution inside the event type
		for _, attribute := range eventAttributes {
			switch string(attribute.Key) {
			case "reward_address":
				contractDistributionAddress = HandleAddress(attribute.Value)
			case "contract_rewards":
				contractDistributionRewards, err = HandleRewards(attribute.Value)
				if err != nil {
					return fmt.Errorf("error while parsing contract rewards (distribution event): %s", err)
				}
			case "leftover_rewards":
				leftoverRewards, err = HandleRewards(attribute.Value)
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
				contractDistributionAddress,
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

func HandleGas(value []byte) (string, error) {
	reg, err := regexp.Compile(`[^0-9]`)
	if err != nil {
		return "", fmt.Errorf("error while converting gas_consumed (calculation event): %s", err)
	}
	num := reg.ReplaceAllString(string(value), "")
	return num, nil
	//return strconv.ParseInt(num, 10, 64)
}

func HandleRewards(value sdk.Coins) (sdk.Coins, error) {
	return value, nil
	//return getGasTrackerRewardFromString(string(value))
}

func HandleMetadata(value []byte) (*types.MetadataReward, error) {
	var metadata types.MetadataReward
	err := json.Unmarshal(value, &metadata)
	fmt.Print(err)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func GetEventJson(event *tmabcitypes.Event) ([]byte, error) {
	// Try to Marshall event
	return event.Marshal()
}
