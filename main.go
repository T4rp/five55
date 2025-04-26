package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	MAIN_CHANNEL         = "1349197130912235554"
	TOGGLE_CATAGORY      = "1349197130912235552"
	TOGGLE_ROLE          = "1349197130912235551"
	CHATTING_PERIMSSIONS = discordgo.PermissionSendMessages |
		discordgo.PermissionSendMessagesInThreads |
		discordgo.PermissionAddReactions |
		discordgo.PermissionVoiceConnect |
		discordgo.PermissionVoiceSpeak |
		discordgo.PermissionVoiceStreamVideo |
		discordgo.PermissionUseActivities
)

const (
	START_HOUR   = 12 + 5 + 4 // 5 pm EST
	START_MINUTE = 55
	END_HOUR     = 12 + 7 + 4 // 7 pm EST
	END_MINUTE   = 55
)

func getChannelGlobalPermissions(session *discordgo.Session) (*discordgo.PermissionOverwrite, error) {
	channel, err := session.Channel(MAIN_CHANNEL)
	if err != nil {
		return nil, err
	}

	permissions := channel.PermissionOverwrites

	for _, perm := range permissions {
		if perm.ID == TOGGLE_ROLE {
			return perm, nil
		}
	}

	return nil, nil
}

func updateChannels(session *discordgo.Session) {
	perms, err := getChannelGlobalPermissions(session)
	if err != nil {
		log.Fatal(err)
	}

	if isTime() {
		session.ChannelPermissionSet(
			TOGGLE_CATAGORY,
			TOGGLE_ROLE,
			discordgo.PermissionOverwriteTypeRole,
			perms.Allow|CHATTING_PERIMSSIONS,
			perms.Deny&(^CHATTING_PERIMSSIONS),
		)

		session.ChannelMessageSend(MAIN_CHANNEL, "The server is open @evenyone")
	} else {
		session.ChannelPermissionSet(TOGGLE_CATAGORY,
			TOGGLE_ROLE,
			discordgo.PermissionOverwriteTypeRole,
			perms.Allow&(^CHATTING_PERIMSSIONS),
			perms.Deny|CHATTING_PERIMSSIONS,
		)

		session.ChannelMessageSend(MAIN_CHANNEL, "The server has been locked. Please wait until 5:55 PM EST to chat")
	}
}

func isTime() bool {
	currentTime := time.Now().UTC()

	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()

	isAfterStart := currentHour > START_HOUR || (currentHour == START_HOUR && currentMinute >= START_MINUTE)
	isBeforeEnd := currentHour < END_HOUR || (currentHour == END_HOUR && currentMinute < END_MINUTE)

	return isAfterStart && isBeforeEnd
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
				lastState = newState
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
