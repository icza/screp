// This file contains the Replay type and its components which model a complete
// SC:BW replay.

package rep

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/icza/gox/stringsx"
	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/rep/repcore"
	"github.com/icza/screp/repparser/repdecoder"
)

// Replay models an SC:BW replay.
type Replay struct {
	// Stored here for decoding purposes only.
	RepFormat repdecoder.RepFormat `json:"-"`

	// Header of the replay
	Header *Header

	// Commands of the players
	Commands *Commands

	// MapData describes the map and objects on it
	MapData *MapData

	// Computed contains data that is computed / derived from other parts of the
	// replay.
	Computed *Computed

	// ShieldBattery holds info if game was played on ShieldBattery
	ShieldBattery *ShieldBattery `json:",omitempty"`
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
			mapName := r.Header.Map
			if r.MapData != nil {
				mapName = r.MapData.Name
			}
			// counter-examples: " \aai \x04hunters \x02remastered \x062.0", "\x03(XB2)\x06 Big Game Hunters"
			mapName = strings.ToLower(stringsx.Clean(mapName))
			// "[ai]" maps are special, we can do better than in general:
			switch {
			case mapName == "  hunters kespa soulclan ai" || mapName == ":da hunters ai" ||
				mapName == "(xb2) big game hunters" || strings.HasPrefix(mapName, "王牌猎人") ||
				strings.Contains(mapName, "[ai]") || strings.Contains(mapName, "ai hunters") || strings.Contains(mapName, "bgh random teams"):
				r.detectObservers(pidBuilds, obsProfileUMSAI)
				r.computeUMSTeamsAI()

			default:
				r.computeUMSTeams()
			}
		case repcore.GameTypeMelee:
			r.detectObservers(pidBuilds, obsProfileMelee)
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
//
//	-there are only 2 human players on team 1, having train or build commands
//	-all other players are on a different team, and they have no train nor build commands
//
// If this case is detected, the players on team 1 are split into team 1 and 2,
// and all players (observers) on the (original) team 2 are assiged to team 3, and marked as observers.
func (r *Replay) computeUMSTeams() {
	// We'll have to check player commands later, so if it's not parsed, don't waste any time:
	if r.Commands == nil {
		return
	}

	players := r.Header.Players
	if len(players) < 2 {
		return
	}

	playerCandidateIDs, obsCandidateIDs := map[byte]bool{}, map[byte]bool{}

	for i, p := range players {
		if p.Type != repcore.PlayerTypeHuman {
			return // Non-human involved, don't get involved!
		}
		if i < 2 { // candidates for 1v1 players
			if p.Team != 1 {
				return
			}
			playerCandidateIDs[p.ID] = true
		} else { // candidates for observers
			if p.Team == 1 {
				return
			}
			obsCandidateIDs[p.ID] = true
		}
	}

	// Check if player candidates have train or build commands, and obs candidates don't.
	playerTrainBuildCount := 0
	noObsCandidates := len(obsCandidateIDs) == 0

cmdLoop:
	for _, cmd := range r.Commands.Cmds {
		switch cmd.(type) {
		case *repcmd.TrainCmd, *repcmd.BuildCmd:
			if playerCandidateIDs[cmd.BaseCmd().PlayerID] {
				playerTrainBuildCount++
				if noObsCandidates {
					break cmdLoop // We got what we want, no obs candidates, no need to continue
				}
			} else if obsCandidateIDs[cmd.BaseCmd().PlayerID] {
				return // An obs candidate have a train or build command, this is not the special case we're looking for
			}
		}
	}

	if playerTrainBuildCount == 0 {
		return // Player candidates have no train nor build commands, this is not the special case we're looking for
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

// computeUMSTeamsAI computes the teams in UMS AI games.
//
// Maps having "[AI]" in their name are special: they create random teams after start,
// with optional observers. Random team arragement usually happens 18 seconds after game start.
// Commands selecting teams are not recorded, but since teams are created randomly, players very often check alliance
// to see who their allies are. This reasults in Alliance commands recording the team setup at the time
// of the checks. We will use these to detect teams.
// There is no guarantee players check alliance and they may also change the (initial) alliance arranged by the map.
//
// Different AI maps handle observers differently.
// Some set alliance from players to observers too (observers will be in team 1 and team 2 too), and observers are allied with each other only.
// Other AI maps set alliance only between team members and observers separately.
// Observers may also be allied with all players / slots.
//
// Alliance commands in the first 115 seconds are checked; if they consistently denote 2 teams (and an optional observer team),
// players are assigned team 1 and 2 respectively (and observers are assigned team 3, and marked as observers).
//
// If teams can be computed, also rearranges Header.Players and Computed.PlayerDescs
// according to new teams.
func (r *Replay) computeUMSTeamsAI() {
	// We'll have to check player commands later, so if it's not parsed, don't waste any time:
	if r.Commands == nil {
		return
	}

	players := r.Header.Players
	if len(players) < 2 {
		return
	}

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

	// Set of slotIDs belonging to players (non-observers):
	playerSlotIDs := map[byte]bool{}
	for _, p := range players {
		if !p.Observer {
			playerSlotIDs[byte(p.SlotID)] = true
		}
	}
	filterOutObserverSlotIDs := func(slotIDs []byte) (result []byte) {
		for _, slotID := range slotIDs {
			if playerSlotIDs[slotID] {
				result = append(result, slotID)
			}
		}
		return
	}

	// Slot IDs of player's last Alliance commands, observers filtered out:
	pidSlotIDs := map[byte][]byte{}

	// If there are 2 players only, it's unlikely they'll check alliances and thus below team detection would fail.
	// To make it still work, initialize with self-alliance:
	if len(playerSlotIDs) == 2 {
		for _, p := range players {
			if !p.Observer {
				pidSlotIDs[p.ID] = []byte{byte(p.SlotID)}
			}
		}
	}

	// Stop after ~115 seconds: use the "initial" teams
	frameMaxLimit := repcore.Duration2Frame(115 * time.Second)
	frameMinLimit := repcore.Duration2Frame(18 * time.Second)
	for _, cmd := range r.Commands.Cmds {
		if cmd.BaseCmd().Frame > frameMaxLimit {
			break
		}
		if ac, ok := cmd.(*repcmd.AllianceCmd); ok {
			if p := r.Header.PIDPlayers[ac.PlayerID]; p != nil && p.Observer {
				continue
			}
			filteredSlotIDs := filterOutObserverSlotIDs(ac.SlotIDs) // Note: first filter because on "BGH Random Teams" this also includes the obs computer!
			if len(filteredSlotIDs) == 1 && cmd.BaseCmd().Frame < frameMinLimit {
				continue // Random team arrangement has likely not done, do not count!
			}
			pidSlotIDs[ac.PlayerID] = filteredSlotIDs
		}
	}

	// Since observers are filtered out, there should be exactly 2 teams, with equal size,
	// disjunct players. And all other players must be observers.

	// We use the string representation of the slots as the virtual team ID
	// which will be something like "[0 2 3]"
	virtualTeamIDSlotIDs := map[string][]byte{}
	for _, slotIDs := range pidSlotIDs {
		if len(slotIDs) == 0 {
			continue
		}
		virtualID := fmt.Sprint(slotIDs)
		virtualTeamIDSlotIDs[virtualID] = slotIDs
	}
	if len(virtualTeamIDSlotIDs) != 2 {
		return // not 2 teams exactly
	}

	var team1SlotIDs, team2SlotIDs []byte
	for _, slotIDs := range virtualTeamIDSlotIDs {
		if team1SlotIDs == nil {
			team1SlotIDs = slotIDs
		} else {
			team2SlotIDs = slotIDs
		}
	}
	// Use consistent team order (order by first slot ID):
	if team2SlotIDs[0] < team1SlotIDs[0] {
		team1SlotIDs, team2SlotIDs = team2SlotIDs, team1SlotIDs
	}
	// Check if teams are disjuct:
	for _, slotIDA := range team1SlotIDs {
		if bytes.IndexByte(team2SlotIDs, slotIDA) >= 0 {
			return // slotIDA is in both teams
		}
	}
	// Check if all non-observers are in one of the 2 teams:
	slotIDTeams := map[byte]byte{}
	for _, slotID := range team1SlotIDs {
		slotIDTeams[slotID] = 1
	}
	for _, slotID := range team2SlotIDs {
		slotIDTeams[slotID] = 2
	}
	if len(playerSlotIDs) != len(slotIDTeams) {
		return // Not all player assigned to team 1 or 2
	}

	// Assign new teams
	for _, p := range players {
		if p.Observer {
			p.Team = 3
		} else {
			p.Team = slotIDTeams[byte(p.SlotID)]
		}
	}

	// Re-sort Header.Players and Computed.PlayerDescs
	r.rearrangePlayers()
}

// obsProfile holds data for observer rules in different scenarios.
type obsProfile struct {
	apmLimit        int32         // Human obs must be below this APM limit
	buildCmdsLimit  int           // Human obs must be below this build command limit
	earlyLeaveFrame repcore.Frame // consider early leavers as observer
	computer        bool          // Classify computer as observer (BGH Random Teams map)
}

var (
	obsProfileMelee = &obsProfile{apmLimit: 25, buildCmdsLimit: 5}
	obsProfileUMSAI = &obsProfile{apmLimit: 40, buildCmdsLimit: 2, earlyLeaveFrame: repcore.Duration2Frame(18 * time.Second), computer: true}
)

// detectObservers detects observers based on the given obs profile.
func (r *Replay) detectObservers(pidBuilds map[byte]int, obsProf *obsProfile) {
	c := r.Computed

	// Criteria for observers:
	//   - Human
	//       and
	//         - APM < obsProf.apmLimit
	//         - Has less than obsProf.buildCmdsLimit build commands
	//       or
	//         - obsProf.earlyLeaveFrame is not zero and the player left earlier

	numObs := 0
	for i, p := range r.Header.Players {
		if p.Type == repcore.PlayerTypeHuman &&
			(c.PlayerDescs[i].APM < obsProf.apmLimit && pidBuilds[p.ID] < obsProf.buildCmdsLimit ||
				obsProf.earlyLeaveFrame > 0 && c.PlayerDescs[i].LastCmdFrame < obsProf.earlyLeaveFrame) ||
			(obsProf.computer && p.Type == repcore.PlayerTypeComputer) {
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
	// We'll have to check player commands later, so if it's not parsed, don't waste any time:
	if r.Commands == nil {
		return
	}

	players := r.Header.Players
	if len(players) < 2 {
		return
	}

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
	r.rearrangePlayers()
}

// rearrangePlayers rearranges Header.Players and Computed.PlayerDescs to be in "team order".
// Teams may be assigned / changed by team detection algorithms, this helper function
// rearranges the players so the order will be in team-order.
func (r *Replay) rearrangePlayers() {
	players := r.Header.Players
	pds := r.Computed.PlayerDescs

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
	// The essence of this is to process Leave game commands and track remaining team sizes.
	// Problems:
	//   -Leave game commands are not recorded for computers
	//   -Leave game commands are not recorded for the replay saver

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
			// Add virutal leave game cmd
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
//   - PI/2 => 12 (o'clock)
//   - 0 => 3 (o'clock)
//   - PI => 9 (o'clock)
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
