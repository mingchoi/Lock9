package vote

import (
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

// AdvVoteHandler create a advanced vote from command
func AdvVoteHandler(m *tb.Message, bot *tb.Bot, db *s2s.DB) {
	options := strings.Fields(m.Text)
	if len(options) < 5 || !(options[1] == "single" || options[1] == "multiple") {
		bot.Send(m.Sender, "Please follow format: /vote {single|multiple} MyTitle option1 option2")
		return

	} else if len(options) > 9 {
		bot.Send(m.Sender, "Maximum of options is 6")
		return
	}

	// Vote Content
	vote := Vote{
		Title:     options[2],
		Options:   "",
		Multiple:  (options[1] == "multiple"),
		Voters:    0,
		CreatedAt: time.Now(),
	}
	vote.Options = strings.Join(options[3:], "\n")

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
