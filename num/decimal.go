package num

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrDecScan = errors.New("expected bytes when scanning decimal field")

	noopDec = Dec{}

	nullDecStr = "null"
	nullDec    = []byte(nullDecStr)

	Zero       = NewDecFromInt(0)
	One        = NewDecFromInt(1)
	MinusOne   = NewDecFromInt(-1)
	Nil        = NewDecNil()
	OneHundred = NewDecFromInt(100)
)

// Dec represents atani broker decimal type. Take into account de/serialization could be optimized by employing
// the unsafe package when transforming bytes to string and vice-versa.
type Dec struct {
	dec decimal.Decimal

	isset bool
}

func (d *Dec) scanBytes(data []byte) (err error) {
	// data (bytes from db) will have greater capacity than needed. Better not to keep that underlying array dangling
	raw := make([]byte, len(data))
	copy(raw, data)
	d.dec, err = decimal.NewFromString(string(raw))
	return
}

func (d *Dec) Scan(x any) (err error) {
	if x == nil {
		return
	}

	switch val := x.(type) {
	// Handle bytes
	case []uint8:
		err = d.scanBytes(val)
	// Handle string
	case string:
		d.dec, err = decimal.NewFromString(val)
	case float64:
		d.dec = decimal.NewFromFloat(val)
	case int64:
		d.dec = decimal.NewFromInt(val)
	default:
		err = ErrDecScan
	}

	if err == nil {
		d.isset = true
	}

	return
}

func (d Dec) MarshalJSON() ([]byte, error) {
	if !d.isset {
		return nullDec, nil
	}
	str := "\"" + d.dec.String() + "\""
	return []byte(str), nil
}

func (d *Dec) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullDec) {
		return nil
	}

	if err := d.dec.UnmarshalJSON(data); err != nil {
		return err
	}

	d.isset = true
	return nil
}

func (d Dec) Equal(other Dec) bool {
	if d.isset != other.Isset() {
		return false
	}
	return d.dec.Equal(other.dec)
}

func (d Dec) Equals(other Dec) bool {
	return d.Equal(other)
}

func (d Dec) NotEquals(other Dec) bool {
	return !d.Equal(other)
}

func (d Dec) GreaterThan(other Dec) bool {
	return d.dec.GreaterThan(other.dec)
}

func (d Dec) GreaterThanOrEqual(other Dec) bool {
	return d.dec.GreaterThanOrEqual(other.dec)
}

func (d Dec) LessThan(other Dec) bool {
	return d.dec.LessThan(other.dec)
}

func (d Dec) LessThanOrEqual(other Dec) bool {
	return d.dec.LessThanOrEqual(other.dec)
}

func (d Dec) Sub(other Dec) Dec {
	return Dec{dec: d.dec.Sub(other.dec), isset: d.isset}
}

func (d Dec) Add(other Dec) Dec {
	return Dec{dec: d.dec.Add(other.dec), isset: d.isset}
}

func (d Dec) AddOr(other Dec) Dec {
	if d.Isset() {
		return d.Add(other)
	}
	return other
}

func (d Dec) Mul(other Dec) Dec {
	return Dec{dec: d.dec.Mul(other.dec), isset: d.isset}
}

func (d Dec) Mod(other Dec) Dec {
	return Dec{dec: d.dec.Mod(other.dec), isset: d.isset}
}

func (d Dec) Div(other Dec) Dec {
	return Dec{dec: d.dec.Div(other.dec), isset: d.isset}
}

// Percent assumes that current amount represents a percentage from 0 to 1. Returns percentage value.
func (d Dec) Percent() Dec {
	return d.Div(OneHundred)
}

func (d Dec) AddPercent(p Dec) Dec {
	return d.Add(d.ApplyPercent(p))
}

func (d Dec) SubPercent(p Dec) Dec {
	return d.Sub(d.ApplyPercent(p))
}

func (d Dec) ApplyPercent(p Dec) Dec {
	return d.Mul(p.Percent())
}

func (d Dec) Inverse() Dec {
	return One.Div(d)
}

func (d Dec) Exponent() int32 {
	return d.dec.Exponent()
}

func (d Dec) NumberOfDecimals() int32 {
	f, _ := d.Float64()
	df := NewDecFromFloat(f)

	return Abs[int32](df.Exponent())
}

func (d Dec) Clone() Dec {
	return Dec{dec: d.dec.Copy(), isset: d.isset}
}

func (d Dec) Neg() Dec {
	return Dec{dec: d.dec.Neg(), isset: d.isset}
}

func (d Dec) MustNeg() Dec {
	if d.dec.IsNegative() {
		return Dec{dec: d.dec, isset: d.isset}
	}
	return Dec{dec: d.dec.Neg(), isset: d.isset}
}

func (d Dec) MustPos() Dec {
	if d.dec.IsNegative() {
		return Dec{dec: d.dec.Neg(), isset: d.isset}
	}
	return Dec{dec: d.dec, isset: d.isset}
}

func (d Dec) IsNil() bool {
	return !d.isset
}

func (d Dec) Isset() bool {
	return d.isset
}

func (d Dec) Abs() Dec {
	return Dec{dec: d.dec.Abs(), isset: d.isset}
}

func (d Dec) Floor() Dec {
	return Dec{dec: d.dec.Floor(), isset: d.isset}
}

func (d Dec) Float64() (float64, bool) {
	if !d.isset {
		return float64(0), true
	}

	return d.dec.Float64()
}

func (d Dec) String() string {
	if !d.isset {
		return nullDecStr
	}

	return d.dec.String()
}

func (d Dec) LatinString() string {
	return strings.ReplaceAll(d.String(), ".", ",")
}

func (d Dec) Value() (driver.Value, error) {
	if !d.isset {
		return nil, nil
	}

	return d.dec.String(), nil
}

func (d Dec) ToPostgres() any {
	if !d.isset {
		return nil
	}

	return d.dec.String()
}

// IsPositive return
//
//	true if d > 0
//	false if d == 0
//	false if d < 0
func (d Dec) IsPositive() bool {
	if !d.isset {
		return false
	}

	return d.dec.IsPositive()
}

func (d Dec) IsPositiveOrZero() bool {
	if !d.isset {
		return false
	}

	return !d.dec.IsNegative()
}

func (d Dec) IsNegativeOrZero() bool {
	if !d.isset {
		return false
	}

	return !d.dec.IsPositive()
}

func (d Dec) IntPart() int64 {
	if !d.isset {
		return int64(0)
	}

	return d.dec.IntPart()
}

// IsNegative return
//
//	true if d < 0
//	false if d == 0
//	false if d > 0
func (d Dec) IsNegative() bool {
	if !d.isset {
		return false
	}

	return d.dec.IsNegative()
}

func (d Dec) IsZero() bool {
	if !d.isset {
		return false
	}

	return d.dec.IsZero()
}

func (d Dec) IfNilOrZeroThen(other Dec) Dec {
	return d.Map(func(dec Dec) Dec {
		if !dec.Isset() || dec.IsZero() {
			return other
		}

		return dec
	})

}

func (d Dec) Map(fn func(Dec) Dec) Dec {
	return fn(d)
}

func (d Dec) IsLowerThanZero() bool {
	if !d.isset {
		return false
	}

	return d.dec.LessThan(decimal.Zero)
}

func (d Dec) IsGreaterThanZero() bool {
	if !d.isset {
		return false
	}

	return d.dec.GreaterThan(decimal.Zero)
}

func (d Dec) IsLowerThanOrEqualsZero() bool {
	return !d.IsGreaterThanZero()
}

func (d Dec) Truncate(precision int32) Dec {
	return Dec{dec: d.dec.Truncate(precision), isset: d.isset}
}

func (d Dec) OrElse(other Dec) Dec {
	if d.Isset() {
		return d
	}

	return other
}

func (d Dec) OrZero() Dec {
	return d.OrElse(Zero)
}

func (d Dec) Match(
	onNil func() Dec,
	onZero func() Dec,
	onValue func() Dec,
) Dec {
	switch {
	case !d.isset:
		return onNil()
	case d.dec.IsZero():
		return onZero()
	default:
		return onValue()
	}
}

func (d Dec) Round() Dec {
	return Dec{dec: d.dec.Round(12), isset: d.isset}
}

func (d Dec) RoundTo(places int32) Dec {
	return Dec{dec: d.dec.Round(places), isset: d.isset}
}

func MustDecFromString(str string) Dec {
	dec, err := decimal.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return Dec{dec: dec, isset: true}
}

func NewDecFromString(str string) (Dec, error) {
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return noopDec, err
	}
	return Dec{dec: dec, isset: true}, nil
}

func NewDecFromAny(num any) (Dec, error) {
	switch val := num.(type) {
	case int:
		return NewDecFromInt(int64(val)), nil
	case int64:
		return NewDecFromInt(val), nil
	case float64:
		return NewDecFromFloat(val), nil
	case string:
		return MustDecFromString(val), nil
	default:
		return noopDec, errors.New("invalid type. The valid types are int, int64, float64 and string")
	}
}

func MustDecFromAny(num any) Dec {
	dec, err := NewDecFromAny(num)
	if err != nil {
		panic("MustDecFromAny: " + err.Error())
	}

	return dec
}

func NewDecFromInt(n int64) Dec {
	return Dec{dec: decimal.NewFromInt(n), isset: true}
}

func NewDecFromFloat(n float64) Dec {
	return Dec{dec: decimal.NewFromFloat(n), isset: true}
}

func NewDecNil() Dec {
	return Dec{isset: false}
}

func DecSum(decs []Dec) Dec {
	res := Zero.Clone()

	for _, dec := range decs {
		res = res.Add(dec)
	}

	return res
}
