package commands

import "github.com/bwmarrin/discordgo"

func HasRole(interaction *discordgo.InteractionCreate, roleID string) bool {
	for _, role := range interaction.Member.Roles {
		if role == roleID {
			return true
		}
	}

	return false
}

func InChannel(interaction *discordgo.InteractionCreate, channelID string) bool {
	return interaction.ChannelID == channelID
}
