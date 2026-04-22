package rep

const (
	LimitsUnitsOriginal = 1700
	LimitsUnitsExtended = 3400
)

// Limits models the limits section added in Remastered.
type Limits struct {
	// Units limit
	Units uint32 `json:"units"`
}
