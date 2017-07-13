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

	// TODO ideas: "guess" player results (on a largest remaining team wins basis)
	// APM of players
}

// Result describes the result of a player (e.g. win or loss).
type Result struct {
	repcore.Enum
}

// Results is an enumeration of the possible results.
var Results = []*Result{
	{e("Win")},
	{e("Loss")},
	{e("Draw")},
	{e("Unknown")},
}

// e creates a new Enum value.
func e(name string) repcore.Enum {
	return repcore.Enum{Name: name}
}
