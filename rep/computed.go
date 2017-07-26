// This file contains the types describing the computed / derived data.

package rep

import (
	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/rep/repcore"
)

// Computed contains computed, derived data from other parts of the replay.
type Computed struct {
	// LeaveGameCmds of the players.
	LeaveGameCmds []*repcmd.LeaveGameCmd

	// ChatCmds is a collection of the player chat.
	ChatCmds []*repcmd.ChatCmd

	// WinnerTeam if can be detected by the "largest remaining team wins"
	// algorithm. It's 0 if winner team is unknown.
	WinnerTeam byte

	// PlayerDescs contains player descriptions in team order.
	PlayerDescs []*PlayerDesc

	// PIDPlayerDescs maps from player ID to PlayerDesc.
	PIDPlayerDescs map[byte]*PlayerDesc `json:"-"`
}

// PlayerDesc contains computed / derived data for a player.
type PlayerDesc struct {
	// PlayerID this PlayerDesc belongs to.
	PlayerID byte

	// LastCmdFrame is the frame of the last command of the player.
	LastCmdFrame repcore.Frame

	// CmdCount is the number of commands of the player.
	CmdCount uint32

	// APM is the APM (Actions Per Minute) of the player.
	APM int32

	// StartLocation of the player
	StartLocation *repcore.Point

	// StartDirection is the direction of the start location of the player
	// compared to the center of the map, expressed using the clock,
	// e.g. 1 o'clock, 6 o'clock etc.
	StartDirection int32
}
