// This file contains techs.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Tech describes the tech (research).
type Tech struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// Techs is an enumeration of the possible techs.
var Techs = []*Tech{
	{e("Stim Packs"), 0x00},
	{e("Lockdown"), 0x01},
	{e("EMP Shockwave"), 0x02},
	{e("Spider Mines"), 0x03},
	{e("Scanner Sweep"), 0x04},
	{e("Tank Siege Mode"), 0x05},
	{e("Defensive Matrix"), 0x06},
	{e("Irradiate"), 0x07},
	{e("Yamato Gun"), 0x08},
	{e("Cloaking Field"), 0x09},
	{e("Personnel Cloaking"), 0x0a},
	{e("Burrowing"), 0x0b},
	{e("Infestation"), 0x0c},
	{e("Spawn Broodlings"), 0x0d},
	{e("Dark Swarm"), 0x0e},
	{e("Plague"), 0x0f},
	{e("Consume"), 0x10},
	{e("Ensnare"), 0x11},
	{e("Parasite"), 0x12},
	{e("Psionic Storm"), 0x13},
	{e("Hallucination"), 0x14},
	{e("Recall"), 0x15},
	{e("Stasis Field"), 0x16},
	{e("Archon Warp"), 0x17},
	{e("Restoration"), 0x18},
	{e("Disruption Web"), 0x19},
	{e("Unused 26"), 0x1a},
	{e("Mind Control"), 0x1b},
	{e("Dark Archon Meld"), 0x1c},
	{e("Feedback"), 0x1d},
	{e("Optical Flare"), 0x1e},
	{e("Maelstrom"), 0x1f},
	{e("Lurker Aspect"), 0x20},
	{e("Unused 33"), 0x21},
	{e("Healing"), 0x22},
}

// TechByID returns the Tech for a given ID.
// A new Tech with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func TechByID(ID byte) *Tech {
	if int(ID) < len(Techs) {
		return Techs[ID]
	}
	return &Tech{repcore.UnknownEnum(ID), ID}
}
