// This file contains the Replay type and its components which model a complete
// SC:BW replay.

package rep

import "github.com/icza/screp/rep/repcmd"

// Replay models an SC:BW replay.
type Replay struct {
	// Header of the replay
	Header *Header

	// Commands of the players
	Commands *Commands

	// MapData describes the map and objects on it
	MapData *MapData

	// Computed contains data that is computed / derived from other parts of the
	// replay.
	Computed *Computed
}

// Compute creates and computes the Computed field.
func (r *Replay) Compute() {
	if r.Computed != nil {
		return
	}

	c := &Computed{}

	if r.Commands != nil {
		for _, cmd := range r.Commands.Cmds {
			switch x := cmd.(type) {
			case *repcmd.LeaveGameCmd:
				c.LeaveGameCmds = append(c.LeaveGameCmds, x)
			case *repcmd.ChatCmd:
				c.ChatCmds = append(c.ChatCmds, x)
			}
		}
	}

	r.Computed = c
}
