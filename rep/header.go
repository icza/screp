// This file contains the types describing the replay header.

package rep

import (
	"fmt"
	"time"

	"github.com/icza/screp/rep/repcore"
)

// Header models the replay header.
type Header struct {
	// Engine used to play the game and save the replay
	Engine *repcore.Engine

	// Frames is the number of frames. There are approximately ~23.81 frames in
	// a second. (1 frame = 0.042 second to be exact).
	Frames repcore.Frame

	// StartTime is the timestamp when the game started
	StartTime time.Time

	// Title is the game name / title
	Title string

	// Size of the map
	MapWidth, MapHeight uint16

	// AvailSlotsCount is the number of available slots
	AvailSlotsCount byte

	// Speed is the game speed
	Speed *repcore.Speed

	// Type is the game type
	Type *repcore.GameType

	// SubType indicates the size of the "Home" team.
	// For example, in case of 3v5 this is 3, in case of 7v1 this is 7.
	SubType uint16

	// Host is the game creator's name.
	Host string

	// Map name
	Map string

	// Slots contains all players of the game (including open/closed slots)
	Slots []*Player `json:"-"`

	// Players contains the actual ("real") players of the game
	Players []*Player
}

// Duration returns the game duration.
func (h *Header) Duration() time.Duration {
	return h.Frames.Duration()
}

// MapSize returns the map size in widthxheight format, e.g. "64x64".
func (h *Header) MapSize() string {
	return fmt.Sprint(h.MapWidth, "x", h.MapHeight)
}

// Player represents a player of the game.
type Player struct {
	// SlotID is the slot ID
	SlotID uint16

	// ID of the player
	ID byte

	// Type is the player type
	Type *repcore.PlayerType

	// Race of the player
	Race *repcore.Race

	// Team of the player
	Team byte

	// Name of the player
	Name string

	// Color of the player
	Color *repcore.Color
}
