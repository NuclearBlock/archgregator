package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nuclearblock/archgregator/types"
)

// sumGasTxs returns the total gas consumed by a set of transactions.
func sumGasTxs(txs []*types.Tx) uint64 {
	var totalGas uint64
	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed)
	}
	return totalGas
}

// This is a Cosmologger solution to get correct rewards data
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

	//numValue, err := strconv.ParseFloat(tmpMap["amount"].(string), 64)
	// if err != nil {
	// 	return types.GasTrackerReward{}, err
	// }

	return types.GasTrackerReward{
		Denom:  tmpMap["denom"].(string),
		Amount: tmpMap["amount"].(string),
	}, nil
}
