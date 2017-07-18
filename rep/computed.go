// This file contains the types describing the computed / derived data.

package rep

import "github.com/icza/screp/rep/repcmd"

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

	// LastCmd is the last command of the player.
	LastCmd repcmd.Cmd

	// TODO also count and exclude commands in the first 2 minutes?

	// CmdCount is the number of commands of the player.
	CmdCount int

	// APM is the APM (Actions Per Minute) of the player.
	APM int
}
