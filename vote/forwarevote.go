package vote

import (
	"strconv"
	"strings"

	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

// ForwareVoteHandler forward a vote message from command
func ForwareVoteHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	options := strings.Fields(m.Text)
	if len(options) != 2 {
		bot.Send(m.Sender, "Please follow format: /forwardvote {VoteID}")
		return
	}
	voteid, err := strconv.Atoi(options[1])
	if err != nil {
		panic(err)
	}

	// Search Vote
	var vote Vote
	err = db.Select(&vote, "id = ? ", voteid)
	if err != nil {
		panic(err)
	}
	if vote.ID == 0 {
		bot.Send(m.Sender, "Vote not found.")
		return
	}

	// Load Choices
	choices := []Choice{}
	err = db.Select(&choices, "voteid = ?", vote.ID)
	if err != nil {
		panic(err)
	}
	vote.Voters = len(choices)

	// Response
	m, err = bot.Send(
		m.Chat,
		vote.String(choices),
		&tb.ReplyMarkup{InlineKeyboard: vote.GenButton()},
	)
	if err != nil {
		panic(err)
	}

	// Add reference to database
	ref := VoteRef{
		VoteID:    vote.ID,
		ChatID:    int(m.Chat.ID),
		MessageID: m.ID,
	}
	err = db.Insert(&ref)
	if err != nil {
		panic(err)
	}
}
