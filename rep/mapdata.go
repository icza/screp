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

	// Tiles is the tile data of the map (within the tile set): width x height elements.
	// 1 Tile is 32 units (pixel)
	Tiles []uint16

	// Mineral field locations on the map
	MineralFields []repcore.Point

	// Geyser locations on the map
	Geysers []repcore.Point

	// StartLocations on the map
	StartLocations []StartLocation

	// Debug holds optional debug info.
	Debug *MapDataDebug `json:"-"`
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
