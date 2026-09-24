package utils

import "strings"

var correction = strings.NewReplacer(
	// Paradox / plugin prefixes
	"§2[§7Paradox§2]§o§7", "Paradox",
	"§2[§7Member§2]", "Member",

	// Formatting
	"§l", "", // bold
	"§o", "", // italic
	"§k", "", // obfuscated
	"§r", "", // reset

	// Color codes
	"§0", "",
	"§1", "",
	"§2", "",
	"§3", "",
	"§4", "",
	"§5", "",
	"§6", "",
	"§7", "",
	"§8", "",
	"§9", "",
	"§a", "",
	"§b", "",
	"§c", "",
	"§d", "",
	"§e", "",
	"§f", "",
	"§g", "",
	"§h", "",
	"§i", "",
	"§j", "",
	"§m", "",
	"§n", "",
	"§p", "",
	"§q", "",
	"§s", "",
	"§t", "",
	"§u", "",

	// Custom prefixes
	"", "[Console]",
	"", "[PC]",
	"", "[Mobile]",

	// Rare / malformed formatting
	"§¶", "",
)

func AutoCorrect(text string) string {
	return correction.Replace(text)
}
