package types

type GasTrackerReward struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
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

	GasConsumed int64

	ContractRewards  GasTrackerReward
	InflationRewards GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged int64

	DataCalculationJson []byte
	Height              int64
}

// NewContractRewardCalculation allows to easily create a new ContractRewardCalculation
func NewContractRewardCalculation(
	contractAddress string,
	rewardAddress string,
	developerAddress string,
	gasConsumed int64,
	contractReward GasTrackerReward,
	inflationRewards GasTrackerReward,
	collectPremium bool,
	gasRebateToUser bool,
	premiumPercentageCharged int64,
	dataCalculationJson []byte,
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
		DataCalculationJson:      dataCalculationJson,
		Height:                   height,
	}
}

type ContractRewardDistribution struct {
	ContractAddress      string
	ContractRewards      GasTrackerReward
	LeftoverRewards      GasTrackerReward
	DataDistributionJson []byte
	Height               int64
}

// NewContractRewardDistribution allows to easily create a new ContractRewardDistribution
func NewContractRewardDistribution(
	contractAddress string,
	contractReward GasTrackerReward,
	leftoverRewards GasTrackerReward,
	dataDistributionJson []byte,
	height int64,
) ContractRewardDistribution {
	return ContractRewardDistribution{
		ContractAddress:      contractAddress,
		ContractRewards:      contractReward,
		LeftoverRewards:      leftoverRewards,
		DataDistributionJson: dataDistributionJson,
		Height:               height,
	}
}

type ContractReward struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string

	GasConsumed int64

	ContractRewards  GasTrackerReward
	InflationRewards GasTrackerReward
	LeftoverRewards  GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged int64

	DataCalculationJson  []byte
	DataDistributionJson []byte

	Height int64
}
