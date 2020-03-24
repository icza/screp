// This file contains types that model the different commands.

package repcmd

import (
	"fmt"

	"github.com/icza/screp/rep/repcore"
)

// e creates a new Enum value.
func e(name string) repcore.Enum {
	return repcore.Enum{Name: name}
}

// Cmd is the command interface.
type Cmd interface {
	// Base returns the base command.
	BaseCmd() *Base

	// Params returns human-readable concrete command-specific parameters.
	Params() string
}

// Base is the base of all player commands.
type Base struct {
	// Frame at which the command was issued
	Frame repcore.Frame

	// PlayerID this command was issued by
	PlayerID byte

	// Type of the command
	Type *Type
}

// BaseCmd implements Cmd.BaseCmd().
func (b *Base) BaseCmd() *Base {
	return b
}

// Params implements Cmd.Params().
func (b *Base) Params() string {
	return ""
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
func (pec *ParseErrCmd) Params() string {
	prevBase := pec.PrevCmd.BaseCmd()
	return fmt.Sprintf("PrevCmd: [Frame: %d, PlayerID: %d, Type: %s, Params: [%s]",
		prevBase.Frame, prevBase.PlayerID, prevBase.Type, pec.PrevCmd.Params())
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
func (gc *GeneralCmd) Params() string {
	return fmt.Sprintf("Data: [% x]", gc.Data)
}

// SelectCmd describes commands of types: TypeSelect, TypeSelectAdd, TypeSelectRemove
type SelectCmd struct {
	*Base

	// UnitTags contains the unit tags involved in the select command.
	UnitTags []UnitTag
}

// Params implements Cmd.Params().
func (sc *SelectCmd) Params() string {
	return fmt.Sprintf("UnitTags: %x", sc.UnitTags)
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
func (bc *BuildCmd) Params() string {
	return fmt.Sprintf("Order: %v, Pos: (%v), Unit: %v", bc.Order, bc.Pos, bc.Unit)
}

// GameSpeedCmd describes a set game speed command. Type: TypeGameSpeed
type GameSpeedCmd struct {
	*Base

	// Speed is the new game speed.
	Speed *repcore.Speed
}

// Params implements Cmd.Params().
func (gc *GameSpeedCmd) Params() string {
	return fmt.Sprintf("Speed: %v", gc.Speed)
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
func (hc *HotkeyCmd) Params() string {
	return fmt.Sprintf("HotkeyType: %v, Group: %d", hc.HotkeyType, hc.Group)
}

// LeaveGameCmd describes a leave game command. Type: TypeLeaveGame
type LeaveGameCmd struct {
	*Base

	// Speed is the new game speed.
	Reason *LeaveReason
}

// Params implements Cmd.Params().
func (lgc *LeaveGameCmd) Params() string {
	return fmt.Sprintf("Reason: %v", lgc.Reason)
}

// TrainCmd describes a train command. Type: TypeTrain, TypeUnitMorph
type TrainCmd struct {
	*Base

	// Unit is the trained unit.
	Unit *Unit
}

// Params implements Cmd.Params().
func (tc *TrainCmd) Params() string {
	return fmt.Sprintf("Unit: %v", tc.Unit)
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
func (qc *QueueableCmd) Params() string {
	return fmt.Sprintf("Queued: %t", qc.Queued)
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
func (rcc *RightClickCmd) Params() string {
	return fmt.Sprintf("Pos: (%v), UnitTag: %x, Unit: %v, Queued: %t", rcc.Pos, rcc.UnitTag, rcc.Unit, rcc.Queued)
}

// UnloadCmd describes an unload command.
type UnloadCmd struct {
	*Base

	// UnitTag is the unloaded unit's tag if it's valid.
	UnitTag UnitTag
}

// Params implements Cmd.Params().
func (uc *UnloadCmd) Params() string {
	return fmt.Sprintf(" UnitTag: %x", uc.UnitTag)
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
func (toc *TargetedOrderCmd) Params() string {
	return fmt.Sprintf("Pos: (%v), UnitTag: %x, Unit: %v, Order: %v, Queued: %t", toc.Pos, toc.UnitTag, toc.Unit, toc.Order, toc.Queued)
}

// MinimapPingCmd describes a minimap ping command. Type: TypeMinimapPing
type MinimapPingCmd struct {
	*Base

	// Pos tells the pinged location.
	Pos repcore.Point
}

// Params implements Cmd.Params().
func (mpc *MinimapPingCmd) Params() string {
	return fmt.Sprintf("Pos: (%v)", mpc.Pos)
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
func (cc *ChatCmd) Params() string {
	return fmt.Sprintf("SenderSlotID: %d, Message: %q", cc.SenderSlotID, cc.Message)
}

// CancelTrainCmd describes a cancel train command. Type: TypeCancelTrain
type CancelTrainCmd struct {
	*Base

	// UnitTag is the cancelled unit tag.
	UnitTag UnitTag
}

// Params implements Cmd.Params().
func (ctc *CancelTrainCmd) Params() string {
	return fmt.Sprintf("UnitTag: %x", ctc.UnitTag)
}

// BuildingMorphCmd describes a building morph command. Type: TypeBuildingMorph
type BuildingMorphCmd struct {
	*Base

	// Unit is the unit to morph into (e.g. Lair from Hatchery).
	Unit *Unit
}

// Params implements Cmd.Params().
func (bmc *BuildingMorphCmd) Params() string {
	return fmt.Sprintf("Unit: %v", bmc.Unit)
}

// LiftOffCmd describes a lift off command. Type: TypeLiftOff
type LiftOffCmd struct {
	*Base

	// Pos tells the location of the lift off.
	Pos repcore.Point
}

// Params implements Cmd.Params().
func (loc *LiftOffCmd) Params() string {
	return fmt.Sprintf("Pos: (%v)", loc.Pos)
}

// TechCmd describes a tech (research) command. Type: TypeTech
type TechCmd struct {
	*Base

	// Tech that was started.
	Tech *Tech
}

// Params implements Cmd.Params().
func (tc *TechCmd) Params() string {
	return fmt.Sprintf("Tech: %v", tc.Tech)
}

// UpgradeCmd describes an upgrade command. Type: TypeUpgrade
type UpgradeCmd struct {
	*Base

	// Upgrade that was started.
	Upgrade *Upgrade
}

// Params implements Cmd.Params().
func (uc *UpgradeCmd) Params() string {
	return fmt.Sprintf("Upgrade: %v", uc.Upgrade)
}

// LatencyCmd describes a latency change command. Type: TypeLatency
type LatencyCmd struct {
	*Base

	// Latency is the new latency.
	Latency *Latency
}

// Params implements Cmd.Params().
func (lc *LatencyCmd) Params() string {
	return fmt.Sprintf("Latency: %v", lc.Latency)
}
