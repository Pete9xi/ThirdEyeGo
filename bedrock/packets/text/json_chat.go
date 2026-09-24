package text

import (
	"encoding/json"
	"fmt"
	"strings"

	"thirdeyego/bedrock/utils"
	"thirdeyego/discord"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ignoredPrefixes contains Minecraft messages that should not be forwarded
// to Discord because they are system/server messages.
var ignoredPrefixes = []string{
	"§2[§7Available Commands§2]§r",
	"§2[§7Paradox§2]§o§7",
	"?link:",
}

// RawTextMessage represents the JSON structure used by Minecraft's rawtext
// messages.
type RawTextMessage struct {
	RawText []struct {
		Text      string `json:"text"`
		Translate string `json:"translate"`
	} `json:"rawtext"`
}

// handleObject handles Minecraft TextTypeObject packets.
//
// Object packets contain a JSON message in p.Message. The JSON is parsed,
// converted into normal text, filtered, corrected, and then sent to Discord.
func (h *Handler) handleObject(p *packet.Text) {
	var obj RawTextMessage

	if err := json.Unmarshal([]byte(p.Message), &obj); err != nil {
		fmt.Println("Failed to parse JSON chat:", err)
		return
	}

	var text strings.Builder

	// Combine all text entries from the rawtext array.
	for _, part := range obj.RawText {
		// Ignore messages that use translation keys.
		if part.Translate != "" {
			return
		}

		text.WriteString(part.Text)
	}

	// Get the complete raw message.
	message := text.String()

	// Ignore empty messages.
	if message == "" {
		return
	}

	// Ignore known Minecraft/server system messages.
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(message, prefix) {
			return
		}
	}

	// Apply Minecraft formatting corrections.
	message = utils.AutoCorrect(message)

	fmt.Println("[In Game]", message)

	// Forward the message to the configured Discord channel.
	if err := discord.SendMessage(h.cfg.Channel, "[In Game] "+message); err != nil {
		fmt.Println("Failed to send Discord message:", err)
	}
}
