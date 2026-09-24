package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session

func Start(token string) error {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	Session = session

	if err := Session.Open(); err != nil {
		return err
	}

	fmt.Println("Discord bot connected!")

	return nil
}

func SendMessage(channelID string, message string) error {
	_, err := Session.ChannelMessageSend(channelID, message)
	return err
}
