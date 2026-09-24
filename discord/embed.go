package discord

import "github.com/bwmarrin/discordgo"

type EmbedOptions struct {
	Title         string
	Description   string
	Color         []int
	AuthorName    string
	AuthorIconURL string
	ThumbnailURL  string
}

func CreateEmbed(options EmbedOptions) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}

	if options.Title != "" {
		embed.Title = options.Title
	}

	if options.Description != "" {
		embed.Description = options.Description
	}

	if len(options.Color) >= 3 {
		embed.Color = RGB(
			options.Color[0],
			options.Color[1],
			options.Color[2],
		)
	}

	if options.AuthorName != "" {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    options.AuthorName,
			IconURL: options.AuthorIconURL,
		}
	}

	if options.ThumbnailURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: options.ThumbnailURL,
		}
	}

	return embed
}

func RGB(r, g, b int) int {
	return (r << 16) | (g << 8) | b
}
