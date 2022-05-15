package types

import (
	"time"

	gastrackertypes "github.com/archway-network/archway/x/gastracker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GasTrackerContractMetadata struct {
	Sender          string
	ContractAddress string
	Metadata        gastrackertypes.ContractInstanceMetadata
	MetadataJson    []byte
	TxHash          string
	SavedAt         time.Time
	Height          int64
}

func NewGasTrackerContractMetadata(
	msg *gastrackertypes.MsgSetContractMetadata,
	txHash string,
	savedAt time.Time,
	height int64,
) GasTrackerContractMetadata {
	metadataJson := []byte("{}")

	return GasTrackerContractMetadata{
		Sender:          msg.Sender,
		ContractAddress: msg.ContractAddress,
		Metadata:        *msg.Metadata,
		MetadataJson:    metadataJson,
		TxHash:          txHash,
		SavedAt:         savedAt,
		Height:          height,
	}
}

type ContractRewardCalculation struct {
	ContractAddress  string
	GasConsumed      uint64
	ContractRewards  sdk.DecCoins
	InflationRewards sdk.DecCoin
	Height           int64
}

// NewContractRewardCalculation allows to easily create a new ContractRewardCalculation
func NewContractRewardCalculation(
	contractAddress string,
	gasConsumed uint64,
	contractReward *sdk.DecCoins,
	inflationRewards *sdk.DecCoin,
	height int64,
) ContractRewardCalculation {
	return ContractRewardCalculation{
		ContractAddress:  contractAddress,
		GasConsumed:      gasConsumed,
		ContractRewards:  *contractReward,
		InflationRewards: *inflationRewards,
		Height:           height,
	}
}

type ContractRewardDistribution struct {
	RewardAddress      string
	DistributedRewards sdk.Coins
	LeftoverRewards    sdk.DecCoin
	Height             int64
}

// NewContractRewardDistribution allows to easily create a new ContractRewardDistribution
func NewContractRewardDistribution(
	rewardAddress string,
	distributedRewards *sdk.Coins,
	leftoverRewards *sdk.DecCoin,
	height int64,
) ContractRewardDistribution {
	return ContractRewardDistribution{
		RewardAddress:      rewardAddress,
		DistributedRewards: *distributedRewards,
		LeftoverRewards:    *leftoverRewards,
		Height:             height,
	}
}

type ContractReward struct {
	ContractAddress    string
	GasConsumed        int64
	ContractRewards    []*sdk.DecCoin
	InflationRewards   *sdk.DecCoin
	DistributedRewards []*sdk.Coin
	LeftoverRewards    []*sdk.DecCoin
	Height             int64
}
