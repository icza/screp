// This file contains the algorithm implementation for EAPM classification and calculation.

package rep

import "github.com/icza/screp/rep/repcmd"

const (
	// EAPMVersion is a Semver2 compatible version of the EAPM algorithm.
	EAPMVersion = "v1.0.1"
)

// IsCmdEffective tells if a command is considered effective so it can be included in EAPM calculation.
//
// cmds must contain commands of the cmd's player only. It may be a partially filled slice, but must contain
// the player's all commands up to the command in question: len(cmds) > i must hold.
func IsCmdEffective(cmds []repcmd.Cmd, i int) bool {
	if i == 0 {
		return true // First command is effective whatever it is
	}

	// Try to "prove" command is ineffective. If we can't, it's effective.

	cmd := cmds[i]
	tid := cmd.BaseCmd().Type.ID

	// Unit queue overflow
	if tid == repcmd.TypeIDTrain || tid == repcmd.TypeIDTrainFighter {
		if countSameCmds(cmds, i, cmd) >= 6 {
			return false
		}
	}

	frame := cmd.BaseCmd().Frame

	prevCmd := cmds[i-1] // i > 0
	prevTid := prevCmd.BaseCmd().Type.ID
	prevFrame := prevCmd.BaseCmd().Frame

	// Too fast cancel
	if frame-prevFrame <= 20 {
		switch {
		case tid == repcmd.TypeIDTrain && prevTid == repcmd.TypeIDCancelTrain:
			return false
		case (tid == repcmd.TypeIDUnitMorph || tid == repcmd.TypeIDBuildingMorph) && prevTid == repcmd.TypeIDCancelMorph:
			return false
		case tid == repcmd.TypeIDUpgrade && prevTid == repcmd.TypeIDCancelUpgrade:
			return false
		}
	}

	// TODO

	return true // If we got this far: it's effective
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
