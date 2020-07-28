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

	// ChatCmds is a collection of the received chat messages.
	ChatCmds []*repcmd.ChatCmd

	// WinnerTeam if can be detected by the "largest remaining team wins"
	// algorithm. It's 0 if winner team is unknown.
	WinnerTeam byte

	// PlayerDescs contains player descriptions in team order.
	PlayerDescs []*PlayerDesc

	// PIDPlayerDescs maps from player ID to PlayerDesc.
	// Note: all computer players have ID=255, so this won't be accurate for
	// computer players.
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

	// EffectiveCmdCount is the number of effective commands of the player.
	EffectiveCmdCount uint32

	// EAPM is the EAPM (Effective Actions Per Minute) of the player.
	EAPM int32

	// StartLocation of the player
	StartLocation *repcore.Point

	// StartDirection is the direction of the start location of the player
	// compared to the center of the map, expressed using the clock,
	// e.g. 1 o'clock, 6 o'clock etc.
	StartDirection int32

	// Observer tells if the player only observes the game and should be excluded
	// from matchup.
	Observer bool
}

// Redundancy returns the redundancy percent of the player's commands.
// A command is redundant if its ineffective.
func (pd *PlayerDesc) Redundancy() int {
	if pd.CmdCount == 0 {
		return 0
	}
	return int(float64(pd.CmdCount-pd.EffectiveCmdCount)*100/float64(pd.CmdCount) + 0.5)
}
