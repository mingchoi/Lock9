package main

import (
	"fmt"
	"strconv"
	"strings"

	tb "github.com/tucnak/telebot"
)

// VoteRef is
type VoteRef struct {
	VoteID    int
	ChatID    int
	MessageID int
}

// Vote is
type Vote struct {
	ID          int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	Title       string
	Description string
	Options     string
	Multiple    bool
	Voters      int
}

// Choice is
type Choice struct {
	ID       int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	VoteID   int
	UserID   int
	UserName string
	Option   int
}

// String is
func (vote *Vote) String(choices []Choice) string {
	var str string
	if vote.Description != "" {
		str = fmt.Sprintf("📊 **%s**\n%s\n\n", vote.Title, vote.Description)
	} else {
		str = fmt.Sprintf("📊 %s\n\n", vote.Title)
	}

	options := strings.Split(vote.Options, "\n")
	for i := range options {
		count := 0
		nameList := ""
		for j := range choices {
			if choices[j].Option == i+1 {
				count++
				nameList += "\n    - " + choices[j].UserName
			}
		}
		str += strconv.Itoa(i+1) + ". " + options[i] + " [" + strconv.Itoa(count) + "]" + nameList + "\n\n"
	}
	str += "Votes: " + strconv.Itoa(vote.Voters)
	return str
}

// GenButton is
func (vote *Vote) GenButton() [][]tb.InlineButton {
	btns := make([][]tb.InlineButton, 0)
	options := strings.Split(vote.Options, "\n")
	for i := range options {
		btns = append(btns, []tb.InlineButton{
			tb.InlineButton{
				Unique: "a" + strconv.Itoa(i+1),
				Text:   options[i],
			},
		})
	}
	return btns
}

func quickVoteHandler(m *tb.Message) {
	options := strings.Fields(m.Text)
	if len(options) < 4 {
		b.Send(m.Sender, "Please follow format: /vote MyTitle option1 option2")
		return

	} else if len(options) > 8 {
		b.Send(m.Sender, "Maximum of options is 6")
		return
	}

	// Vote Content
	vote := Vote{
		Title:   options[1],
		Options: "",
		Voters:  0,
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
	m, err = b.Send(
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

func advVoteHandler(m *tb.Message) {
	options := strings.Fields(m.Text)
	if len(options) < 5 || !(options[1] == "single" || options[1] == "multiple") {
		b.Send(m.Sender, "Please follow format: /vote {single|multiple} MyTitle option1 option2")
		return

	} else if len(options) > 9 {
		b.Send(m.Sender, "Maximum of options is 6")
		return
	}

	// Vote Content
	vote := Vote{
		Title:    options[2],
		Options:  "",
		Multiple: (options[1] == "multiple"),
		Voters:   0,
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
	m, err = b.Send(
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

func forwareVoteHandler(m *tb.Message) {
	options := strings.Fields(m.Text)
	if len(options) != 2 {
		b.Send(m.Sender, "Please follow format: /forwardvote {VoteID}")
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
		b.Send(m.Sender, "Vote not found.")
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
	m, err = b.Send(
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

func handleVote(chatID int, messageID int, userID int, userName string, option int) {
	// Search for Vote
	var ref VoteRef
	err := db.Select(&ref, "chatid = ? AND messageid = ?", chatID, messageID)
	if err != nil {
		panic(err)
	}
	var vote Vote
	err = db.Select(&vote, "id = ?", ref.VoteID)
	if err != nil {
		panic(err)
	}
	if vote.ID == 0 {
		panic("Vote not found")
	}

	if vote.Multiple {
		var choice Choice
		err := db.Select(&choice, "voteid = ? AND userid = ? AND option = ?", vote.ID, userID, option)
		if err != nil {
			panic(err)
		}

		if choice.ID != 0 {
			err := db.Delete(&choice, "id = ?", choice.ID)
			if err != nil {
				panic(err)
			}
		} else {
			err := db.Insert(&Choice{
				VoteID:   vote.ID,
				UserID:   userID,
				UserName: userName,
				Option:   option,
			})
			if err != nil {
				panic(err)
			}
		}
	} else {
		var choice Choice
		err := db.Select(&choice, "voteid = ? AND userid = ?", vote.ID, userID)
		if err != nil {
			panic(err)
		}

		if choice.ID != 0 {
			if choice.Option == option {
				err := db.Delete(&choice, "id = ?", choice.ID)
				if err != nil {
					panic(err)
				}
			} else {
				choice.Option = option
				err := db.Update(&choice, "id = ?", choice.ID)
				if err != nil {
					panic(err)
				}
			}
		} else {
			err := db.Insert(&Choice{
				VoteID:   vote.ID,
				UserID:   userID,
				UserName: userName,
				Option:   option,
			})
			if err != nil {
				panic(err)
			}
		}
	}

	// Update message
	updateVoteMessage(vote)
}

func updateVoteMessage(vote Vote) {
	// Load Choices
	choices := []Choice{}
	err := db.Select(&choices, "voteid = ?", vote.ID)
	if err != nil {
		panic(err)
	}
	vote.Voters = len(choices)

	// Look for message reference
	var refs []VoteRef
	err = db.Select(&refs, "voteid = ?", vote.ID)
	if err != nil {
		panic(err)
	}

	// Edit message
	for i := range refs {
		_, err = b.Edit(
			tb.StoredMessage{
				ChatID:    int64(refs[i].ChatID),
				MessageID: strconv.Itoa(refs[i].MessageID),
			},
			vote.String(choices),
			&tb.ReplyMarkup{InlineKeyboard: vote.GenButton()},
		)
		if err != nil {
			panic(err)
		}
	}
}

func btnHandler1(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		1,
	)
}
func btnHandler2(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		2,
	)
}
func btnHandler3(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		3,
	)
}
func btnHandler4(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		4,
	)
}
func btnHandler5(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		5,
	)
}
func btnHandler6(c *tb.Callback) {
	handleVote(
		int(c.Message.Chat.ID),
		c.Message.ID, c.Sender.ID,
		c.Sender.FirstName+" "+c.Sender.LastName,
		6,
	)
}
