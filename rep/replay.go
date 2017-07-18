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

	c := &Computed{
		PlayerDescs: make([]*PlayerDesc, len(r.Header.Players)),
	}

	for i, p := range r.Header.Players {
		c.PlayerDescs[i] = &PlayerDesc{
			PlayerID: p.ID,
		}
	}

	// For winners detection, keep track of team sizes:
	teamSizes := map[byte]int{}
	for _, p := range r.Header.Players {
		teamSizes[p.Team]++
	}

	if r.Commands != nil {
		for _, cmd := range r.Commands.Cmds {
			switch x := cmd.(type) {
			case *repcmd.LeaveGameCmd:
				c.LeaveGameCmds = append(c.LeaveGameCmds, x)
				teamSizes[r.Header.PIDPlayers[x.PlayerID].ID]--
			case *repcmd.ChatCmd:
				c.ChatCmds = append(c.ChatCmds, x)
			}
		}
	}

	// Complete winners detection: largest remaining team wins:
	maxTeam, maxSize := byte(0), -1
	for team, size := range teamSizes {
		if size > maxSize {
			maxTeam, maxSize = team, size
		}
	}
	// Are winners detectable?
	if maxSize > 0 {
		// Are there only one team with max size?
		count := 0
		for _, size := range teamSizes {
			if size == maxSize {
				count++
			}
		}
		if count == 1 {
			// We have our winners!
			c.WinnerTeam = maxTeam
		}
	}

	r.Computed = c
}
