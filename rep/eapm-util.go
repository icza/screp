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

	cmd := cmds[i]

	_ = cmd
	// TODO
	return false
}
