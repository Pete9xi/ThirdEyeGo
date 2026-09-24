package text

import (
	"fmt"
	"thirdeyego/discord"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (h *Handler) handlePlayerJoinLeave(p *packet.Text) {
	var message string
	var color []int

	switch p.Message {
	case "§e%multiplayer.player.joined":
		if len(p.Parameters) == 0 {
			return
		}

		playerName := p.Parameters[0]
		message = fmt.Sprintf("[In Game] %s joined the server", playerName)
		color = []int{0, 255, 0}

	case "§e%multiplayer.player.left":
		if len(p.Parameters) == 0 {
			return
		}

		playerName := p.Parameters[0]
		message = fmt.Sprintf("[In Game] %s left the server", playerName)
		color = []int{255, 0, 0}

	default:
		return
	}

	fmt.Println(message)

	if h.cfg.UseEmbed {
		embed := discord.CreateEmbed(discord.EmbedOptions{
			Title:       h.cfg.SetTitle,
			Description: message,
			Color:       color,
		})

		if err := h.discord.SendEmbed(h.cfg.Channel, embed); err != nil {
			fmt.Println("Failed to send Discord embed:", err)
		}

		return
	}

	if err := h.discord.SendMessage(h.cfg.Channel, message); err != nil {
		fmt.Println("Failed to send Discord message:", err)
	}
}
