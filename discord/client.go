package discord

import (
	"fmt"

	"thirdeyego/bedrock/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/sandertv/gophertunnel/minecraft"
)

type Client struct {
	Session *discordgo.Session
	Conn    *minecraft.Conn
}

func Start(token string, conn *minecraft.Conn) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Session: session,
		Conn:    conn,
	}

	session.AddHandler(client.handleMessageCreate)

	if err := session.Open(); err != nil {
		return nil, err
	}

	fmt.Println("Discord bot connected!")

	return client, nil
}

func (client *Client) SendMessage(channelID string, message string) error {
	_, err := client.Session.ChannelMessageSend(channelID, message)
	return err
}
func (client *Client) SendEmbed(channelID string, embed *discordgo.MessageEmbed) error {
	_, err := client.Session.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func (client *Client) handleMessageCreate(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
) {
	if message.Author.Bot {
		return
	}

	fmt.Printf("[Discord] %s: %s\n", message.Author.Username, message.Content)
	cmd := fmt.Sprintf(
		`/tellraw @a {"rawtext":[{"text":"§8[§9Discord§8] §7%s: §f%s"}]}`,
		message.Author.Username,
		message.Content,
	)

	utils.RunCMD(client.Conn, cmd)
}
