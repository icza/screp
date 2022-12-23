// This file contains the types describing the map data.

package rep

import "github.com/icza/screp/rep/repcore"

// MapData describes the map and objects on it.
type MapData struct {
	// Version of the map.
	// 0x2f: StarCraft beta
	// 0x3b: 1.00-1.03 StarCraft and above ("hybrid")
	// 0x3f: 1.04 StarCraft and above ("hybrid")
	// 0x40: StarCraft Remastered
	// 0xcd: Brood War
	// 0xce: Brood War Remastered
	Version uint16

	// TileSet defines the tile set used on the map.
	TileSet *repcore.TileSet

	// Scenario name
	Name string

	// Scenario description
	Description string

	// PlayerOwners defines the player types (player owners).
	PlayerOwners []*repcore.PlayerOwner

	// PlayerSides defines the player sides (player races).
	PlayerSides []*repcore.PlayerSide

	// Tiles is the tile data of the map (within the tile set): width x height elements.
	// 1 Tile is 32 units (pixel)
	Tiles []uint16 `json:",omitempty"`

	// Mineral field locations on the map
	MineralFields []Resource `json:",omitempty"`

	// Geyser locations on the map
	Geysers []Resource `json:",omitempty"`

	// StartLocations on the map
	StartLocations []StartLocation

	// MapGraphics holds data for map image rendering.
	MapGraphics *MapGraphics `json:",omitempty"`

	// Debug holds optional debug info.
	Debug *MapDataDebug `json:"-"`
}

// MaxHumanPlayers returns the max number of human players on the map.
func (md *MapData) MaxHumanPlayers() (count int) {
	for _, owner := range md.PlayerOwners {
		if owner == repcore.PlayerOwnerHumanOpenSlot {
			count++
		}
	}
	return
}

// Resource describes a resource (mineral field of vespene geyser).
type Resource struct {
	// Location of the resource
	repcore.Point

	// Amount of the resource
	Amount uint32
}

// StartLocation describes a player start location on the map
type StartLocation struct {
	repcore.Point

	// SlotID of the owner of this start location;
	// Belongs to the Player with matching Player.SlotID
	SlotID byte
}

// MapDataDebug holds debug info for the map data section.
type MapDataDebug struct {
	// Data is the raw, uncompressed data of the section.
	Data []byte
}

// MapGraphics holds info usually required only for map image rendering.
type MapGraphics struct {
	// PlacedUnits contains all placed units on the map.
	// This also includes mineral fields, geysers, startlocations.
	PlacedUnits []*PlacedUnit

	// Sprites contains additional visual sprites on the map.
	Sprites []*Sprite
}

type PlacedUnit struct {
	repcore.Point

	// UnitID is the unit id. This value is used in repcmd.Unit.UnitID.
	UnitID uint16

	// SlotID of the owner of this unit.
	// Belongs to the Player with matching Player.SlotID
	SlotID byte

	// ResourceAmount of if it's a resource
	ResourceAmount uint32 `json:",omitempty"`

	// Sprite tells if this unit is a sprite.
	Sprite bool `json:",omitempty"`
}

type Sprite struct {
	repcore.Point

	// SpriteID is the sprite id.
	SpriteID uint16
}
