package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
delete - admin function
*/

func main() {
	var err error

	// config database
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

	// check tables exist
	_, err = db.Exec("SELECT 1 FROM vote LIMIT 1")
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

	// config bot
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

	b.Handle("\fa1", btnHandler1)
	b.Handle("\fa2", btnHandler2)
	b.Handle("\fa3", btnHandler3)
	b.Handle("\fa4", btnHandler4)
	b.Handle("\fa5", btnHandler5)
	b.Handle("\fa6", btnHandler6)

	b.Handle("/vote", quickVoteHandler)
	b.Handle("/voteadv", advVoteHandler)
	b.Handle("/forwardvote", forwareVoteHandler)
	b.Handle("/delete", removeMessageHandler)

	fmt.Println("Bot starting...")
	b.Start()

}

func handleError(err error) {
	fmt.Println("Error: ", err)
}

func removeMessageHandler(m *tb.Message) {
	if m.Sender.ID != 195152664 {
		b.Send(m.Sender, "You have no permission to do that")
		return
	}
	options := strings.Split(m.Text, " ")
	if len(options) != 3 {
		b.Send(m.Sender, "Please follow format: /delete {ChatID} {MessageID}")
		return
	}

	chatid, err := strconv.Atoi(options[1])
	if err != nil {
		panic(err)
	}

	b.Delete(
		tb.StoredMessage{
			ChatID:    int64(chatid),
			MessageID: options[2],
		})

	var ref VoteRef
	err = db.Select(&ref, "chatid = ? AND messageid = ?", chatid, options[2])
	if err != nil {
		panic(err)
	}
	if ref.VoteID != 0 {
		_, err = db.Exec("DELETE FROM voteref WHERE voteid = ? AND chatid = ? AND messageid = ?", ref.VoteID, ref.ChatID, ref.MessageID)
		if err != nil {
			panic(err)
		}
	}

}
