package commands

import (
	"fmt"
	"thirdeyego/bedrock/utils"

	"github.com/bwmarrin/discordgo"
)

var minecraftCommand = Command{
	Data: &discordgo.ApplicationCommand{
		Name:        "minecraftCommand",
		Description: "executes a command on the Minecraft server",

		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "The Minecraft command to execute",
				Required:    true,
			},
		},
	},

	Execute: func(
		ctx *Context,
		interaction *discordgo.InteractionCreate,
	) {
		// Role check
		if !HasRole(
			interaction,
			ctx.Config.OperatorsRoleID,
		) {
			_ = ctx.Discord.InteractionRespond(
				interaction.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ You do not have permission to use this command.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				},
			)

			return
		}
		options := interaction.ApplicationCommandData().Options

		command := options[0].StringValue()

		fmt.Println("[Discord Command]", command)

		err := utils.RunCMD(
			ctx.Minecraft,
			"/"+command,
		)

		if err != nil {
			fmt.Println("Failed to execute Minecraft command:", err)

			_ = ctx.Discord.InteractionRespond(
				interaction.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ Failed to execute Minecraft command.",
					},
				},
			)

			return
		}

		_ = ctx.Discord.InteractionRespond(
			interaction.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Command executed on the Minecraft server.",
				},
			},
		)
	},
}
