package fees

import (
	"github.com/Oneledger/protocol/data/balance"
	"math/big"
)

type FeeOption struct {
	FeeCurrency   balance.Currency `json:"feeCurrency"`
	MinFeeDecimal int64            `json:"minFeeDecimal"`

	minimalFee *balance.Coin
}

func (fo *FeeOption) MinFee() *balance.Coin {
	if fo.minimalFee == nil {
		amount := balance.Amount{Int: *big.NewInt(0).Exp(big.NewInt(10), big.NewInt(fo.FeeCurrency.Decimal-fo.MinFeeDecimal), nil)}
		coin := fo.FeeCurrency.NewCoinFromAmount(amount)
		fo.minimalFee = &coin
	}
	return fo.minimalFee
}
