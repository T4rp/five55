package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
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

var startTime time.Time

type Config struct {
	token     string
	appId     string
	testGuild string

	mainChannel    string
	toggleCatagory string
	toggleRole     string
}

func assertEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("%s not set", key)
	}

	return val
}

func parseConfig() Config {
	token := assertEnv("TOKEN")
	appId := assertEnv("APP_ID")
	testGuild := assertEnv("TEST_GUILD")

	mainChannel := assertEnv("MAIN_CHANNEL")
	toggleCatagory := assertEnv("TOGGLE_CATAGORY")
	toggleRole := assertEnv("TOGGLE_ROLE")

	return Config{
		token:          token,
		appId:          appId,
		testGuild:      testGuild,
		mainChannel:    mainChannel,
		toggleCatagory: toggleCatagory,
		toggleRole:     toggleRole,
	}
}

func getChannelRolePermissions(session *discordgo.Session, channelId string, roleId string) (*discordgo.PermissionOverwrite, error) {
	channel, err := session.Channel(channelId)
	if err != nil {
		return nil, err
	}

	permissions := channel.PermissionOverwrites

	for _, perm := range permissions {
		if perm.ID == roleId {
			return perm, nil
		}
	}

	return nil, nil
}

func updateChannels(session *discordgo.Session, config *Config) {
	perms, err := getChannelRolePermissions(session, config.toggleCatagory, config.toggleRole)
	if err != nil {
		log.Fatal(err)
	}

	if isTime() {
		session.ChannelPermissionSet(
			config.toggleCatagory,
			config.toggleRole,
			discordgo.PermissionOverwriteTypeRole,
			perms.Allow|CHATTING_PERIMSSIONS,
			perms.Deny&(^CHATTING_PERIMSSIONS),
		)
	} else {
		session.ChannelPermissionSet(
			config.toggleCatagory,
			config.toggleRole,
			discordgo.PermissionOverwriteTypeRole,
			perms.Allow&(^CHATTING_PERIMSSIONS),
			perms.Deny|CHATTING_PERIMSSIONS,
		)
	}
}

func sendOpenMessage(session *discordgo.Session, config *Config) {
	timeNow := time.Now()
	closeDate := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), END_HOUR, END_MINUTE, 0, 0, time.UTC)

	session.ChannelMessageSendComplex(config.mainChannel, &discordgo.MessageSend{
		Content: "@everyone",
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "Five55 is now open!",
			Description: fmt.Sprintf("Five55 will remain open until <t:%d:t>.", closeDate.Unix()),
			Color:       0xe5eb42,
		}},
	})
}

func sendCloseMessage(session *discordgo.Session, config *Config) {
	timeNow := time.Now()
	openDate := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), START_HOUR, START_MINUTE, 0, 0, time.UTC)
	openDate = openDate.Add(24 * time.Hour)

	session.ChannelMessageSendComplex(config.mainChannel, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "Five55 is now closed",
			Description: fmt.Sprintf("Five55 will open again on <t:%d:f>.", openDate.Unix()),
			Color:       0x3348bd,
		}},
	})
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
	startTime = time.Now()

	config := parseConfig()

	token := config.token
	appId := config.appId
	testGuild := config.testGuild

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.Interaction.ApplicationCommandData().Name

		if commandName == "info" {
			uptime := time.Since(startTime)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{{
						Title:       "Five55",
						Description: "Five55 is a bot that manages the five55 discord server.",
						Color:       0x990000,
						Fields: []*discordgo.MessageEmbedField{{
							Name:  "Uptime",
							Value: uptime.Abs().String(),
						}},
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
		updateChannels(session, &config)

		lastState := isTime()

		for {
			newState := isTime()

			if newState != lastState {
				lastState = newState
				updateChannels(session, &config)

				if newState {
					sendOpenMessage(session, &config)
				} else {
					sendCloseMessage(session, &config)
				}
			}

			time.Sleep(10 * time.Second)
		}
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	log.Println("Stopping...")
}
