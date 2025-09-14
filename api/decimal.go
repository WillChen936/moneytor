package api

import (
	"strings"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	decimal, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	d.Decimal = decimal
	return nil
}
