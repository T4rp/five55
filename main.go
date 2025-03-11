package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatal("No token provided")
	}

	test_guild := os.Getenv("TEST_GUILD")
	if test_guild == "" {
		log.Fatal("No test guild id provided")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: m.Content,
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Rahh",
				Description: "Rahhhhh",
				Color:       0xaaff96,
			}},
		})
		if err != nil {
			log.Println(err)
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
}
