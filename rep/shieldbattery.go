// This file contains the types describing info parsed from the ShieldBattery custom section.

package rep

import "encoding/base64"

// ShieldBattery models the data parsed from the ShieldBattery custom section.
type ShieldBattery struct {
	FormatVersion        uint16
	StarCraftExeBuild    uint32
	ShieldBatteryVersion string
	UserIDs              []uint32
	RawGameID            []byte
	GameID               string

	// Available for FormatVersion >= 1
	GameLogicVersion *uint16 `json:",omitempty"`

	// DisabledLogicFixes is a bitmap, each bit has its own meaning. Available for FormatVersion >= 2
	DisabledLogicFixes *uint32 `json:",omitempty"`

	PidGameID map[byte]uint32 `json:"-"`
}

// UserID returns the ShieldBattery user ID for the given player ID (pid).
func (sb *ShieldBattery) UserID(pid byte) uint32 {
	if uid, ok := sb.PidGameID[pid]; ok {
		return uid
	}
	return 0
}

func (sb *ShieldBattery) RawGameIDURLString() string {
	return base64.RawURLEncoding.EncodeToString(sb.RawGameID)
}
