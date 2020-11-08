package automoji

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hackebrot/turtle"
)

/*
Run comments to discord and reacts to all messages with vaguely related emoji.

Command: auto-emote
*/
func Run(args []string, _ io.Reader) error {
	if len(args) > 1 {
		for _, arg := range args[1:] {
			fmt.Println(messageToEmoji(arg))
		}
		return nil
	}

	token := os.Getenv("LM_DISCORD_TOKEN")
	if token == "" {
		return errors.New("set LM_DISCORD_TOKEN to use auto-emote")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	if err := dg.Open(); err != nil {
		return err
	}

	x := make(chan bool)

	<-x

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || m.Message == nil {
		return
	}

	emoji := messageToEmoji(m.Message.Content)
	for i, e := range emoji {
		// 20 is max, so limit to half the total amount
		if i == 10 {
			break
		}
		s.MessageReactionAdd(
			m.ChannelID,
			m.ID,
			e,
		)
	}
}

type emojiSet map[string]bool

func (s emojiSet) add(e *turtle.Emoji) {
	if e.Category == "flags" {
		return
	}

	s[e.Char] = true
}

func (s emojiSet) all() []string {
	ret := make([]string, 0, len(s))
	for e := range s {
		ret = append(ret, e)
	}

	return ret
}

func messageToEmoji(m string) []string {
	words := strings.Split(strings.ToLower(m), " ")
	s := emojiSet(make(map[string]bool, len(words)))

	for _, word := range words {
		if word == "" {
			continue
		}

		if e, ok := turtle.Emojis[word]; ok {
			s.add(e)
		}

		if es := turtle.Category(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}

		if es := turtle.Keyword(word); es != nil {
			for _, e := range es {
				s.add(e)
			}
		}

		// don't use turtle.Search because it finds *way* too much.
	}

	return s.all()
}
