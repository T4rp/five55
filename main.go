package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("TOKEN")

	if token == "" {
		log.Fatal("No token provided")
	}

	appId := os.Getenv("APP_ID")
	if token == "" {
		log.Fatal("No application id provided")
	}

	testGuild := os.Getenv("TEST_GUILD")
	if testGuild == "" {
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

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.Interaction.ApplicationCommandData().Name

		if commandName == "info" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "hello world",
				},
			})
		} else {
			log.Println("Cannot handle interaction with id: " + commandName)
		}
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Logged in as " + r.User.Username)
	})

	_, err = session.ApplicationCommandCreate(appId, testGuild, &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "ARarararh",
	})
	if err != nil {
		log.Println("Failed to create application command:", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	log.Println("Stopping...")
}
