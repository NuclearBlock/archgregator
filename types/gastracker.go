package types

import (
	"time"

	gastrackertypes "github.com/archway-network/archway/x/gastracker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GasTrackerContractMetadata represents the Gastracker contract Metadata
type GasTrackerContractMetadata struct {
	Sender          string
	ContractAddress string
	Metadata        gastrackertypes.ContractInstanceMetadata
	MetadataJson    []byte
	TxHash          string
	SavedAt         time.Time
	Height          int64
}

// NewGasTrackerContractMetadata allows to easily create a new GasTrackerContractMetadata
func NewGasTrackerContractMetadata(
	msg *gastrackertypes.MsgSetContractMetadata,
	txHash string,
	savedAt time.Time,
	height int64,
) GasTrackerContractMetadata {
	metadataJson := []byte("{}") // TO-DO

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

// ContractRewardCalculation represents the Gastracker reward calculation data
type ContractRewardCalculation struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string
	GasConsumed      uint64
	ContractRewards  sdk.DecCoins
	InflationRewards sdk.DecCoins
	Height           int64
}

// NewContractRewardCalculation allows to easily create a new ContractRewardCalculation
func NewContractRewardCalculation(
	contractAddress string,
	rewardAddress string,
	developerAddress string,
	gasConsumed uint64,
	contractRewards []*sdk.DecCoin,
	inflationRewards *sdk.DecCoin,
	height int64,
) ContractRewardCalculation {
	return ContractRewardCalculation{
		ContractAddress:  contractAddress,
		RewardAddress:    rewardAddress,
		DeveloperAddress: developerAddress,
		GasConsumed:      gasConsumed,
		ContractRewards:  sdk.NewDecCoins(*contractRewards[0]),
		InflationRewards: sdk.NewDecCoins(*inflationRewards),
		Height:           height,
	}
}

// ContractRewardDistribution represents the Gastracker reward distribution data
type ContractRewardDistribution struct {
	RewardAddress      string
	DistributedRewards sdk.Coins
	LeftoverRewards    sdk.DecCoins
	Height             int64
}

// NewContractRewardDistribution allows to easily create a new ContractRewardDistribution
func NewContractRewardDistribution(
	rewardAddress string,
	distributedRewards []*sdk.Coin,
	leftoverRewards []*sdk.DecCoin,
	height int64,
) ContractRewardDistribution {
	return ContractRewardDistribution{
		RewardAddress:      rewardAddress,
		DistributedRewards: sdk.NewCoins(*distributedRewards[0]),
		LeftoverRewards:    sdk.NewDecCoins(*leftoverRewards[0]),
		Height:             height,
	}
}
