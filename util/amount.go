package util

import (
	"errors"
	"math/big"
)

// CGAAmount wraps a raw amount.
type CGAAmount struct {
	Raw *big.Int
}

func (CGAAmount) exp() *big.Int {
	x := big.NewInt(10)
	// Change RAW value
	return x.Exp(x, big.NewInt(29), nil)
}

// CGAAmountFromString parses CGA amounts in strings.
func CGAAmountFromString(s string) (n CGAAmount, err error) {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		err = errors.New("unable to parse CGA amount")
		return
	}
	r = r.Mul(r, new(big.Rat).SetInt(n.exp()))
	if !r.IsInt() {
		err = errors.New("unable to parse CGA amount")
		return
	}
	n.Raw = r.Num()
	return
}

func (n CGAAmount) String() string {
	r := new(big.Rat).SetFrac(n.Raw, n.exp())
	s := r.FloatString(29)
	return s[:len(s)-24]
}
