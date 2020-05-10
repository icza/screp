// This file contains the Replay type and its components which model a complete
// SC:BW replay.

package rep

import (
	"math"

	"github.com/icza/screp/rep/repcmd"
)

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
		PlayerDescs:    make([]*PlayerDesc, len(r.Header.Players)),
		PIDPlayerDescs: make(map[byte]*PlayerDesc, len(r.Header.Players)),
	}

	for i, p := range r.Header.Players {
		pd := &PlayerDesc{
			PlayerID: p.ID,
		}
		c.PlayerDescs[i] = pd
		c.PIDPlayerDescs[p.ID] = pd
	}

	// For winners detection, keep track of team sizes:
	teamSizes := map[byte]int{}
	for _, p := range r.Header.Players {
		teamSizes[p.Team]++
	}

	if r.Commands != nil {
		cmds := r.Commands.Cmds
		for _, cmd := range cmds {
			// Observers' commands (e.g. chat) have PlayerID starting with 128 (2nd obs 129 etc.)
			// We don't have PlayerDescs for them, so must check:
			if pd := c.PIDPlayerDescs[cmd.BaseCmd().PlayerID]; pd != nil {
				c.PIDPlayerDescs[cmd.BaseCmd().PlayerID].CmdCount++
			}
			switch x := cmd.(type) {
			case *repcmd.LeaveGameCmd:
				c.LeaveGameCmds = append(c.LeaveGameCmds, x)
				if pid := r.Header.PIDPlayers[x.PlayerID]; pid != nil {
					teamSizes[pid.Team]--
				}
			case *repcmd.ChatCmd:
				c.ChatCmds = append(c.ChatCmds, x)
			}
		}

		// Search for last commands:
		// Make a local copy of the PIDPlayerDescs map to keep track of
		// players we still need this info for:
		pidPlayerDescs := make(map[byte]*PlayerDesc, len(r.Header.Players))
		for pid, pd := range c.PIDPlayerDescs {
			// Optimization: Only include players that do have commands:
			if pd.CmdCount > 0 {
				pidPlayerDescs[pid] = pd
			}
		}
		for i := len(cmds) - 1; i >= 0; i-- {
			cmd := cmds[i]
			baseCmd := cmd.BaseCmd()
			pd := pidPlayerDescs[baseCmd.PlayerID]
			if pd == nil {
				continue
			}
			if baseCmd.Frame > r.Header.Frames {
				// Bad parsing or corrupted replay may result in invalid frames,
				// do not use such a bad frame.
				continue
			}
			pd.LastCmdFrame = baseCmd.Frame
			// Optimization: If this was the last player, break:
			if len(pidPlayerDescs) == 1 {
				break
			}
			delete(pidPlayerDescs, pd.PlayerID)
		}

		// Complete winners detection: largest remaining team wins
		// (if there were Leave game commands)
		if len(c.LeaveGameCmds) > 0 {
			maxTeam, maxSize := byte(0), -1
			for team, size := range teamSizes {
				if size > maxSize {
					maxTeam, maxSize = team, size
				}
			}
			// Are winners detectable?
			if maxSize > 0 {
				// Is there only one team with max size?
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
		}
		// If winners are not detectable and there are no Leave game commands,
		// it's most likely due to the replay saver left first.
		// Replay saver is the one who receives the chat messages.
		// If we have chat commands, declare the replay saver's team the loser.
		// (Note chat is saved since patch 1.16, released on 2008-11-25.)
		// If there is only one team besides the loser (2 teams altogether), we have our winner.
		if c.WinnerTeam == 0 && len(c.LeaveGameCmds) == 0 && len(c.ChatCmds) > 0 && len(teamSizes) == 2 {
			// rep saver might be an observer, so must check if there's a player for him/her:
			if repSaver := r.Header.PIDPlayers[c.ChatCmds[0].PlayerID]; repSaver != nil {
				loserTeam := repSaver.Team
				for team := range teamSizes {
					if team != loserTeam {
						c.WinnerTeam = team
						break
					}
				}
			}
		}
		// Also if there are 2 players and 2 Game leave commands,
		// and they are on different teams, declare the 2nd leaver the winner
		// (this might be the case if an obs saved the replay).
		if c.WinnerTeam == 0 && len(r.Header.Players) == 2 && len(c.LeaveGameCmds) == 2 {
			p1 := r.Header.PIDPlayers[c.LeaveGameCmds[0].PlayerID]
			p2 := r.Header.PIDPlayers[c.LeaveGameCmds[1].PlayerID]
			if p1 != nil && p2 != nil && p1.Team != p2.Team {
				c.WinnerTeam = p2.Team
			}
		}

		// Calculate APMs:
		for _, pd := range c.PlayerDescs {
			if pd.LastCmdFrame == 0 {
				continue
			}
			mins := pd.LastCmdFrame.Duration().Minutes()
			pd.APM = int32(float64(pd.CmdCount)/mins + 0.5)
		}
	}

	if r.MapData != nil {
		// 1 tile is 32 pixels, so half is x*16:
		cx, cy := float64(r.Header.MapWidth*16), float64(r.Header.MapHeight*16)
		// Lookup start location of players
		sls := r.MapData.StartLocations
		for i, p := range r.Header.Players {
			for j := range sls {
				if p.SlotID == uint16(sls[j].SlotID) {
					pt := &sls[j].Point
					c.PlayerDescs[i].StartLocation = pt
					// Map Y coordinate grows from top to bottom:
					c.PlayerDescs[i].StartDirection = angleToClock(
						math.Atan2(cy-float64(pt.Y), float64(pt.X)-cx),
					)
					break
				}
			}
		}
	}

	r.Computed = c
}

// angleToClock converts an angle given in radian to an hour clock value
// in the range of 1..12.
//
// Examples:
//  - PI/2 => 12 (o'clock)
//  - 0 => 3 (o'clock)
//  - PI => 9 (o'clock)
func angleToClock(angle float64) int32 {
	// The algorithm below computes clock value in the range of 0..11 where
	// 0 corresponds to 12.

	// 1 hour is PI/6 angle range
	const oneHour = math.Pi / 6

	// Shift by 3:30 (0 or 12 o-clock starts at 11:30)
	// and invert direction (clockwise):
	angle = -angle + oneHour*3.5

	// Put in range of 0..2*PI
	for angle < 0 {
		angle += oneHour * 12
	}
	for angle >= oneHour*12 {
		angle -= oneHour * 12
	}

	// And convert to a clock value:
	hour := int32(angle / oneHour)
	if hour == 0 {
		return 12
	}
	return hour
}
