// This file contains the algorithm implementation for EAPM classification and calculation.

package rep

import (
	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/rep/repcore"
)

const (
	// EAPMVersion is a Semver2 compatible version of the EAPM algorithm.
	EAPMVersion = "v1.0.3"
)

// IsCmdEffective tells if a command is considered effective so it can be included in EAPM calculation.
//
// cmds must contain commands of the cmd's player only. It may be a partially filled slice, but must contain
// the player's all commands up to the command in question: len(cmds) > i must hold.
func IsCmdEffective(cmds []repcmd.Cmd, i int) bool {
	return CmdIneffKind(cmds, i) == repcore.IneffKindEffective
}

// CmdIneffKind returns the IneffKind classification of the given command.
//
// cmds must contain commands of the cmd's player only. It may be a partially filled slice, but must contain
// the player's all commands up to the command in question: len(cmds) > i must hold.
func CmdIneffKind(cmds []repcmd.Cmd, i int) repcore.IneffKind {
	if i == 0 {
		return repcore.IneffKindEffective // First command is effective whatever it is
	}

	// Try to "prove" command is ineffective. If we can't, it's effective.

	cmd := cmds[i]
	tid := cmd.BaseCmd().Type.ID

	// Unit queue overflow
	switch tid {
	case repcmd.TypeIDTrain, repcmd.TypeIDTrainFighter, repcmd.TypeIDCancelTrain:
		if countSameCmds(cmds, i, cmd) >= 6 {
			return repcore.IneffKindUnitQueueOverflow
		}
	}

	prevCmd := cmds[i-1] // i > 0
	prevTid := prevCmd.BaseCmd().Type.ID

	deltaFrame := cmd.BaseCmd().Frame - prevCmd.BaseCmd().Frame

	// Too fast cancel
	if deltaFrame <= 20 {
		switch {
		case (prevTid == repcmd.TypeIDTrain || prevTid == repcmd.TypeIDTrainFighter) && tid == repcmd.TypeIDCancelTrain:
			return repcore.IneffKindFastCancel
		case (prevTid == repcmd.TypeIDUnitMorph || prevTid == repcmd.TypeIDBuildingMorph) && tid == repcmd.TypeIDCancelMorph:
			return repcore.IneffKindFastCancel
		case prevTid == repcmd.TypeIDUpgrade && tid == repcmd.TypeIDCancelUpgrade:
			return repcore.IneffKindFastCancel
		case prevTid == repcmd.TypeIDTech && tid == repcmd.TypeIDCancelTech:
			return repcore.IneffKindFastCancel
		}
	}

	// Too fast repetition of certain commands in a short period of time
	// (regardless of their destinations, if destinations are different/far, then the first one was useless)
	if deltaFrame <= 10 && tid == prevTid {
		switch tid {
		case repcmd.TypeIDStop, repcmd.TypeIDHoldPosition:
			return repcore.IneffKindFastRepetition
		case repcmd.TypeIDTargetedOrder, repcmd.TypeIDTargetedOrder121:
			oid, prevOid := cmd.(*repcmd.TargetedOrderCmd).Order.ID, prevCmd.(*repcmd.TargetedOrderCmd).Order.ID
			if oid == prevOid {
				if repcmd.IsOrderIDKindStop(oid) || repcmd.IsOrderIDKindAttack(oid) || repcmd.IsOrderIDKindHold(oid) {
					return repcore.IneffKindFastRepetition
				}
				switch oid {
				case repcmd.OrderIDMove, repcmd.OrderIDRallyPointUnit, repcmd.OrderIDRallyPointTile:
					return repcore.IneffKindFastRepetition
				}
			}
		}
	}

	// Too fast switch away from or reselecting the same selected unit = no use of selecting it.
	// By too fast I mean it's not even enough to check the units' state.
	if deltaFrame <= 8 && isSelectionChanger(cmd) && isSelectionChanger(prevCmd) {
		// Exclude double tapping the same hotkey: it's only ineffective if tapped more than 3 times
		// (double tapping is used to center the group)
		doubleTap := false
		if he, ok := cmd.(*repcmd.HotkeyCmd); ok {
			if he2, ok2 := prevCmd.(*repcmd.HotkeyCmd); ok2 {
				if he.Group == he2.Group {
					doubleTap = true
					// Is it repeated fast at least 3 times?
					if i >= 2 {
						prevPrevCmd := cmds[i-2]
						if he3, ok3 := prevPrevCmd.(*repcmd.HotkeyCmd); ok3 &&
							he3.HotkeyType.ID == repcmd.HotkeyTypeIDSelect && he3.Group == he.Group &&
							he2.Base.Frame-he3.Base.Frame <= 8 {
							return repcore.IneffKindFastReselection // Same hotkey (select) pressed at least 3 times
						}
					}
				}
			}
		}
		if !doubleTap {
			return repcore.IneffKindFastReselection
		}
	}

	// Repetition of certain commands without time restriction
	if tid == prevTid {
		switch tid {
		case repcmd.TypeIDUnitMorph, repcmd.TypeIDBuildingMorph, repcmd.TypeIDUpgrade, repcmd.TypeIDBuild,
			repcmd.TypeIDMergeArchon, repcmd.TypeIDMergeDarkArchon, repcmd.TypeIDLiftOff,
			repcmd.TypeIDCancelAddon, repcmd.TypeIDCancelBuild, repcmd.TypeIDCancelMorph, repcmd.TypeIDCancelNuke,
			repcmd.TypeIDCancelTech, repcmd.TypeIDCancelUpgrade:
			return repcore.IneffKindRepetition
		}
	}

	// Repetition of the same hotkey assign or add
	if he, ok := cmd.(*repcmd.HotkeyCmd); ok && he.HotkeyType.ID != repcmd.HotkeyTypeIDSelect {
		if he2, ok2 := prevCmd.(*repcmd.HotkeyCmd); ok2 && he2.HotkeyType.ID == he.HotkeyType.ID {
			if he.Group == he2.Group {
				return repcore.IneffKindRepetitionHotkeyAddAssign
			}
		}
	}

	return repcore.IneffKindEffective // If we got this far, classify it as effective
}

// countSameCmds counts how many times the given command is repeated on the same selected units
// without about 1 second.
//
// Counting is capped at 6: even if the command is repeated more times, 6 is returned.
//
// cmd must be cmds[i].
func countSameCmds(cmds []repcmd.Cmd, i int, cmd repcmd.Cmd) (count int) {
	baseCmd := cmd.BaseCmd()
	frameLimit := baseCmd.Frame - 25 // About 1 second

	for ; i >= 0; i-- {
		cmd2 := cmds[i]
		baseCmd2 := cmd2.BaseCmd()
		if baseCmd2.Frame < frameLimit {
			break
		}

		if baseCmd2.Type == baseCmd.Type {
			count++
			if count == 6 {
				break
			}
		} else if isSelectionChanger(cmd2) {
			break
		}
	}

	return
}

// isSelectionChanger tells if the given command (may) change the current selection.
func isSelectionChanger(cmd repcmd.Cmd) bool {
	switch cmd.BaseCmd().Type.ID {
	case repcmd.TypeIDSelect, repcmd.TypeIDSelectAdd, repcmd.TypeIDSelectRemove,
		repcmd.TypeIDSelect121, repcmd.TypeIDSelectAdd121, repcmd.TypeIDSelectRemove121:
		return true
	case repcmd.TypeIDHotkey:
		if cmd.(*repcmd.HotkeyCmd).HotkeyType.ID == repcmd.HotkeyTypeIDSelect {
			return true
		}
	}
	return false
}
