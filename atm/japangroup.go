package atm

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/shopspring/decimal"
)

// JapanGroupAccountForm is the model of google form
type JapanGroupAccountForm struct {
	Time    time.Time
	Title   string
	Marco   decimal.Decimal
	Kariko  decimal.Decimal
	Kitchen decimal.Decimal
	Ming    decimal.Decimal
	Rory    decimal.Decimal
	Howard  decimal.Decimal
	KC      decimal.Decimal
	Kevee   decimal.Decimal
	Melody  decimal.Decimal
	Steven  decimal.Decimal
	Chung   decimal.Decimal
}

// JapanGroupPOSTPayment post data to google form
func JapanGroupPOSTPayment(payments []Payment) {
	form := JapanGroupAccountForm{
		Time:  payments[0].CreatedAt,
		Title: payments[0].Title,
	}
	// Calculate
	for _, p := range payments {
		switch p.Payee {
		case "@husky":
			form.Marco = form.Marco.Add(p.Amount)
		case "@Kariko23":
			form.Kariko = form.Kariko.Add(p.Amount)
		case "@kitchen4848":
			form.Kitchen = form.Kitchen.Add(p.Amount)
		case "@winampmaker":
			form.Ming = form.Ming.Add(p.Amount)
		case "@roriiiiii":
			form.Rory = form.Rory.Add(p.Amount)
		case "@Eternal_Ha":
			form.Howard = form.Howard.Add(p.Amount)
		case "@toki_usagi":
			form.KC = form.KC.Add(p.Amount)
		case "@KillerKaster":
			form.Kevee = form.Kevee.Add(p.Amount)
		case "@misakikasim":
			form.Melody = form.Melody.Add(p.Amount)
		case "@azumiwaki":
			form.Steven = form.Steven.Add(p.Amount)
		case "@chungnnn":
			form.Chung = form.Chung.Add(p.Amount)
		}

		switch p.Payer {
		case "@husky":
			form.Marco = form.Marco.Add(p.Amount.Neg())
		case "@Kariko23":
			form.Kariko = form.Kariko.Add(p.Amount.Neg())
		case "@kitchen4848":
			form.Kitchen = form.Kitchen.Add(p.Amount.Neg())
		case "@winampmaker":
			form.Ming = form.Ming.Add(p.Amount.Neg())
		case "@roriiiiii":
			form.Rory = form.Rory.Add(p.Amount.Neg())
		case "@Eternal_Ha":
			form.Howard = form.Howard.Add(p.Amount.Neg())
		case "@toki_usagi":
			form.KC = form.KC.Add(p.Amount.Neg())
		case "@KillerKaster":
			form.Kevee = form.Kevee.Add(p.Amount.Neg())
		case "@misakikasim":
			form.Melody = form.Melody.Add(p.Amount.Neg())
		case "@azumiwaki":
			form.Steven = form.Steven.Add(p.Amount.Neg())
		case "@chungnnn":
			form.Chung = form.Chung.Add(p.Amount.Neg())
		}
	}

	// Post Form
	link := os.Getenv("LOCK9_JAPANGROUP_FORMURL")
	d := url.Values{}
	d.Add("entry.504739006", form.Title)
	d.Add("entry.870492299", DecimalToStringOrEmpty(form.Marco))
	d.Add("entry.118891274", DecimalToStringOrEmpty(form.Kariko))
	d.Add("entry.97956543", DecimalToStringOrEmpty(form.Kitchen))
	d.Add("entry.1643391874", DecimalToStringOrEmpty(form.Ming))
	d.Add("entry.1977093094", DecimalToStringOrEmpty(form.Rory))
	d.Add("entry.1460219666", DecimalToStringOrEmpty(form.Howard))
	d.Add("entry.275984198", DecimalToStringOrEmpty(form.KC))
	d.Add("entry.1819777074", DecimalToStringOrEmpty(form.Kevee))
	d.Add("entry.548985734", DecimalToStringOrEmpty(form.Melody))
	d.Add("entry.2035442548", DecimalToStringOrEmpty(form.Steven))
	d.Add("entry.943940356", DecimalToStringOrEmpty(form.Chung))
	_, err := http.PostForm(link, d)
	if err != nil {
		panic(err)
	}
}

// DecimalToStringOrEmpty convert decimal to string, empty if value is 0
func DecimalToStringOrEmpty(d decimal.Decimal) string {
	if d == (decimal.Decimal{}) {
		return ""
	}
	return d.StringFixed(0)
}
