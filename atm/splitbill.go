package atm

import (
	"math/big"
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	"github.com/shopspring/decimal"

	tb "github.com/tucnak/telebot"
)

// SplitBillInfo contains bill info
type SplitBillInfo struct {
	Title    string
	Payer    string
	Total    decimal.Decimal
	Currency string
	Payee    []string
	EachPay  decimal.Decimal
}

// New create a SplitBillInfo from command
func (info *SplitBillInfo) New(message string) (err error) {
	// Validate input
	input := strings.Fields(message)
	if len(input) < 5 {
		return ErrCommandFormat
	}

	// Fill info
	info.Title = input[1]
	info.Payer = input[2]
	info.Payee = input[4:]
	info.Total, info.Currency, err = ParseMoney(input[3])
	if err != nil {
		return err
	}

	// Validate username
	for i := range info.Payee {
		if info.Payee[i][0] != "@"[0] {
			return ErrUsernameIncorrect
		}
	}

	// Calulate Payment
	divideBy := decimal.NewFromBigInt(big.NewInt(int64(1+len(info.Payee))), 0)
	info.EachPay = info.Total.Div(divideBy)

	return nil
}

// Make a payment record
func (info *SplitBillInfo) Make() (payments []Payment, err error) {
	payments = make([]Payment, 0)
	for i := range info.Payee {
		payment := Payment{
			Title:     info.Title,
			Action:    "Lend",
			Payer:     info.Payer,
			Payee:     info.Payee[i],
			Amount:    info.EachPay,
			Currency:  info.Currency,
			CreatedAt: time.Now(),
		}
		payments = append(payments, payment)
	}
	return
}

// SplitBillHandler Handle a split bill request
func SplitBillHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	info := SplitBillInfo{}
	err := info.New(m.Text)
	if err != nil {
		switch err {
		case ErrEmpty:
			bot.Send(m.Chat, "分單, Format: \n/split 食咩飯 @邊個出錢 3000 @食家A @食家B...")
			//bot.Send(m.Chat, "Split a bill to all people, Format: \n/split Title @payer 3000 @lenderA @lenderB...")
			return
		case ErrCommandFormat:
			bot.Send(m.Chat, "Error: "+err.Error()+"\n/split Title @payer 3000 @lenderA @lenderB...")
			return
		case ErrUsernameIncorrect:
			bot.Send(m.Chat, "Error: "+err.Error()+"\nPlease check the username")
			return
		default:
			bot.Send(m.Chat, "Error: "+err.Error())
		}
		return
	}

	pays, _ := info.Make()
	bot.Send(m.Chat, PrintPayments(pays))

	JapanGroupPOSTPayment(pays)
	return
}
