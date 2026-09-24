package packets

import "strings"

// spammyPackets contains packet types that are received very frequently
// and would make debug output difficult to read.
var spammyPackets = []string{
	"packet.MovePlayer",
	"packet.PlayerAuthInput",
	"packet.NetworkStackLatency",
	"packet.MoveActorDelta",
	"packet.UpdateSubChunkBlocks",
	"packet.NetworkChunkPublisherUpdate",
	"packet.SetActorData",
	"packet.LevelSoundEvent",
	"packet.BlockActorData",
	"packet.SetTitle",
	"packet.LevelChunk",
	"packet.SetActorMotion",
	"packet.AddActor",
	"packet.AddItemActor",
	"packet.MobEquipment",
	"packet.MobArmourEquipment",
	"packet.AddPlayer",
	"packet.CurrentStructureFeature",
	"packet.LocatorBar",
	"packet.UpdateAttributes",
	"packet.RemoveActor",
	"packet.UpdateBlock",
	"packet.LevelEvent",
	"packet.MobEffect",
	"packet.TakeItemActor",
	"packet.PlaySound",
}

// isSpammyPacket determines whether a packet should be excluded from
// debug logging.
func isSpammyPacket(packetName string) bool {
	for _, spammy := range spammyPackets {
		if strings.Contains(packetName, spammy) {
			return true
		}
	}

	return false
}
