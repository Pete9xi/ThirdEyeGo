package packets

import (
	"fmt"
	"thirdeyego/bedrock/packets/text"
	"thirdeyego/config"
	"thirdeyego/discord"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Handler struct {
	cfg         config.Config
	textHandler *text.Handler
	discord     *discord.Client
}

func NewHandler(cfg config.Config, discordClient *discord.Client) *Handler {
	return &Handler{
		cfg:         cfg,
		textHandler: text.NewHandler(cfg, discordClient),
		discord:     discordClient,
	}
}

func StartPacketLoop(
	conn *minecraft.Conn,
	cfg config.Config,
	discordClient *discord.Client,
) {
	fmt.Println("Starting packet handler...")

	handler := NewHandler(cfg, discordClient)

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			fmt.Printf("Connection closed: %v\n", err)
			return
		}

		handler.HandlePacket(pk)
	}
}

func (h *Handler) HandlePacket(pk any) {
	h.debugPacket(pk)

	switch p := pk.(type) {

	case *packet.Text:
		h.textHandler.HandleText(p)

	default:
		// Ignore packets we don't care about.
	}
}

// debugPacket logs incoming packets when debug mode is enabled.
func (h *Handler) debugPacket(pk any) {
	if !h.cfg.Debug {
		return
	}

	packetName := fmt.Sprintf("%T", pk)

	if isSpammyPacket(packetName) {
		return
	}

	if p, ok := pk.(*packet.Text); ok {
		fmt.Printf("[Packet] %s\n", packetName)
		fmt.Printf("  Type: %d\n", p.TextType)
		fmt.Printf("  Source: %s\n", p.SourceName)
		fmt.Printf("  Message: %s\n", p.Message)
		fmt.Printf("  Parameters: %v\n", p.Parameters)
	}
}
