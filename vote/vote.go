package vote

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	s2s "github.com/mingchoi/struct2sql"
	tb "github.com/tucnak/telebot"
)

// VoteRef is
type VoteRef struct {
	VoteID    int `foreignkey:"vote(id)"`
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Choice is
type Choice struct {
	ID        int `keyword:"NOT NULL AUTO_INCREMENT" primarykey:"true"`
	VoteID    int `foreignkey:"vote(id)"`
	UserID    int
	UserName  string
	Option    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// String prints a vote status to string
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

// GenButton create button for vote options
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

// Run when user vote by clicking button
func handleVote(chatID int, messageID int, userID int, userName string, option int, bot *tb.Bot, db *s2s.DB) {
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
				VoteID:    vote.ID,
				UserID:    userID,
				UserName:  userName,
				Option:    option,
				CreatedAt: time.Now(),
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
				choice.UpdatedAt = time.Now()
				err := db.Update(&choice, "id = ?", choice.ID)
				if err != nil {
					panic(err)
				}
			}
		} else {
			err := db.Insert(&Choice{
				VoteID:    vote.ID,
				UserID:    userID,
				UserName:  userName,
				Option:    option,
				CreatedAt: time.Now(),
			})
			if err != nil {
				panic(err)
			}
		}
	}

	// Update message
	updateVoteMessage(vote, bot, db)
}

// Update vote message
func updateVoteMessage(vote Vote, bot *tb.Bot, db *s2s.DB) {
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
		_, err = bot.Edit(
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
