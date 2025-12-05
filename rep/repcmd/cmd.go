// This file contains types that model the different commands.

package repcmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/icza/screp/rep/repcore"
)

// Bytes is a []byte that JSON-marshals itself as a number array.
type Bytes []byte

// MarshalJSON marshals the byte slice as a number array.
func (bs Bytes) MarshalJSON() ([]byte, error) {
	if bs == nil {
		return []byte("null"), nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(bs)*3))
	buf.WriteByte('[')
	for i, v := range bs {
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprint(buf, v)
	}
	buf.WriteByte(']')

	return buf.Bytes(), nil
}

// e creates a new Enum value.
func e(name string) repcore.Enum {
	return repcore.Enum{Name: name}
}

// Cmd is the command interface.
type Cmd interface {
	// Base returns the base command.
	BaseCmd() *Base

	// Params returns human-readable concrete command-specific parameters.
	Params(verbose bool) string
}

// Base is the base of all player commands.
type Base struct {
	// Frame at which the command was issued
	Frame repcore.Frame

	// PlayerID this command was issued by
	PlayerID byte

	// Type of the command
	Type *Type

	// IneffKind classification of the command
	IneffKind repcore.IneffKind `json:",omitempty"`
}

// BaseCmd implements Cmd.BaseCmd().
func (b *Base) BaseCmd() *Base {
	return b
}

// Params implements Cmd.Params().
func (b *Base) Params(verbose bool) string {
	return ""
}

// c is a helper function to choose between 2 formats based on verbosity.
func c(verbose bool, verboseFmt, nonVerboseFmt string) string {
	if verbose {
		return verboseFmt
	}
	return nonVerboseFmt
}

// ParseErrCmd represents a command where parsing error encountered.
// It stores a reference to the preceding command for debugging purposes
// (often a parse error is the result of improperly parsing the preceding command).
type ParseErrCmd struct {
	*Base

	// PrevCmd is the command preceding the parse error command.
	PrevCmd Cmd
}

// Params implements Cmd.Params().
func (pec *ParseErrCmd) Params(verbose bool) string {
	prevBase := pec.PrevCmd.BaseCmd()
	return fmt.Sprintf(
		c(verbose,
			"PrevCmd: [Frame: %d, PlayerID: %d, Type: %s, Params: [%s]",
			"[%d, %d, %s, [%s]",
		),
		prevBase.Frame, prevBase.PlayerID, prevBase.Type, pec.PrevCmd.Params(verbose),
	)
}

// UnitTag itentifies a unit in the game (engine). Contains its in-game ID and
// a recycle counter.
type UnitTag uint16

// Index returns the unit's tag index (in-game ID).
func (ut UnitTag) Index() uint16 {
	return uint16(ut) & 0x7ff
}

// Recycle returns the tag resycle.
func (ut UnitTag) Recycle() byte {
	return byte(uint16(ut) >> 12)
}

// Valid tells if this is a valid unit tag.
func (ut UnitTag) Valid() bool {
	return ut != 0xffff
}

// GeneralCmd represents a general command whose parameters
// are not handled / cared for.
type GeneralCmd struct {
	*Base

	// Data is the "raw" parameters of the command.
	Data []byte
}

// Params implements Cmd.Params().
func (gc *GeneralCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Data: [% x]",
			"[% x]",
		),
		gc.Data,
	)
}

// SelectCmd describes commands of types: TypeSelect, TypeSelectAdd, TypeSelectRemove
type SelectCmd struct {
	*Base

	// UnitTags contains the unit tags involved in the select command.
	UnitTags []UnitTag
}

// Params implements Cmd.Params().
func (sc *SelectCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"UnitTags: %x",
			"%x",
		),
		sc.UnitTags,
	)
}

// BuildCmd describes a build command. Type: TypeBuild
type BuildCmd struct {
	*Base

	// Order type
	Order *Order

	// Pos tells the point where the building is placed.
	Pos repcore.Point

	// Unit is the building issued to be built.
	Unit *Unit
}

// Params implements Cmd.Params().
func (bc *BuildCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("Order: %v, Pos: (%v), Unit: %v", bc.Order, bc.Pos, bc.Unit)
	}

	// Order is "redundant" (e.g. PlaceProtossBuilding, DroneStartBuild)
	return fmt.Sprintf("(%v), %v", bc.Pos, bc.Unit)
}

// GameSpeedCmd describes a set game speed command. Type: TypeGameSpeed
type GameSpeedCmd struct {
	*Base

	// Speed is the new game speed.
	Speed *repcore.Speed
}

// Params implements Cmd.Params().
func (gc *GameSpeedCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Speed: %v",
			"%v",
		),
		gc.Speed,
	)
}

// HotkeyCmd describes a hotkey command. Type: TypeHotkey
type HotkeyCmd struct {
	*Base

	// HotkeyType is the type of the hotkey command
	// (named like this to avoid same name from Base.Type).
	HotkeyType *HotkeyType

	// Group (the "number"): 0..9.
	Group byte
}

// Params implements Cmd.Params().
func (hc *HotkeyCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"HotkeyType: %v, Group: %d",
			"%v, %d",
		),
		hc.HotkeyType, hc.Group,
	)
}

// LeaveGameCmd describes a leave game command. Type: TypeLeaveGame
type LeaveGameCmd struct {
	*Base

	// Reasom why the player left.
	Reason *LeaveReason
}

// Params implements Cmd.Params().
func (lgc *LeaveGameCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Reason: %v",
			"%v",
		), lgc.Reason,
	)
}

// TrainCmd describes a train command. Type: TypeTrain, TypeUnitMorph
type TrainCmd struct {
	*Base

	// Unit is the trained unit.
	Unit *Unit
}

// Params implements Cmd.Params().
func (tc *TrainCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Unit: %v",
			"%v",
		),
		tc.Unit,
	)
}

// QueueableCmd describes a generic command that holds whether it is queued.
// Types: TypeStop, TypeReturnCargo, TypeUnloadAll, TypeHoldPosition,
// TypeBurrow, TypeUnburrow, TypeSiege, TypeUnsiege, TypeCloack, TypeDecloack
type QueueableCmd struct {
	*Base

	// Queued tells if the command is queued. If not, it's instant.
	Queued bool
}

// Params implements Cmd.Params().
func (qc *QueueableCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("Queued: %t", qc.Queued)
	}
	if qc.Queued {
		return "Queued"
	}
	return ""
}

// RightClickCmd represents a right click command. Type: TypeRightClick
type RightClickCmd struct {
	*Base

	// Pos tells the right-clicked target point.
	Pos repcore.Point

	// UnitTag is the right-clicked unit's unit tag if it's valid.
	UnitTag UnitTag

	// Unit is the right-clicked unit (if UnitTag is valid).
	Unit *Unit

	// Queued tells if the command is queued. If not, it's instant.
	Queued bool
}

// Params implements Cmd.Params().
func (rcc *RightClickCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("Pos: (%v), UnitTag: %x, Unit: %v, Queued: %t", rcc.Pos, rcc.UnitTag, rcc.Unit, rcc.Queued)
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "(%v)", rcc.Pos)
	if rcc.UnitTag != 0 {
		fmt.Fprintf(b, ", %x", rcc.UnitTag)
	}
	if rcc.Unit.ID != UnitIDNone {
		fmt.Fprintf(b, ", %v", rcc.Unit)
	}
	if rcc.Queued {
		b.WriteString(", Queued")
	}
	return b.String()
}

// UnloadCmd describes an unload command.
type UnloadCmd struct {
	*Base

	// UnitTag is the unloaded unit's tag if it's valid.
	UnitTag UnitTag
}

// Params implements Cmd.Params().
func (uc *UnloadCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			" UnitTag: %x",
			"%x",
		),
		uc.UnitTag,
	)
}

// TargetedOrderCmd describes a targeted order command. Type: TypeTargetedOrder
type TargetedOrderCmd struct {
	*Base

	// Pos tells the targeted order's target point.
	Pos repcore.Point

	// UnitTag is the targeted order's unit tag if it's valid.
	UnitTag UnitTag

	// Unit is the targeted order's unit (if UnitTag is valid).
	Unit *Unit

	// Order type
	Order *Order

	// Queued tells if the command is queued. If not, it's instant.
	Queued bool
}

// Params implements Cmd.Params().
func (toc *TargetedOrderCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("Pos: (%v), UnitTag: %x, Unit: %v, Order: %v, Queued: %t", toc.Pos, toc.UnitTag, toc.Unit, toc.Order, toc.Queued)
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "(%v)", toc.Pos)
	if toc.UnitTag != 0 {
		fmt.Fprintf(b, ", %x", toc.UnitTag)
	}
	if toc.Unit.ID != UnitIDNone {
		fmt.Fprintf(b, ", %v", toc.Unit)
	}
	fmt.Fprintf(b, ", %v", toc.Order)
	if toc.Queued {
		b.WriteString(", Queued")
	}
	return b.String()
}

// MinimapPingCmd describes a minimap ping command. Type: TypeMinimapPing
type MinimapPingCmd struct {
	*Base

	// Pos tells the pinged location.
	Pos repcore.Point
}

// Params implements Cmd.Params().
func (mpc *MinimapPingCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Pos: (%v)",
			"(%v)",
		),
		mpc.Pos,
	)
}

// ChatCmd describes an in-game receive chat command. Type: TypeChat
// Owner of the command receives the message sent by the user identified by SenderSlotID.
type ChatCmd struct {
	*Base

	// SenderSlotID tells the slot ID of the message sender.
	SenderSlotID byte

	// Message sent.
	Message string
}

// Params implements Cmd.Params().
func (cc *ChatCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"SenderSlotID: %d, Message: %q",
			"%d, %q",
		),
		cc.SenderSlotID, cc.Message,
	)
}

// VisionCmd describes the share vision command. Type: TypeIDVision
type VisionCmd struct {
	*Base

	// SlotIDs lists slot IDs the owner shared shared vision with
	SlotIDs Bytes
}

// Params implements Cmd.Params().
func (vc *VisionCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"SlotIDs: %v",
			"%v",
		),
		vc.SlotIDs,
	)
}

// AllianceCmd describes the set alliance command. Type: TypeIDAlliance
type AllianceCmd struct {
	*Base

	// SlotIDs lists slot IDs the owner is allied to.
	// It contains slot IDs in increasing order.
	SlotIDs Bytes

	// AlliedVictory tells if Allied Victory is set.
	AlliedVictory bool
}

// Params implements Cmd.Params().
func (ac *AllianceCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("SlotIDs: %v, AlliedVictory: %t", ac.SlotIDs, ac.AlliedVictory)
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "%v", ac.SlotIDs)
	if ac.AlliedVictory {
		b.WriteString(", AlliedVictory")
	}
	return b.String()
}

// CancelTrainCmd describes a cancel train command. Type: TypeCancelTrain
type CancelTrainCmd struct {
	*Base

	// UnitTag is the cancelled unit tag.
	UnitTag UnitTag
}

// Params implements Cmd.Params().
func (ctc *CancelTrainCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"UnitTag: %x",
			"%x",
		),
		ctc.UnitTag,
	)
}

// BuildingMorphCmd describes a building morph command. Type: TypeBuildingMorph
type BuildingMorphCmd struct {
	*Base

	// Unit is the unit to morph into (e.g. Lair from Hatchery).
	Unit *Unit
}

// Params implements Cmd.Params().
func (bmc *BuildingMorphCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Unit: %v",
			"%v",
		),
		bmc.Unit,
	)
}

// LiftOffCmd describes a lift off command. Type: TypeLiftOff
type LiftOffCmd struct {
	*Base

	// Pos tells the location of the lift off.
	Pos repcore.Point
}

// Params implements Cmd.Params().
func (loc *LiftOffCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Pos: (%v)",
			"(%v)",
		), loc.Pos,
	)
}

// LandCmd describes a land command. Type: TypeBuild
type LandCmd struct {
	*Base

	// Order type
	Order *Order

	// Pos tells the point where the building is landed.
	Pos repcore.Point

	// Unit is the building issued to be landed.
	Unit *Unit
}

// Params implements Cmd.Params().
func (bc *LandCmd) Params(verbose bool) string {
	if verbose {
		return fmt.Sprintf("Order: %v, Pos: (%v), Unit: %v", bc.Order, bc.Pos, bc.Unit)
	}

	// Order is "redundant" (it's always BuildingLand)
	return fmt.Sprintf("(%v), %v", bc.Pos, bc.Unit)
}

// TechCmd describes a tech (research) command. Type: TypeTech
type TechCmd struct {
	*Base

	// Tech that was started.
	Tech *Tech
}

// Params implements Cmd.Params().
func (tc *TechCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Tech: %v",
			"%v",
		),
		tc.Tech,
	)
}

// UpgradeCmd describes an upgrade command. Type: TypeUpgrade
type UpgradeCmd struct {
	*Base

	// Upgrade that was started.
	Upgrade *Upgrade
}

// Params implements Cmd.Params().
func (uc *UpgradeCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Upgrade: %v",
			"%v",
		), uc.Upgrade,
	)
}

// LatencyCmd describes a latency change command. Type: TypeLatency
type LatencyCmd struct {
	*Base

	// Latency is the new latency.
	Latency *Latency
}

// Params implements Cmd.Params().
func (lc *LatencyCmd) Params(verbose bool) string {
	return fmt.Sprintf(
		c(verbose,
			"Latency: %v",
			"%v",
		), lc.Latency,
	)
}

// CheatCmd describes a use cheat command. Type: TypeCheat
type CheatCmd struct {
	*Base

	CheatsBitmap uint32
	CheatCodes   []*CheatCode
}

// Params implements Cmd.Params().
func (lc *CheatCmd) Params(verbose bool) string {
	cheatCodes := ""
	if len(lc.CheatCodes) == 0 {
		cheatCodes = "Cheats Disabled"
	} else {
		sb := &strings.Builder{}
		for i, cc := range lc.CheatCodes {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(cc.Name)
		}
		cheatCodes = sb.String()
	}

	return fmt.Sprintf(
		c(verbose,
			"Cheats: %v",
			"%v",
		), cheatCodes,
	)
}
