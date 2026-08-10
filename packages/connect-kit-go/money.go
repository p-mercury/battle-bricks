package connectkit

import (
	"math/big"

	"google.golang.org/genproto/googleapis/type/money"
)

func MultiplyMoney(a *money.Money, factor int64) *money.Money {
	if a == nil {
		return nil
	}

	billion := big.NewInt(1e9)
	totalNanos := new(big.Int).Add(
		new(big.Int).Mul(big.NewInt(a.Units), billion),
		big.NewInt(int64(a.Nanos)),
	)
	totalNanos.Mul(totalNanos, big.NewInt(factor))

	units := new(big.Int).Quo(totalNanos, billion)
	nanos := new(big.Int).Rem(totalNanos, billion)

	return &money.Money{
		CurrencyCode: a.CurrencyCode,
		Units:        units.Int64(),
		Nanos:        int32(nanos.Int64()),
	}
}

func AddMoney(a, b *money.Money) *money.Money {
	if a == nil || b == nil {
		return nil
	}

	billion := big.NewInt(1e9)

	aNanos := new(big.Int).Add(
		new(big.Int).Mul(big.NewInt(a.Units), billion),
		big.NewInt(int64(a.Nanos)),
	)
	bNanos := new(big.Int).Add(
		new(big.Int).Mul(big.NewInt(b.Units), billion),
		big.NewInt(int64(b.Nanos)),
	)

	totalNanos := new(big.Int).Add(aNanos, bNanos)

	units := new(big.Int).Quo(totalNanos, billion)
	nanos := new(big.Int).Rem(totalNanos, billion)

	return &money.Money{
		CurrencyCode: a.CurrencyCode,
		Units:        units.Int64(),
		Nanos:        int32(nanos.Int64()),
	}
}

func SubtractMoney(a, b *money.Money) *money.Money {
	if a == nil || b == nil {
		return nil
	}

	billion := big.NewInt(1e9)

	aNanos := new(big.Int).Add(
		new(big.Int).Mul(big.NewInt(a.Units), billion),
		big.NewInt(int64(a.Nanos)),
	)
	bNanos := new(big.Int).Add(
		new(big.Int).Mul(big.NewInt(b.Units), billion),
		big.NewInt(int64(b.Nanos)),
	)

	totalNanos := new(big.Int).Sub(aNanos, bNanos)

	units := new(big.Int).Quo(totalNanos, billion)
	nanos := new(big.Int).Rem(totalNanos, billion)

	return &money.Money{
		CurrencyCode: a.CurrencyCode,
		Units:        units.Int64(),
		Nanos:        int32(nanos.Int64()),
	}
}
