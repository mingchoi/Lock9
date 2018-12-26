package atm

import (
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	"github.com/shopspring/decimal"
	tb "github.com/tucnak/telebot"
)

type ATMInfo struct {
	Title    string
	Payer    string
	Amount   decimal.Decimal
	Currency string
	Payee    string
}

// New create a ATMInfo from command
func (info *ATMInfo) New(message string) (err error) {
	// Validate input
	input := strings.Fields(message)
	if len(input) == 1 {
		return ErrEmpty
	}
	if len(input) != 5 {
		return ErrCommandFormat
	}

	// Fill info
	info.Title = input[1]
	info.Payer = input[2]
	info.Payee = input[4]
	info.Amount, info.Currency, err = ParseMoney(input[3])
	if err != nil {
		return err
	}

	// Validate username
	if info.Payer[0] != "@"[0] || info.Payee[0] != "@"[0] {
		return ErrUsernameIncorrect
	}

	return nil
}

// Make a payment record
func (info *ATMInfo) Make() (payment Payment, err error) {
	payment = Payment{
		Title:     info.Title,
		Action:    "Transfer",
		Payer:     info.Payer,
		Payee:     info.Payee,
		Amount:    info.Amount,
		Currency:  info.Currency,
		CreatedAt: time.Now(),
	}

	return
}

// AtmHandler Handle a transfer request
func AtmHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	info := ATMInfo{}
	err := info.New(m.Text)
	if err != nil {
		switch err {
		case ErrEmpty:
			bot.Send(m.Chat, "過數比人, Format: \n/atm 原因 @邊個比錢 1500 @邊個收錢")
			//bot.Send(m.Chat, "Transfer money to someone, Format: \n/atm Title @payer 1500yen @payee")
			return
		case ErrCommandFormat:
			bot.Send(m.Chat, "訓撚醒未: "+err.Error()+"\n睇清楚格式啦: \n/atm Title @payer 1500yen @payee")
			//bot.Send(m.Chat, "Error: "+err.Error()+"\nPlease follow format: \n/atm Title @payer 1500yen @payee")
			return
		case ErrUsernameIncorrect:
			bot.Send(m.Chat, "訓啦柒頭: "+err.Error()+"\n你tag緊邊條柒頭？")
			//bot.Send(m.Chat, "Error: "+err.Error()+"\nPlease check the username")
			return
		default:
			bot.Send(m.Chat, "Error: "+err.Error())
			return
		}
	}

	pay, _ := info.Make()
	bot.Send(m.Chat, PrintPayments([]Payment{pay}))

	JapanGroupPOSTPayment([]Payment{pay})

	return
}
