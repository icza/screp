// This file contains the types describing the computed / derived data.

package rep

import "github.com/icza/screp/rep/repcmd"

// Computed contains computed, derived data from other parts of the replay.
type Computed struct {
	// LeaveGameCmds of the players.
	LeaveGameCmds []*repcmd.LeaveGameCmd

	// ChatCmds is a collection of the player chat.
	ChatCmds []*repcmd.ChatCmd

	// TODO ideas: "guess" player results (on a largest remaining team wins basis)
	// APM of players

	// WinnerTeam if can be detected by the "largest remaining team wins"
	// algorithm. It's 0 if winner team is unknown.
	WinnerTeam byte

	// PlayerDescs contains player descriptions in team order.
	PlayerDescs []*PlayerDesc
}

// PlayerDesc contains computed / derived data for a player.
type PlayerDesc struct {
	// PlayerID this PlayerDesc belongs to.
	PlayerID byte
}
