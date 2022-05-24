// This file contains the Replay type and its components which model a complete
// SC:BW replay.

package rep

import (
	"math"
	"sort"
	"time"

	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/rep/repcore"
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

	players := r.Header.Players
	numPlayers := len(players)

	c := &Computed{
		PlayerDescs:    make([]*PlayerDesc, numPlayers),
		PIDPlayerDescs: make(map[byte]*PlayerDesc, numPlayers),
	}
	r.Computed = c

	for i, p := range players {
		pd := &PlayerDesc{
			PlayerID: p.ID,
		}
		c.PlayerDescs[i] = pd
		c.PIDPlayerDescs[p.ID] = pd
	}

	if r.Commands != nil {
		// We need to gather player's commands separately for EAPM calculation.
		// We could use a map, mapping from pid to player's commands, but then when building it,
		// we would have to always reassign the slice. Instead we use a pointer to a wrapper struct:
		type pidCmdsWrapper struct {
			cmds []repcmd.Cmd
		}
		pidCmdsWrappers := make(map[byte]*pidCmdsWrapper, numPlayers)
		pidBuilds := make(map[byte]int, numPlayers) // Build commands count per player
		for _, p := range players {
			pidCmdsWrappers[p.ID] = &pidCmdsWrapper{
				cmds: make([]repcmd.Cmd, 0, len(r.Commands.Cmds)/numPlayers), // Estimate even cmd distribution for fewer reallocations
			}
		}

		cmds := r.Commands.Cmds
		for _, cmd := range cmds {
			// Observers' commands (e.g. chat) have PlayerID starting with 128 (2nd obs 129 etc.)
			// We don't have PlayerDescs for them, so must check:
			baseCmd := cmd.BaseCmd()
			if pd := c.PIDPlayerDescs[baseCmd.PlayerID]; pd != nil {
				pd.CmdCount++
				pidCmdsWrapper := pidCmdsWrappers[baseCmd.PlayerID]
				pidCmdsWrapper.cmds = append(pidCmdsWrapper.cmds, cmd)
				baseCmd.IneffKind = CmdIneffKind(pidCmdsWrapper.cmds, len(pidCmdsWrapper.cmds)-1)
				if baseCmd.IneffKind.Effective() {
					pd.EffectiveCmdCount++
				}
			}
			switch x := cmd.(type) {
			case *repcmd.LeaveGameCmd:
				c.LeaveGameCmds = append(c.LeaveGameCmds, x)
			case *repcmd.ChatCmd:
				c.ChatCmds = append(c.ChatCmds, x)
			case *repcmd.BuildCmd:
				pidBuilds[baseCmd.PlayerID]++
			}
		}

		// Detect replay saver:
		// Replay saver is the one who receives the chat messages.
		// (Note chat is saved since patch 1.16, released on 2008-11-25.)
		if len(c.ChatCmds) > 0 {
			c.RepSaverPlayerID = &c.ChatCmds[0].PlayerID
		}

		// Search for last commands:
		// Make a local copy of the PIDPlayerDescs map to keep track of
		// players we still need this info for:
		pidPlayerDescs := make(map[byte]*PlayerDesc, numPlayers)
		for pid, pd := range c.PIDPlayerDescs {
			// Only include players that do have commands:
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
			if baseCmd.Frame > r.Header.Frames || baseCmd.Frame < 0 {
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

		// Calculate APMs and EAPMs:
		for _, pd := range c.PlayerDescs {
			if pd.LastCmdFrame == 0 {
				continue
			}
			mins := pd.LastCmdFrame.Duration().Minutes()
			pd.APM = int32(float64(pd.CmdCount)/mins + 0.5)
			pd.EAPM = int32(float64(pd.EffectiveCmdCount)/mins + 0.5)
		}

		switch r.Header.Type {
		case repcore.GameTypeUMS:
			r.computeUMSTeams()
		case repcore.GameTypeMelee:
			r.detectMeleeObservers(pidBuilds)
			r.computeMeleeTeams()
		}

		r.computeWinners()
	}

	if r.MapData != nil {
		// 1 tile is 32 pixels, so half is x*16:
		cx, cy := float64(r.Header.MapWidth*16), float64(r.Header.MapHeight*16)
		// Lookup start location of players
		sls := r.MapData.StartLocations
		for i, p := range players {
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
}

// computeUMSTeams computes the teams in UMS games.
//
// Handles a special case: 1v1 game with observers.
// Rules to detect this case:
//   -there are only 2 human players on team 1
//   -all other players are on team 2, and they have no train nor build commands
//
// If this case is detected, the players on team 1 are split into team 1 and 2,
// and all players (observers) on the (original) team 2 are assiged to team 3, and marked as observers.
func (r *Replay) computeUMSTeams() {
	players := r.Header.Players
	if len(players) < 2 {
		return
	}

	obsCandidateIDs := map[byte]bool{}

	for i, p := range players {
		if p.Type != repcore.PlayerTypeHuman {
			return // Non-human involved, don't get involved!
		}
		switch {
		case i < 2: // candidates for 1v1 players
			if p.Team != 1 {
				return
			}
		default: // candidates for observers
			if p.Team == 1 {
				return
			}
			obsCandidateIDs[p.ID] = true
		}
	}

	// Check if observers have no train or build commands
	if len(players) > 2 {
		if r.Commands == nil {
			return
		}
		for _, cmd := range r.Commands.Cmds {
			switch cmd.(type) {
			case *repcmd.TrainCmd, *repcmd.BuildCmd:
				if obsCandidateIDs[cmd.BaseCmd().PlayerID] {
					return // An obs candidate have a train or build command, this is not the special case we're looking for
				}
			}
		}
	}

	// Special case detected, proceed to re-teaming.

	// 1v1 players
	players[0].Team = 1
	players[1].Team = 2

	// Observers
	for _, p := range players[2:] {
		p.Team = 3
		p.Observer = true
	}
}

// detectMeleeObservers detects observers in Melee games.
func (r *Replay) detectMeleeObservers(pidBuilds map[byte]int) {
	c := r.Computed

	// Criteria for observers:
	//   - Human
	//   - APM < 25
	//   - Has less than 5 build commands

	numObs := 0
	for i, p := range r.Header.Players {
		if p.Type == repcore.PlayerTypeHuman && c.PlayerDescs[i].APM < 25 && pidBuilds[p.ID] < 5 {
			p.Observer = true
			numObs++
		}
	}

	// If less than 2 non-obs players remained, undo:
	if len(r.Header.Players)-numObs < 2 {
		for _, p := range r.Header.Players {
			p.Observer = false
		}
	}
}

// computeMeleeTeams computes the teams in melee games based on player Alliance commands.
//
// If teams can be computed, also rearranges Header.Players and Computed.PlayerDescs
// according to new teams.
func (r *Replay) computeMeleeTeams() {
	players := r.Header.Players
	if len(players) < 2 {
		return
	}

	c := r.Computed
	pds := c.PlayerDescs

	// Only compute if we don't yet have team info (if all teams are the same):
	var nonObsPlayer *Player
	for _, p := range players {
		if p.Observer {
			continue
		}
		if nonObsPlayer == nil {
			nonObsPlayer = p
		} else {
			if p.Team != nonObsPlayer.Team {
				return
			}
		}
	}

	// NOTE: all computers have pid=255, but since they don't set alliance
	// and they can't be allied with, they won't cause trouble.
	// Only when their team is set, don't try set teams of "faulty" teammates.

	pidSlotIDs := map[byte][]byte{}
	// By default all players are allied to themselves only:
	for _, p := range players {
		if p.Observer {
			continue
		}
		pidSlotIDs[p.ID] = []byte{byte(p.SlotID)}
	}

	// Stop after ~90 seconds: use the "initial" teams
	frameLimit := repcore.Duration2Frame(90 * time.Second)
	for _, cmd := range r.Commands.Cmds {
		if cmd.BaseCmd().Frame > frameLimit {
			break
		}
		if ac, ok := cmd.(*repcmd.AllianceCmd); ok {
			if p := r.Header.PIDPlayers[ac.PlayerID]; p != nil && p.Observer {
				continue
			}
			pidSlotIDs[ac.PlayerID] = ac.SlotIDs
		}
	}

	// Check if set alliances are consistent:
	// For each A=>B alliance there must be a B=>A
	// Build maps for fast lookups:
	slotIDPlayers := map[byte]*Player{}
	for _, p := range players {
		if !p.Observer {
			slotIDPlayers[byte(p.SlotID)] = p
		}
	}
	slotIDSlotIDs := map[byte][]byte{}
	for pid, slotIDs := range pidSlotIDs {
		if p := r.Header.PIDPlayers[pid]; p != nil {
			slotIDSlotIDs[byte(p.SlotID)] = slotIDs
		}
	}
	// Now check the consistency:
	for pid, slotIDs := range pidSlotIDs {
		p := r.Header.PIDPlayers[pid]
		if p == nil {
			continue
		}
		slotIDA := byte(p.SlotID)
		for _, slotIDB := range slotIDs {
			if slotIDA == slotIDB {
				continue
			}
			if p := slotIDPlayers[slotIDB]; p == nil || p.Observer {
				continue
			}
			// There is a slotIDA => slotIDB alliance, there must be a slotIDB => slotIDA:
			found := false
			for _, slotIDC := range slotIDSlotIDs[slotIDB] {
				if slotIDC == slotIDA {
					// found!
					found = true
					break
				}
			}
			if !found {
				// Alliance is inconsistent, do not change teams:
				return
			}
		}
	}

	// Found matching alliances! Assign new teams.
	// Start clean:
	for _, p := range players {
		p.Team = 0
	}
	team := byte(1)
	for _, p := range players {
		if p.Observer {
			continue // We handle observers last
		}
		if p.Team != 0 {
			continue // Already assigned
		}
		p.Team = team
		if p.Type != repcore.PlayerTypeComputer { // pidSlotIDs is not valid for computers.
			// All teammates get the same team
			for _, slotID := range pidSlotIDs[p.ID] {
				if p := slotIDPlayers[slotID]; p != nil && !p.Observer {
					p.Team = team
				}
			}
		}
		team++
	}
	// Last assign highest team to observers:
	for _, p := range players {
		if p.Observer {
			p.Team = team
		}
	}

	// Re-sort Header.Players and Computed.PlayerDescs
	type wrapper struct {
		p  *Player
		pd *PlayerDesc
	}
	ws := make([]wrapper, len(players))
	for i, p := range players {
		ws[i] = wrapper{p: p, pd: pds[i]}
	}
	sort.SliceStable(ws, func(i, j int) bool {
		return ws[i].p.Team < ws[j].p.Team
	})
	for i := range ws {
		players[i] = ws[i].p
		pds[i] = ws[i].pd
	}
}

// computeWinners attempts to compute winners using "largest remaining team wins" principle.
func (r *Replay) computeWinners() {
	// Situation: game result (winners / losers) is not recorded in replays.
	// We try to determine the winners based on the "largest remaining team wins" principle.
	// The essence of this is to procedd Leave game commands and track remaining team sizes.
	// Problems:
	//   -Leave game commands are not recorded for computers
	//   -Leave game c ommands are not recorded for the replay saver

	c := r.Computed

	// Keep track of team sizes and computer counts:
	nonObsPlayersCount := 0
	teamSizes := map[byte]int{}      // Excluding computers
	teamCompsCount := map[byte]int{} // Including only computers

	for _, p := range r.Header.Players {
		if !p.Observer {
			if p.Type == repcore.PlayerTypeComputer {
				teamCompsCount[p.Team]++
			} else {
				teamSizes[p.Team]++
			}
			nonObsPlayersCount++
		}
	}

	// If there is a team full of only computers, we can't detect winners.
	for team := range teamCompsCount {
		if teamSizes[team] == 0 {
			return // This team only consists of computers
		}
	}

	// Computers never leave, so use only non-computer sizes (teamSizes) ongoing.

	// Keep only leave game commands of non-observers, which matters if / when we check the last of them.
	leaveGameCmds := make([]*repcmd.LeaveGameCmd, 0, len(c.LeaveGameCmds)+1)
	for _, lgcmd := range c.LeaveGameCmds {
		if p := r.Header.PIDPlayers[lgcmd.PlayerID]; p != nil {
			if !p.Observer {
				leaveGameCmds = append(leaveGameCmds, lgcmd)
			}
		}
	}

	// There is no Leave game command recorded for the replay saver.
	// If we know the replay saver, "simulate" a leave game command
	// for him/her as the last leave game command.
	if c.RepSaverPlayerID != nil {
		// rep saver might be an observer, so must check if there's a player for him/her:
		if repSaver := r.Header.PIDPlayers[*c.RepSaverPlayerID]; repSaver != nil && !repSaver.Observer {
			// Add virutal leave cmd
			leaveGameCmds = append(leaveGameCmds, &repcmd.LeaveGameCmd{
				Base: &repcmd.Base{
					PlayerID: repSaver.ID, // Only PlayerID is needed / used
				},
			})
		}
	}

	for _, lgcmd := range leaveGameCmds {
		// lgcmd.PlayerID exists in PIDPlayers, was checked when assembled leaveGameCmds
		teamSizes[r.Header.PIDPlayers[lgcmd.PlayerID].Team]--
	}

	if len(teamSizes) < 2 || // There are no multiple teams
		len(leaveGameCmds) == 0 { // There were no Leave game commands, not even a "virtual" one,
		// we just don't know who the winners are.
		return
	}

	// Complete winners detection: largest remaining team wins
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
			return
		}
	}

	// There is no single largest team.
	// If there are multiple teams (not just one), and if all (non-obs) players left (we have a leave game command for all),
	// declare the last leaver's team the winner team.
	// Often this happens if an observer saves the replay, and he/she is the one last leaving (there's no leave game command for observers).
	if len(leaveGameCmds) == nonObsPlayersCount {
		playerID := leaveGameCmds[len(leaveGameCmds)-1].PlayerID
		c.WinnerTeam = r.Header.PIDPlayers[playerID].Team
		return
	}
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
