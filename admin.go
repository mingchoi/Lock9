package main

import (
	"strconv"
	"strings"

	tb "github.com/tucnak/telebot"
)

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
