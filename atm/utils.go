package atm

import (
	"strings"

	"github.com/shopspring/decimal"
)

// ParseMoney parse string to amount and currency
func ParseMoney(input string) (total decimal.Decimal, currency string, err error) {
	currency = "JPY"
	total, err = decimal.NewFromString(input)
	if err != nil {
		total, err = decimal.NewFromString(input[:len(input)-3])
		currency = input[len(input)-3:]
		if err != nil {
			return
		}
	}

	currency = strings.ToUpper(currency)
	if currency == "YEN" {
		currency = "JPY"
	}

	return
}
