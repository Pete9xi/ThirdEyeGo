package text

import (
	"thirdeyego/config"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Handler struct {
	cfg config.Config
}

func NewHandler(cfg config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) HandleText(p *packet.Text) {
	switch p.TextType {

	case packet.TextTypeChat:
		// h.handleChat(p)

	case packet.TextTypeWhisper:
		// h.handleWhisper(p)

	case packet.TextTypeObject:
		h.handleObject(p)
		h.handleAntiCheatPacket(p)

	case packet.TextTypeTranslation:
		h.handlePlayerJoinLeave(p)

	default:
		// Ignore other text types for now.
	}
}
