package types

type GasTrackerReward struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type MetadataReward struct {
	DeveloperAddress         string `json:"developer_address"`
	RewardAddress            string `json:"reward_address"`
	GasRebateToUser          bool   `json:"gas_rebate_to_user"`
	CollectPremium           bool   `json:"collect_premium"`
	PremiumPercentageCharged int64  `json:"premium_percentage_charged"`
}

type ContractRewardCalculation struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string

	GasConsumed uint64

	ContractRewards  GasTrackerReward
	InflationRewards GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged uint64

	DataCalculationJson string
	Height              uint64
}

// NewContractRewardCalculation allows to easily create a new ContractRewardCalculation
func NewContractRewardCalculation(
	contractAddress string,
	rewardAddress string,
	developerAddress string,
	gasConsumed uint64,
	contractReward GasTrackerReward,
	inflationRewards GasTrackerReward,
	collectPremium bool,
	gasRebateToUser bool,
	premiumPercentageCharged uint64,
	dataCalculationJson string,
	height uint64,
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
	DataDistributionJson string
	Height               uint64
}

// NewContractRewardDistribution allows to easily create a new ContractRewardDistribution
func NewContractRewardDistribution(
	contractAddress string,
	contractReward GasTrackerReward,
	leftoverRewards GasTrackerReward,
	dataDistributionJson string,
	height uint64,
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

	GasConsumed uint64

	ContractRewards  GasTrackerReward
	InflationRewards GasTrackerReward
	LeftoverRewards  GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged uint64

	DataCalculationJson  string
	DataDistributionJson string

	Height uint64
}
