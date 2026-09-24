package text

import (
	"encoding/json"
	"fmt"
	"strings"
	"thirdeyego/bedrock/utils"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// handleAntiCheatPacket checks JSON text packets for Paradox anti-cheat messages.
//
// If the message is not an anti-cheat message, the function returns immediately.
func (h *Handler) handleAntiCheatPacket(p *packet.Text) {
	var obj RawTextMessage

	// Ignore invalid JSON.
	if err := json.Unmarshal([]byte(p.Message), &obj); err != nil {
		return
	}

	// Check each rawtext entry for an anti-cheat prefix.
	for _, part := range obj.RawText {
		isAntiCheatMessage :=
			strings.Contains(part.Text, "§2[§7Available Commands§2]§r") ||
				strings.Contains(part.Text, "§2[§7Paradox§2]§o§7")

		if !isAntiCheatMessage {
			continue
		}

		// We found an anti-cheat message.
		// Process it here.
		message := part.Text

		message = utils.AutoCorrect(message)

		// Forward the message to the configured Discord channel.
		if err := h.discord.SendMessage(h.cfg.AntiCheatLogsChannel, message); err != nil {
			fmt.Println("Failed to send Discord message:", err)
		}

		return
	}
}
