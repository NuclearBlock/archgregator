package parser

import (
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
