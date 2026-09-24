package discord

import (
	"fmt"

	"thirdeyego/bedrock/utils"
	"thirdeyego/config"
	"thirdeyego/discord/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/sandertv/gophertunnel/minecraft"
)

type Client struct {
	Session  *discordgo.Session
	Conn     *minecraft.Conn
	Config   config.Config
	Commands map[string]commands.Command
}

func Start(
	token string,
	conn *minecraft.Conn,
	cfg config.Config,
) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Session:  session,
		Conn:     conn,
		Config:   cfg,
		Commands: commands.Load(),
	}

	for name := range client.Commands {
		fmt.Println("Loaded command:", name)
	}

	session.AddHandler(client.handleMessageCreate)
	session.AddHandler(client.handleInteractionCreate)

	if err := session.Open(); err != nil {
		return nil, err
	}

	fmt.Println("Discord bot connected!")

	if err := client.registerCommands(cfg); err != nil {
		return nil, err
	}

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

func (client *Client) handleInteractionCreate(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := interaction.ApplicationCommandData().Name

	command, ok := client.Commands[commandName]
	if !ok {
		fmt.Println("Unknown command:", commandName)
		return
	}

	ctx := &commands.Context{
		Discord:   session,
		Minecraft: client.Conn,
		Config:    client.Config,
	}

	command.Execute(ctx, interaction)
}

func (client *Client) registerCommands(cfg config.Config) error {
	appID := client.Session.State.User.ID

	commandData := commands.Data(client.Commands)

	if cfg.IsDev {
		if cfg.Guild == "" {
			return fmt.Errorf("dev mode requires a guild ID")
		}

		_, err := client.Session.ApplicationCommandBulkOverwrite(
			appID,
			cfg.Guild,
			commandData,
		)

		if err != nil {
			return fmt.Errorf("register guild commands: %w", err)
		}

		fmt.Println("🧪 DEV MODE: Slash commands registered to guild")

		return nil
	}

	_, err := client.Session.ApplicationCommandBulkOverwrite(
		appID,
		"",
		commandData,
	)

	if err != nil {
		return fmt.Errorf("register global commands: %w", err)
	}

	fmt.Println("🚀 PROD MODE: Slash commands registered globally")

	return nil
}
