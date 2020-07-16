package repcore

// IneffKind classifies commands if and why they are ineffective.
type IneffKind byte

const (
	// IneffKindEffective means the command is considered effective.
	IneffKindEffective IneffKind = iota

	// IneffKindUnitQueueOverflow means the command is ineffective due to unit queue overflow
	IneffKindUnitQueueOverflow

	// IneffKindFastCancel means the command is ineffective due to too fast cancel
	IneffKindFastCancel

	// IneffKindFastRepetition means the command is ineffective due to too fast repetition
	IneffKindFastRepetition

	// IneffKindFastReselection means the command is ineffective due to too fast selection change
	// or reselection
	IneffKindFastReselection

	// IneffKindRepetition means the command is ineffective due to repetition
	IneffKindRepetition

	// IneffKindRepetitionHotkeyAddAssign means the command is ineffective due to
	// repeating the same hotkey add or assign
	IneffKindRepetitionHotkeyAddAssign
)

var effectiveKindStrings = []string{
	IneffKindEffective:                 "effective",
	IneffKindUnitQueueOverflow:         "unit queue overflow",
	IneffKindFastCancel:                "too fast cancel",
	IneffKindFastRepetition:            "too fast repetition",
	IneffKindFastReselection:           "too fast selection change or reselection",
	IneffKindRepetition:                "repetition",
	IneffKindRepetitionHotkeyAddAssign: "repeptition of the same hotkey add or assign",
}

// String returns a short string description.
func (k IneffKind) String() string {
	return effectiveKindStrings[k]
}
