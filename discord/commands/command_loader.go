package commands

import "github.com/bwmarrin/discordgo"

func Load() map[string]Command {
	commandList := []Command{
		minecraftCommand,
	}

	commandMap := make(map[string]Command)

	for _, command := range commandList {
		commandMap[command.Data.Name] = command
	}

	return commandMap
}

func Data(commands map[string]Command) []*discordgo.ApplicationCommand {
	data := make([]*discordgo.ApplicationCommand, 0, len(commands))

	for _, command := range commands {
		data = append(data, command.Data)
	}

	return data
}
