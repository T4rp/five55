package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	TOGGLE_CATAGORY = "0"
	TOGGLE_ROLE     = "0"
)

func updateChannels(session *discordgo.Session) {
	if isTime() {
		session.ChannelPermissionSet(TOGGLE_CATAGORY, TOGGLE_ROLE, discordgo.PermissionOverwriteTypeRole, 1, 0)
	} else {
		session.ChannelPermissionSet(TOGGLE_CATAGORY, TOGGLE_ROLE, discordgo.PermissionOverwriteTypeRole, 0, 1)
	}
}

func isTime() bool {
	currentTime := time.Now().UTC()
	isValidTime := currentTime.Hour() >= 21 && currentTime.Minute() >= 55 && currentTime.Hour() < 23

	return isValidTime
}

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

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.Interaction.ApplicationCommandData().Name

		if commandName == "info" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{{
						Title: "five55 Bot",
						Description: "Five55 is a bot that manages the five55 discord server.\n" +
							"This bot is supposed to close the discord server channels at 5:55 PM EST",
						Color: 0x990000,
					}},
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

	go func() {
		updateChannels(session)

		lastState := isTime()

		for {
			newState := isTime()
			if newState != lastState {
				updateChannels(session)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	log.Println("Stopping...")
}
