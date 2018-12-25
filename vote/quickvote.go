package vote

import (
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

// QuickVoteHandler create a quick vote from command
func QuickVoteHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	options := strings.Fields(m.Text)
	if len(options) < 4 {
		bot.Send(m.Chat, "Please follow format: /vote MyTitle option1 option2")
		return

	} else if len(options) > 8 {
		bot.Send(m.Chat, "Maximum of options is 6")
		return
	}

	// Vote Content
	vote := Vote{
		Title:     options[1],
		Options:   "",
		Voters:    0,
		CreatedAt: time.Now(),
	}
	vote.Options = strings.Join(options[2:], "\n")

	err := db.Insert(&vote)
	if err != nil {
		panic(err)
	}
	if vote.ID == 0 {
		panic("Insert vote failed")
	}

	choices := []Choice{}

	// Response
	m, err = bot.Send(
		m.Chat,
		vote.String(choices),
		&tb.ReplyMarkup{InlineKeyboard: vote.GenButton()},
	)
	if err != nil {
		panic(err)
	}

	// Add to Database
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
