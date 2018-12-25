package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mingchoi/lock9/atm"
	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

var b *tb.Bot
var db *s2s.DB

/*
BOT SETTING:
vote - Start a quick vote by: /vote Topic OptionA OptionB
voteadv - Start a vote by: /voteadv {single|multiple} Topic OptionA OptionB
forwardvote - Forward a vote by: /forwardvote {VoteID}

atm - Transfer money to someone by: /atm Title @payer 1500yen @payee
lend - Lend money to peoples by: /lend Title @payer @lenderA 1500yen @lenderB 1800yen...
split - Split a bill to all people by: /split Title @payer 3000 @lenderA @lenderB...

delete - admin function
*/

func main() {
	var err error

	// Config database
	db, err = s2s.Open(
		"mysql",
		os.Getenv("LOCK9_DB_SECRET")+
			"@tcp("+
			strings.Replace(os.Getenv("DB_PORT"), "tcp://", "", 1)+
			")/"+
			os.Getenv("LOCK9_DB_NAME")+
			"?charset=utf8&parseTime=True&loc=Local",
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	checkDBTable()

	// Config bot
	b, err = tb.NewBot(tb.Settings{
		Token:    os.Getenv("LOCK9_API_SECRET"),
		URL:      "https://api.telegram.org",
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: handleError,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle Vote command
	b.Handle("/vote", quickVoteHandler)
	b.Handle("/voteadv", advVoteHandler)
	b.Handle("/forwardvote", forwareVoteHandler)

	// Handle ATM command
	b.Handle("/atm", func(m *tb.Message) { atm.AtmHandler(m, b, db) })
	b.Handle("/lend", func(m *tb.Message) { atm.LendHandler(m, b, db) })
	b.Handle("/split", func(m *tb.Message) { atm.SplitBillHandler(m, b, db) })

	// Handle other command
	b.Handle("/delete", removeMessageHandler)

	// Handle button callback
	b.Handle("\fa1", btnHandler1)
	b.Handle("\fa2", btnHandler2)
	b.Handle("\fa3", btnHandler3)
	b.Handle("\fa4", btnHandler4)
	b.Handle("\fa5", btnHandler5)
	b.Handle("\fa6", btnHandler6)

	// Start bot
	fmt.Println("Bot started")
	b.Start()
}

func handleError(err error) {
	fmt.Println("Error: ", err)
	debug.PrintStack()
}

func checkDBTable() {
	// check if tables exist
	_, err := db.Exec("SELECT 1 FROM vote LIMIT 1")
	if err != nil {
		err = db.CreateTable(&Vote{})
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	_, err = db.Exec("SELECT 1 FROM choice LIMIT 1")
	if err != nil {
		err = db.CreateTable(&Choice{})
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	_, err = db.Exec("SELECT 1 FROM voteref LIMIT 1")
	if err != nil {
		err = db.CreateTable(&VoteRef{})
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
