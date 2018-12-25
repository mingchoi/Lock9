package atm

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Payment holds a transcation record
type Payment struct {
	ID        int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	Title     string
	Action    string
	Payer     string
	Payee     string
	Amount    decimal.Decimal
	Currency  string
	CreatedAt time.Time
}

// PrintPayments print payments to string
func PrintPayments(payments []Payment) string {
	str := payments[0].Title + "\n"
	for _, p := range payments {
		switch p.Action {
		case "Transfer":
			str += fmt.Sprintf("%s 過數 %s %s 比 %s\n", p.Payer, p.Amount.StringFixed(0), p.Currency, p.Payee)
		case "Lend":
			str += fmt.Sprintf("%s 借左 %s %s 比 %s\n", p.Payer, p.Amount.StringFixed(0), p.Currency, p.Payee)
		}
	}
	return str
}
