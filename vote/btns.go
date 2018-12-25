package vote

import (
	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

// BtnHandler1 handle vote from user
func BtnHandler1(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		1,
		bot,
		db,
	)
}

// BtnHandler2 handle vote from user
func BtnHandler2(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		2,
		bot,
		db,
	)
}

// BtnHandler3 handle vote from user
func BtnHandler3(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		3,
		bot,
		db,
	)
}

// BtnHandler4 handle vote from user
func BtnHandler4(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		4,
		bot,
		db,
	)
}

// BtnHandler5 handle vote from user
func BtnHandler5(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		5,
		bot,
		db,
	)
}

// BtnHandler6 handle vote from user
func BtnHandler6(c *tb.Callback, bot *tb.Bot, db *s2s.DB) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		6,
		bot,
		db,
	)
}
