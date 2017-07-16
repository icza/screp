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

	// PlayerDescs contains player descriptions in team order.
	PlayerDescs []*PlayerDesc
}

// PlayerDesc contains computed / derived data for a player.
type PlayerDesc struct {
	// Result of the player.
	Result *Result
}

// Result describes the result of a player (e.g. win or loss).
type Result struct {
	repcore.Enum

	// ID is the (arbitrary) ID of the result.
	ID byte
}

// Results is an enumeration of the possible results.
var Results = []*Result{
	{e("Unknown"), 0x00},
	{e("Win"), 0x01},
	{e("Loss"), 0x02},
	{e("Draw"), 0x03},
}

// Named results
var (
	ResultUnknown = Results[0]
	ResultWin     = Results[1]
	ResultLoss    = Results[2]
	ResultDraw    = Results[3]
)

// e creates a new Enum value.
func e(name string) repcore.Enum {
	return repcore.Enum{Name: name}
}
