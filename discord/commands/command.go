package commands

import (
	"thirdeyego/config"

	"github.com/bwmarrin/discordgo"
	"github.com/sandertv/gophertunnel/minecraft"
)

type Context struct {
	Discord   *discordgo.Session
	Minecraft *minecraft.Conn
	Config    config.Config
}

type Command struct {
	Data    *discordgo.ApplicationCommand
	Execute func(*Context, *discordgo.InteractionCreate)
}
