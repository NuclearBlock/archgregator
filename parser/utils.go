package parser

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/nuclearblock/archgregator/types"
)

// findValidatorByAddr finds a validator by a consensus address given a set of
// Tendermint validators for a particular block. If no validator is found, nil
// is returned.
func findValidatorByAddr(consAddr string, vals *tmctypes.ResultValidators) *tmtypes.Validator {
	for _, val := range vals.Validators {
		if consAddr == sdk.ConsAddress(val.Address).String() {
			return val
		}
	}

	return nil
}

// sumGasTxs returns the total gas consumed by a set of transactions.
func sumGasTxs(txs []*types.Tx) uint64 {
	var totalGas uint64

	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed)
	}

	return totalGas
}

// // GetContractInfo implements wasmsource.Source
// func (s Source) GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error) {
// 	ctx, err := s.LoadHeight(height)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while loading height: %s", err)
// 	}

// 	res, err := s.q.ContractInfo(
// 		sdk.WrapSDKContext(ctx),
// 		&wasmtypes.QueryContractInfoRequest{
// 			Address: contractAddr,
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while getting contract info: %s", err)
// 	}

// 	return res, nil
// }

// // GetContractInfo implements wasmsource.Source
// func (s Source) GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error) {
// 	res, err := s.wasmClient.ContractInfo(
// 		remote.GetHeightRequestContext(s.Ctx, height),
// 		&wasmtypes.QueryContractInfoRequest{
// 			Address: contractAddr,
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("error while getting contract info: %s", err)
// 	}

// 	return res, nil
// }
