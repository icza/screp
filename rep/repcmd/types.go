// This file contains the command types.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Type IDs of command types
const (
	TypeIDKeepAlive          byte = 0x05
	TypeIDSaveGame           byte = 0x06
	TypeIDLoadGame           byte = 0x07
	TypeIDRestartGame        byte = 0x08
	TypeIDSelect             byte = 0x09
	TypeIDSelectAdd          byte = 0x0a
	TypeIDSelectRemove       byte = 0x0b
	TypeIDBuild              byte = 0x0c
	TypeIDVision             byte = 0x0d
	TypeIDAlliance           byte = 0x0e
	TypeIDGameSpeed          byte = 0x0f
	TypeIDPause              byte = 0x10
	TypeIDResume             byte = 0x11
	TypeIDCheat              byte = 0x12
	TypeIDHotkey             byte = 0x13
	TypeIDRightClick         byte = 0x14
	TypeIDTargetedOrder      byte = 0x15
	TypeIDCancelBuild        byte = 0x18
	TypeIDCancelMorph        byte = 0x19
	TypeIDStop               byte = 0x1a
	TypeIDCarrierStop        byte = 0x1b
	TypeIDReaverStop         byte = 0x1c
	TypeIDOrderNothing       byte = 0x1d
	TypeIDReturnCargo        byte = 0x1e
	TypeIDTrain              byte = 0x1f
	TypeIDCancelTrain        byte = 0x20
	TypeIDCloack             byte = 0x21
	TypeIDDecloack           byte = 0x22
	TypeIDUnitMorph          byte = 0x23
	TypeIDUnsiege            byte = 0x25
	TypeIDSiege              byte = 0x26
	TypeIDTrainFighter       byte = 0x27 // Build interceptor / scarab
	TypeIDUnloadAll          byte = 0x28
	TypeIDUnload             byte = 0x29
	TypeIDMergeArchon        byte = 0x2a
	TypeIDHoldPosition       byte = 0x2b
	TypeIDBurrow             byte = 0x2c
	TypeIDUnburrow           byte = 0x2d
	TypeIDCancelNuke         byte = 0x2e
	TypeIDLiftOff            byte = 0x2f
	TypeIDTech               byte = 0x30
	TypeIDCancelTech         byte = 0x31
	TypeIDUpgrade            byte = 0x32
	TypeIDCancelUpgrade      byte = 0x33
	TypeIDCancelAddon        byte = 0x34
	TypeIDBuildingMorph      byte = 0x35
	TypeIDStim               byte = 0x36
	TypeIDSync               byte = 0x37
	TypeIDVoiceEnable        byte = 0x38
	TypeIDVoiceDisable       byte = 0x39
	TypeIDVoiceSquelch       byte = 0x3a
	TypeIDVoiceUnsquelch     byte = 0x3b
	TypeIDStartGame          byte = 0x3c
	TypeIDDownloadPercentage byte = 0x3d
	TypeIDChangeGameSlot     byte = 0x3e
	TypeIDNewNetPlayer       byte = 0x3f
	TypeIDJoinedGame         byte = 0x40
	TypeIDChangeRace         byte = 0x41
	TypeIDTeamGameTeam       byte = 0x42
	TypeIDUMSTeam            byte = 0x43
	TypeIDMeleeTeam          byte = 0x44
	TypeIDSwapPlayers        byte = 0x45
	TypeIDSavedData          byte = 0x48
	TypeIDBriefingStart      byte = 0x54
	TypeIDLatency            byte = 0x55
	TypeIDReplaySpeed        byte = 0x56
	TypeIDLeaveGame          byte = 0x57
	TypeIDMinimapPing        byte = 0x58
	TypeIDMergeDarkArchon    byte = 0x5a
	TypeIDMakeGamePublic     byte = 0x5b
	TypeIDChat               byte = 0x5c
	TypeIDSelect121          byte = 0x63
)

// Type describes the command type.
type Type struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// Types is an enumeration of the possible command types
var Types = []*Type{
	{e("Keep Alive"), TypeIDKeepAlive},
	{e("Save Game"), TypeIDSaveGame},
	{e("Load Game"), TypeIDLoadGame},
	{e("Restart Game"), TypeIDRestartGame},
	{e("Select"), TypeIDSelect},
	{e("Select Add"), TypeIDSelectAdd},
	{e("Select Remove"), TypeIDSelectRemove},
	{e("Build"), TypeIDBuild},
	{e("Vision"), TypeIDVision},
	{e("Alliance"), TypeIDAlliance},
	{e("Game Speed"), TypeIDGameSpeed},
	{e("Pause"), TypeIDPause},
	{e("Resume"), TypeIDResume},
	{e("Cheat"), TypeIDCheat},
	{e("Hotkey"), TypeIDHotkey},
	{e("Right Click"), TypeIDRightClick},
	{e("Targeted Order"), TypeIDTargetedOrder},
	{e("Cancel Build"), TypeIDCancelBuild},
	{e("Cancel Morph"), TypeIDCancelMorph},
	{e("Stop"), TypeIDStop},
	{e("Carrier Stop"), TypeIDCarrierStop},
	{e("Reaver Stop"), TypeIDReaverStop},
	{e("Order Nothing"), TypeIDOrderNothing},
	{e("Return Cargo"), TypeIDReturnCargo},
	{e("Train"), TypeIDTrain},
	{e("Cancel Train"), TypeIDCancelTrain},
	{e("Cloack"), TypeIDCloack},
	{e("Decloack"), TypeIDDecloack},
	{e("Unit Morph"), TypeIDUnitMorph},
	{e("Unsiege"), TypeIDUnsiege},
	{e("Siege"), TypeIDSiege},
	{e("Train Fighter"), TypeIDTrainFighter}, // Build interceptor / scarab
	{e("Unload All"), TypeIDUnloadAll},
	{e("Unload"), TypeIDUnload},
	{e("Merge Archon"), TypeIDMergeArchon},
	{e("Hold Position"), TypeIDHoldPosition},
	{e("Burrow"), TypeIDBurrow},
	{e("Unburrow"), TypeIDUnburrow},
	{e("Cancel Nuke"), TypeIDCancelNuke},
	{e("Lift Off"), TypeIDLiftOff},
	{e("Tech"), TypeIDTech},
	{e("Cancel Tech"), TypeIDCancelTech},
	{e("Upgrade"), TypeIDUpgrade},
	{e("Cancel Upgrade"), TypeIDCancelUpgrade},
	{e("Cancel Addon"), TypeIDCancelAddon},
	{e("Building Morph"), TypeIDBuildingMorph},
	{e("Stim"), TypeIDStim},
	{e("Sync"), TypeIDSync},
	{e("Voice Enable"), TypeIDVoiceEnable},
	{e("Voice Disable"), TypeIDVoiceDisable},
	{e("Voice Squelch"), TypeIDVoiceSquelch},
	{e("Voice Unsquelch"), TypeIDVoiceUnsquelch},
	{e("[Lobby] Start Game"), TypeIDStartGame},
	{e("[Lobby] Download Percentage"), TypeIDDownloadPercentage},
	{e("[Lobby] Change Game Slot"), TypeIDChangeGameSlot},
	{e("[Lobby] New Net Player"), TypeIDNewNetPlayer},
	{e("[Lobby] Joined Game"), TypeIDJoinedGame},
	{e("[Lobby] Change Race"), TypeIDChangeRace},
	{e("[Lobby] Team Game Team"), TypeIDTeamGameTeam},
	{e("[Lobby] UMS Team"), TypeIDUMSTeam},
	{e("[Lobby] Melee Team"), TypeIDMeleeTeam},
	{e("[Lobby] Swap Players"), TypeIDSwapPlayers},
	{e("[Lobby] Saved Data"), TypeIDSavedData},
	{e("Briefing Start"), TypeIDBriefingStart},
	{e("Latency"), TypeIDLatency},
	{e("Replay Speed"), TypeIDReplaySpeed},
	{e("Leave Game"), TypeIDLeaveGame},
	{e("Minimap Ping"), TypeIDMinimapPing},
	{e("Merge Dark Archon"), TypeIDMergeDarkArchon},
	{e("Make Game Public"), TypeIDMakeGamePublic},
	{e("Chat"), TypeIDChat},
	{e("Select121"), TypeIDSelect121},
}

// Named command types
var (
	TypeKeepAlive          = Types[0]
	TypeSaveGame           = Types[1]
	TypeLoadGame           = Types[2]
	TypeRestartGame        = Types[3]
	TypeSelect             = Types[4]
	TypeSelectAdd          = Types[5]
	TypeSelectRemove       = Types[6]
	TypeBuild              = Types[7]
	TypeVision             = Types[8]
	TypeAlliance           = Types[9]
	TypeGameSpeed          = Types[10]
	TypePause              = Types[11]
	TypeResume             = Types[12]
	TypeCheat              = Types[13]
	TypeHotkey             = Types[14]
	TypeRightClick         = Types[15]
	TypeTargetedOrder      = Types[16]
	TypeCancelBuild        = Types[17]
	TypeCancelMorph        = Types[18]
	TypeStop               = Types[19]
	TypeCarrierStop        = Types[20]
	TypeReaverStop         = Types[21]
	TypeOrderNothing       = Types[22]
	TypeReturnCargo        = Types[23]
	TypeTrain              = Types[24]
	TypeCancelTrain        = Types[25]
	TypeCloack             = Types[26]
	TypeDecloack           = Types[27]
	TypeUnitMorph          = Types[28]
	TypeUnsiege            = Types[29]
	TypeSiege              = Types[30]
	TypeTrainFighter       = Types[31] // Build interceptor / scarab
	TypeUnloadAll          = Types[32]
	TypeUnload             = Types[33]
	TypeMergeArchon        = Types[34]
	TypeHoldPosition       = Types[35]
	TypeBurrow             = Types[36]
	TypeUnburrow           = Types[37]
	TypeCancelNuke         = Types[38]
	TypeLiftOff            = Types[39]
	TypeTech               = Types[40]
	TypeCancelTech         = Types[41]
	TypeUpgrade            = Types[42]
	TypeCancelUpgrade      = Types[43]
	TypeCancelAddon        = Types[44]
	TypeBuildingMorph      = Types[45]
	TypeStim               = Types[46]
	TypeSync               = Types[47]
	TypeVoiceEnable        = Types[48]
	TypeVoiceDisable       = Types[49]
	TypeVoiceSquelch       = Types[50]
	TypeVoiceUnsquelch     = Types[51]
	TypeStartGame          = Types[52]
	TypeDownloadPercentage = Types[53]
	TypeChangeGameSlot     = Types[54]
	TypeNewNetPlayer       = Types[55]
	TypeJoinedGame         = Types[56]
	TypeChangeRace         = Types[57]
	TypeTeamGameTeam       = Types[58]
	TypeUMSTeam            = Types[59]
	TypeMeleeTeam          = Types[60]
	TypeSwapPlayers        = Types[61]
	TypeSavedData          = Types[62]
	TypeBriefingStart      = Types[63]
	TypeLatency            = Types[64]
	TypeReplaySpeed        = Types[65]
	TypeLeaveGame          = Types[66]
	TypeMinimapPing        = Types[67]
	TypeMergeDarkArchon    = Types[68]
	TypeMakeGamePublic     = Types[69]
	TypeChat               = Types[70]
	TypeSelect121          = Types[71]
)

// typeIDType maps from type ID to type.
var typeIDType = map[byte]*Type{}

func init() {
	for _, t := range Types {
		typeIDType[t.ID] = t
	}
}

// TypeByID returns the Type for a given ID.
// A new Type with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func TypeByID(ID byte) *Type {
	if t := typeIDType[ID]; t != nil {
		return t
	}
	return &Type{repcore.UnknownEnum(ID), ID}
}
