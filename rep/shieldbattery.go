// This file contains the types describing info parsed from the ShieldBattery custom section.

package rep

// ShieldBattery models the data parsed from the ShieldBattery custom section.
type ShieldBattery struct {
	FormatVersion        uint16
	StarCraftExeBuild    uint32
	ShieldBatteryVersion string
	GameID               string

	// Available for FormatVersion >= 1
	GameLogicVersion *uint16 `json:",omitempty"`

	// DisabledLogicFixes is a bitmask, each bit has its own meaning. Available for FormatVersion >= 2
	DisabledLogicFixes *uint32 `json:",omitempty"`
}
