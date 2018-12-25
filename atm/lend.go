package atm

import (
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	"github.com/shopspring/decimal"
	tb "github.com/tucnak/telebot"
)

// LendItem contain a lend record
type LendItem struct {
	Payee    string
	Amount   decimal.Decimal
	Currency string
}

// LendInfo contain multiple lend record from one person
type LendInfo struct {
	Title string
	Payer string
	Lends []LendItem
}

// New create a LendInfo from command
func (info *LendInfo) New(message string) (err error) {
	// Validate input
	input := strings.Fields(message)
	if len(input) == 1 {
		return ErrEmpty
	}
	if len(input) < 5 {
		return ErrCommandFormat
	}
	if len(input)%2 != 1 {
		return ErrCommandFormat
	}

	// Fill info
	info.Title = input[1]
	info.Payer = input[2]

	info.Lends = make([]LendItem, 0)
	lends := input[3:]
	for i := 0; i < len(lends); i += 2 {
		item := LendItem{}
		item.Payee = lends[i]
		item.Amount, item.Currency, err = ParseMoney(lends[i+1])
		if err != nil {
			return err
		}
		info.Lends = append(info.Lends, item)
	}

	// Validate username
	if info.Payer[0] != "@"[0] {
		return ErrUsernameIncorrect
	}
	for i := range info.Lends {
		if info.Lends[i].Payee[0] != "@"[0] {
			return ErrUsernameIncorrect
		}
	}

	return nil
}

// Make a payment record
func (info *LendInfo) Make() (payments []Payment, err error) {
	payments = make([]Payment, 0)
	for i := range info.Lends {
		payment := Payment{
			Title:     info.Title,
			Action:    "Lend",
			Payer:     info.Payer,
			Payee:     info.Lends[i].Payee,
			Amount:    info.Lends[i].Amount,
			Currency:  info.Lends[i].Currency,
			CreatedAt: time.Now(),
		}
		payments = append(payments, payment)
	}
	return
}

// LendHandler Handle a lean request
func LendHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	info := LendInfo{}
	err := info.New(m.Text)
	if err != nil {
		switch err {
		case ErrEmpty:
			bot.Send(m.Chat, "幫人出錢, Format: \n/lend 原因 @邊個出錢 @欠款人A 1500yen @欠款人B 1800yen...")
			//bot.Send(m.Chat, "Lend money to peoples, Format: \n/lend Title @payer @lenderA 1500yen @lenderB 1800yen...")
			return
		case ErrCommandFormat:
			bot.Send(m.Chat, "Error: "+err.Error()+"\nPlease follow format: /lend Title @payer @lenderA 1500yen @lenderB 1800yen...")
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
