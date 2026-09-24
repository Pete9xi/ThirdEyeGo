package text

import (
	"fmt"

	"thirdeyego/discord"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (h *Handler) handlePlayerJoinLeave(p *packet.Text) {
	switch p.Message {
	case "§e%multiplayer.player.joined":
		if len(p.Parameters) == 0 {
			return
		}

		playerName := p.Parameters[0]
		message := fmt.Sprintf("[In Game] %s joined the server", playerName)

		fmt.Println(message)

		if err := discord.SendMessage(h.cfg.Channel, message); err != nil {
			fmt.Println("Failed to send Discord message:", err)
		}

	case "§e%multiplayer.player.left":
		if len(p.Parameters) == 0 {
			return
		}

		playerName := p.Parameters[0]
		message := fmt.Sprintf("[In Game] %s left the server", playerName)

		fmt.Println(message)

		if err := discord.SendMessage(h.cfg.Channel, message); err != nil {
			fmt.Println("Failed to send Discord message:", err)
		}

	}
}
