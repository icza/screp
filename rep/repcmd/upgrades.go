// This file contains upgrades.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Upgrade describes the upgrade.
type Upgrade struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// Upgrades is an enumeration of the possible upgrades.
var Upgrades = []*Upgrade{
	{e("Terran Infantry Armor"), 0x00},
	{e("Terran Vehicle Plating"), 0x01},
	{e("Terran Ship Plating"), 0x02},
	{e("Zerg Carapace"), 0x03},
	{e("Zerg Flyer Carapace"), 0x04},
	{e("Protoss Ground Armor"), 0x05},
	{e("Protoss Air Armor"), 0x06},
	{e("Terran Infantry Weapons"), 0x07},
	{e("Terran Vehicle Weapons"), 0x08},
	{e("Terran Ship Weapons"), 0x09},
	{e("Zerg Melee Attacks"), 0x0A},
	{e("Zerg Missile Attacks"), 0x0B},
	{e("Zerg Flyer Attacks"), 0x0C},
	{e("Protoss Ground Weapons"), 0x0D},
	{e("Protoss Air Weapons"), 0x0E},
	{e("Protoss Plasma Shields"), 0x0F},
	{e("U-238 Shells (Marine Range)"), 0x10},
	{e("Ion Thrusters (Vulture Speed)"), 0x11},
	{e("Titan Reactor (Science Vessel Energy)"), 0x13},
	{e("Ocular Implants (Ghost Sight)"), 0x14},
	{e("Moebius Reactor (Ghost Energy)"), 0x15},
	{e("Apollo Reactor (Wraith Energy)"), 0x16},
	{e("Colossus Reactor (Battle Cruiser Energy)"), 0x17},
	{e("Ventral Sacs (Overlord Transport)"), 0x18},
	{e("Antennae (Overlord Sight)"), 0x19},
	{e("Pneumatized Carapace (Overlord Speed)"), 0x1A},
	{e("Metabolic Boost (Zergling Speed)"), 0x1B},
	{e("Adrenal Glands (Zergling Attack)"), 0x1C},
	{e("Muscular Augments (Hydralisk Speed)"), 0x1D},
	{e("Grooved Spines (Hydralisk Range)"), 0x1E},
	{e("Gamete Meiosis (Queen Energy)"), 0x1F},
	{e("Defiler Energy"), 0x20},
	{e("Singularity Charge (Dragoon Range)"), 0x21},
	{e("Leg Enhancement (Zealot Speed)"), 0x22},
	{e("Scarab Damage"), 0x23},
	{e("Reaver Capacity"), 0x24},
	{e("Gravitic Drive (Shuttle Speed)"), 0x25},
	{e("Sensor Array (Observer Sight)"), 0x26},
	{e("Gravitic Booster (Observer Speed)"), 0x27},
	{e("Khaydarin Amulet (Templar Energy)"), 0x28},
	{e("Apial Sensors (Scout Sight)"), 0x29},
	{e("Gravitic Thrusters (Scout Speed)"), 0x2A},
	{e("Carrier Capacity"), 0x2B},
	{e("Khaydarin Core (Arbiter Energy)"), 0x2C},
	{e("Argus Jewel (Corsair Energy)"), 0x2F},
	{e("Argus Talisman (Dark Archon Energy)"), 0x31},
	{e("Caduceus Reactor (Medic Energy)"), 0x33},
	{e("Chitinous Plating (Ultralisk Armor)"), 0x34},
	{e("Anabolic Synthesis (Ultralisk Speed)"), 0x35},
	{e("Charon Boosters (Goliath Range)"), 0x36},
}

// upgradeIDUpgrade maps from upgrade ID to upgrade.
var upgradeIDUpgrade = map[byte]*Upgrade{}

func init() {
	for _, u := range Upgrades {
		upgradeIDUpgrade[u.ID] = u
	}
}

// UpgradeByID returns the Upgrade for a given ID.
// A new Upgrade with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func UpgradeByID(ID byte) *Upgrade {
	if u := upgradeIDUpgrade[ID]; u != nil {
		return u
	}
	return &Upgrade{repcore.UnknownEnum(ID), ID}
}
