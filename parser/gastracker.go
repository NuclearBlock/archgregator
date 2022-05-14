package parser

import (
	"fmt"
	"time"

	gastrackertypes "github.com/archway-network/archway/x/gastracker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	database "github.com/nuclearblock/archgregator/database"
	types "github.com/nuclearblock/archgregator/types"
	tmabcitypes "github.com/tendermint/tendermint/abci/types"
)

// HandleMsgSetMetadata allows to properly handle a Gastracker MsgSetMetadata
func HandleMsgSetMetadata(index int, tx *types.Tx, msg *gastrackertypes.MsgSetContractMetadata, db database.Database) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return db.SaveGasTrackerContractMetadata(
		types.NewGasTrackerContractMetadata(msg, tx.TxHash, timestamp, tx.Height),
	)
}

// HandleGasTrackerRewards allows to build a new smart contract reward instance from gastracker event
func HandleGasTrackerRewards(event *tmabcitypes.Event, height int64, db database.Database) error {

	// Try to parse acbi event
	typedEvent, err := sdk.ParseTypedEvent(*event)
	if err != nil {
		return fmt.Errorf("error while parsing typed event to proto.message: %s", err)
	}

	// We have to check if this Revard Calculation event or Reward Distribution event
	switch gastrackerEvent := typedEvent.(type) {
	case *gastrackertypes.ContractRewardCalculationEvent:
		// We have to decrement target block height,
		// because reward is always processed in the next BeginBlock
		rewardHeight := height - 1

		return db.SaveContractRewardCalculation(
			types.NewContractRewardCalculation(
				gastrackerEvent.ContractAddress,
				gastrackerEvent.GasConsumed,
				gastrackerEvent.ContractRewards,
				gastrackerEvent.InflationRewards,
				rewardHeight,
			),
		)
	case *gastrackertypes.RewardDistributionEvent:
		// Now try to catch reward distribution event
		// Because rewards distriburion happens in another event, we have to collect this data
		// and update db row previously added with 'ContractRewardCalculationEvent'

		// We have to decrement target block height,
		// because reward is always processed in the next BeginBlock
		// We need this fied to identify correct 'calculation' table row
		distributionHeight := height - 1

		return db.SaveContractRewardDistribution(
			types.NewContractRewardDistribution(
				gastrackerEvent.RewardAddress,
				gastrackerEvent.ContractRewards,
				gastrackerEvent.LeftoverRewards,
				distributionHeight,
			),
		)

	}

	return nil
}
