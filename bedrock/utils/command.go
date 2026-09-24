package utils

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// RunCMD sends a command to the connected Minecraft server.
func RunCMD(conn *minecraft.Conn, cmd string) error {
	return conn.WritePacket(&packet.CommandRequest{
		CommandLine: cmd,
		CommandOrigin: protocol.CommandOrigin{
			Origin:         protocol.CommandOriginPlayer,
			UUID:           uuid.Nil,
			RequestID:      "",
			PlayerUniqueID: 0,
		},
		Internal: false,
		Version:  "latest",
	})
}
