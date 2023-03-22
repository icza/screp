// This file contains the types describing info parsed from the ShieldBattery custom section.

package rep

// ShieldBattery models the data parsed from the ShieldBattery custom section.
type ShieldBattery struct {
	StarCraftExeBuild    uint32
	ShieldBatteryVersion string
	GameID               string
}
