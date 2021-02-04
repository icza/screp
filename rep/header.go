// This file contains the types describing the replay header.

package rep

import (
	"fmt"
	"strings"
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

	// RawTitle is the undecoded Title data. It may differ from Title if the latter is invalid UTF-8.
	RawTitle string `json:"-"`

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

	// RawHost is the undecoded Host data. It may differ from Host if the latter is invalid UTF-8.
	RawHost string `json:"-"`

	// Map name
	Map string

	// RawMap is the undecoded Map data. It may differ from Map if the latter is invalid UTF-8.
	RawMap string `json:"-"`

	// Slots contains all players of the game (including open/closed slots)
	Slots []*Player `json:"-"`

	// OrigPlayers contains the actual ("real") players of the game
	// in the order recorded in the replay.
	OrigPlayers []*Player `json:"-"`

	// Players contains the actual ("real") players of the game
	// in team order.
	Players []*Player

	// PIDPlayers maps from player ID to Player.
	// Note: all computer players have ID=255, so this won't be accurate for
	// computer players.
	PIDPlayers map[byte]*Player `json:"-"`

	// Debug holds optional debug info.
	Debug *HeaderDebug `json:"-"`
}

// Duration returns the game duration.
func (h *Header) Duration() time.Duration {
	return h.Frames.Duration()
}

// MapSize returns the map size in widthxheight format, e.g. "64x64".
func (h *Header) MapSize() string {
	return fmt.Sprint(h.MapWidth, "x", h.MapHeight)
}

// Matchup returns the matchup, the race letters of players in team order,
// inserting 'v' between different teams, e.g. "PvT" or "PTZvZTP".
// Observers are excluded from the matchup.
func (h *Header) Matchup() string {
	m := make([]rune, 0, 9)
	first, prevTeam := true, byte(0)
	for _, p := range h.Players {
		if p.Observer {
			continue
		}
		if !first && p.Team != prevTeam {
			m = append(m, 'v')
		}
		m = append(m, p.Race.Letter)
		first, prevTeam = false, p.Team
	}
	return string(m)
}

// PlayerNames returns a comma separated list of player names in team order,
// inserting " VS " between different teams.
func (h *Header) PlayerNames() string {
	buf := &strings.Builder{}
	var prevTeam byte
	for i, p := range h.Players {
		if i > 0 {
			if p.Team != prevTeam {
				buf.WriteString(" VS ")
			} else {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(p.Name)
		prevTeam = p.Team
	}
	return buf.String()
}

// Player represents a player of the game.
type Player struct {
	// SlotID is the slot ID
	SlotID uint16

	// ID of the player.
	// Computer players all have ID=255.
	ID byte

	// Type is the player type
	Type *repcore.PlayerType

	// Race of the player
	Race *repcore.Race

	// Team of the player
	Team byte

	// Name of the player
	Name string

	// RawName is the undecoded Name data. It may differ from Name if the latter is invalid UTF-8.
	RawName string `json:"-"`

	// Color of the player
	Color *repcore.Color

	// Observer tells if the player only observes the game and should be excluded
	// from matchup.
	// This is not stored in replays, this is a calculated property.
	Observer bool
}

// HeaderDebug holds debug info for the header section.
type HeaderDebug struct {
	// Data is the raw, uncompressed data of the section.
	Data []byte

	// Descriptor fields of the data
	Fields []*DebugFieldDescriptor
}

// DebugFieldDescriptor describes some arbitrary data in a byte slice.
type DebugFieldDescriptor struct {
	Offset int    // Offset of the data field
	Length int    // Length of the data field in bytes
	Name   string // Name of the data field
}
