package types

type GasTrackerReward struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type MetadataReward struct {
	DeveloperAddress         string `json:"developer_address,omitempty"`
	RewardAddress            string `json:"reward_address,omitempty"`
	GasRebateToUser          bool   `json:"gas_rebate_to_user,omitempty"`
	CollectPremium           bool   `json:"collect_premium,omitempty"`
	PremiumPercentageCharged int64  `json:"premium_percentage_charged,string,omitempty"`
}

type ContractRewardCalculation struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string

	GasConsumed string

	ContractRewards  []GasTrackerReward
	InflationRewards []GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged int64

	MetadataJson []byte
	Height       int64
}

// NewContractRewardCalculation allows to easily create a new ContractRewardCalculation
func NewContractRewardCalculation(
	contractAddress string,
	rewardAddress string,
	developerAddress string,
	gasConsumed string,
	contractReward []GasTrackerReward,
	inflationRewards []GasTrackerReward,
	collectPremium bool,
	gasRebateToUser bool,
	premiumPercentageCharged int64,
	metadataJson []byte,
	height int64,
) ContractRewardCalculation {
	return ContractRewardCalculation{
		ContractAddress:          contractAddress,
		RewardAddress:            rewardAddress,
		DeveloperAddress:         developerAddress,
		GasConsumed:              gasConsumed,
		ContractRewards:          contractReward,
		InflationRewards:         inflationRewards,
		CollectPremium:           collectPremium,
		GasRebateToUser:          gasRebateToUser,
		PremiumPercentageCharged: premiumPercentageCharged,
		MetadataJson:             metadataJson,
		Height:                   height,
	}
}

type ContractRewardDistribution struct {
	RewardAddress      string
	DistributedRewards []GasTrackerReward
	LeftoverRewards    []GasTrackerReward
	Height             int64
}

// NewContractRewardDistribution allows to easily create a new ContractRewardDistribution
func NewContractRewardDistribution(
	rewardAddress string,
	distributedRewards []GasTrackerReward,
	leftoverRewards []GasTrackerReward,
	height int64,
) ContractRewardDistribution {
	return ContractRewardDistribution{
		RewardAddress:      rewardAddress,
		DistributedRewards: distributedRewards,
		LeftoverRewards:    leftoverRewards,
		Height:             height,
	}
}

type ContractReward struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string

	GasConsumed int64

	ContractRewards    []GasTrackerReward
	InflationRewards   []GasTrackerReward
	DistributedRewards []GasTrackerReward
	LeftoverRewards    []GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged int64

	MetadataJson []byte

	Height int64
}
