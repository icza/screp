// This file contains unit orders.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Order describes the unit order.
type Order struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// Orders is an enumeration of the possible unit orders.
var Orders = []*Order{
	{e("Die"), 0x00},
	{e("Stop"), 0x01},
	{e("Guard"), 0x02},
	{e("PlayerGuard"), 0x03},
	{e("TurretGuard"), 0x04},
	{e("BunkerGuard"), 0x05},
	{e("Move"), 0x06},
	{e("ReaverStop"), 0x07},
	{e("Attack1"), 0x08},
	{e("Attack2"), 0x09},
	{e("AttackUnit"), 0x0a},
	{e("AttackFixedRange"), 0x0b},
	{e("AttackTile"), 0x0c},
	{e("Hover"), 0x0d},
	{e("AttackMove"), 0x0e},
	{e("InfestedCommandCenter"), 0x0f},
	{e("UnusedNothing"), 0x10},
	{e("UnusedPowerup"), 0x11},
	{e("TowerGuard"), 0x12},
	{e("TowerAttack"), 0x13},
	{e("VultureMine"), 0x14},
	{e("StayInRange"), 0x15},
	{e("TurretAttack"), 0x16},
	{e("Nothing"), 0x17},
	{e("Unused_24"), 0x18},
	{e("DroneStartBuild"), 0x19},
	{e("DroneBuild"), 0x1a},
	{e("CastInfestation"), 0x1b},
	{e("MoveToInfest"), 0x1c},
	{e("InfestingCommandCenter"), 0x1d},
	{e("PlaceBuilding"), 0x1e},
	{e("PlaceProtossBuilding"), 0x1f},
	{e("CreateProtossBuilding"), 0x20},
	{e("ConstructingBuilding"), 0x21},
	{e("Repair"), 0x22},
	{e("MoveToRepair"), 0x23},
	{e("PlaceAddon"), 0x24},
	{e("BuildAddon"), 0x25},
	{e("Train"), 0x26},
	{e("RallyPointUnit"), 0x27},
	{e("RallyPointTile"), 0x28},
	{e("ZergBirth"), 0x29},
	{e("ZergUnitMorph"), 0x2a},
	{e("ZergBuildingMorph"), 0x2b},
	{e("IncompleteBuilding"), 0x2c},
	{e("IncompleteMorphing"), 0x2d},
	{e("BuildNydusExit"), 0x2e},
	{e("EnterNydusCanal"), 0x2f},
	{e("IncompleteWarping"), 0x30},
	{e("Follow"), 0x31},
	{e("Carrier"), 0x32},
	{e("ReaverCarrierMove"), 0x33},
	{e("CarrierStop"), 0x34},
	{e("CarrierAttack"), 0x35},
	{e("CarrierMoveToAttack"), 0x36},
	{e("CarrierIgnore2"), 0x37},
	{e("CarrierFight"), 0x38},
	{e("CarrierHoldPosition"), 0x39},
	{e("Reaver"), 0x3a},
	{e("ReaverAttack"), 0x3b},
	{e("ReaverMoveToAttack"), 0x3c},
	{e("ReaverFight"), 0x3d},
	{e("ReaverHoldPosition"), 0x3e},
	{e("TrainFighter"), 0x3f},
	{e("InterceptorAttack"), 0x40},
	{e("ScarabAttack"), 0x41},
	{e("RechargeShieldsUnit"), 0x42},
	{e("RechargeShieldsBattery"), 0x43},
	{e("ShieldBattery"), 0x44},
	{e("InterceptorReturn"), 0x45},
	{e("DroneLand"), 0x46},
	{e("BuildingLand"), 0x47},
	{e("BuildingLiftOff"), 0x48},
	{e("DroneLiftOff"), 0x49},
	{e("LiftingOff"), 0x4a},
	{e("ResearchTech"), 0x4b},
	{e("Upgrade"), 0x4c},
	{e("Larva"), 0x4d},
	{e("SpawningLarva"), 0x4e},
	{e("Harvest1"), 0x4f},
	{e("Harvest2"), 0x50},
	{e("MoveToGas"), 0x51},
	{e("WaitForGas"), 0x52},
	{e("HarvestGas"), 0x53},
	{e("ReturnGas"), 0x54},
	{e("MoveToMinerals"), 0x55},
	{e("WaitForMinerals"), 0x56},
	{e("MiningMinerals"), 0x57},
	{e("Harvest3"), 0x58},
	{e("Harvest4"), 0x59},
	{e("ReturnMinerals"), 0x5a},
	{e("Interrupted"), 0x5b},
	{e("EnterTransport"), 0x5c},
	{e("PickupIdle"), 0x5d},
	{e("PickupTransport"), 0x5e},
	{e("PickupBunker"), 0x5f},
	{e("Pickup4"), 0x60},
	{e("PowerupIdle"), 0x61},
	{e("Sieging"), 0x62},
	{e("Unsieging"), 0x63},
	{e("WatchTarget"), 0x64},
	{e("InitCreepGrowth"), 0x65},
	{e("SpreadCreep"), 0x66},
	{e("StoppingCreepGrowth"), 0x67},
	{e("GuardianAspect"), 0x68},
	{e("ArchonWarp"), 0x69},
	{e("CompletingArchonSummon"), 0x6a},
	{e("HoldPosition"), 0x6b},
	{e("QueenHoldPosition"), 0x6c},
	{e("Cloak"), 0x6d},
	{e("Decloak"), 0x6e},
	{e("Unload"), 0x6f},
	{e("MoveUnload"), 0x70},
	{e("FireYamatoGun"), 0x71},
	{e("MoveToFireYamatoGun"), 0x72},
	{e("CastLockdown"), 0x73},
	{e("Burrowing"), 0x74},
	{e("Burrowed"), 0x75},
	{e("Unburrowing"), 0x76},
	{e("CastDarkSwarm"), 0x77},
	{e("CastParasite"), 0x78},
	{e("CastSpawnBroodlings"), 0x79},
	{e("CastEMPShockwave"), 0x7a},
	{e("NukeWait"), 0x7b},
	{e("NukeTrain"), 0x7c},
	{e("NukeLaunch"), 0x7d},
	{e("NukePaint"), 0x7e},
	{e("NukeUnit"), 0x7f},
	{e("CastNuclearStrike"), 0x80},
	{e("NukeTrack"), 0x81},
	{e("InitializeArbiter"), 0x82},
	{e("CloakNearbyUnits"), 0x83},
	{e("PlaceMine"), 0x84},
	{e("RightClickAction"), 0x85},
	{e("SuicideUnit"), 0x86},
	{e("SuicideLocation"), 0x87},
	{e("SuicideHoldPosition"), 0x88},
	{e("CastRecall"), 0x89},
	{e("Teleport"), 0x8a},
	{e("CastScannerSweep"), 0x8b},
	{e("Scanner"), 0x8c},
	{e("CastDefensiveMatrix"), 0x8d},
	{e("CastPsionicStorm"), 0x8e},
	{e("CastIrradiate"), 0x8f},
	{e("CastPlague"), 0x90},
	{e("CastConsume"), 0x91},
	{e("CastEnsnare"), 0x92},
	{e("CastStasisField"), 0x93},
	{e("CastHallucination"), 0x94},
	{e("Hallucination2"), 0x95},
	{e("ResetCollision"), 0x96},
	{e("ResetHarvestCollision"), 0x97},
	{e("Patrol"), 0x98},
	{e("CTFCOPInit"), 0x99},
	{e("CTFCOPStarted"), 0x9a},
	{e("CTFCOP2"), 0x9b},
	{e("ComputerAI"), 0x9c},
	{e("AtkMoveEP"), 0x9d},
	{e("HarassMove"), 0x9e},
	{e("AIPatrol"), 0x9f},
	{e("GuardPost"), 0xa0},
	{e("RescuePassive"), 0xa1},
	{e("Neutral"), 0xa2},
	{e("ComputerReturn"), 0xa3},
	{e("InitializePsiProvider"), 0xa4},
	{e("SelfDestructing"), 0xa5},
	{e("Critter"), 0xa6},
	{e("HiddenGun"), 0xa7},
	{e("OpenDoor"), 0xa8},
	{e("CloseDoor"), 0xa9},
	{e("HideTrap"), 0xaa},
	{e("RevealTrap"), 0xab},
	{e("EnableDoodad"), 0xac},
	{e("DisableDoodad"), 0xad},
	{e("WarpIn"), 0xae},
	{e("Medic"), 0xaf},
	{e("MedicHeal"), 0xb0},
	{e("HealMove"), 0xb1},
	{e("MedicHoldPosition"), 0xb2},
	{e("MedicHealToIdle"), 0xb3},
	{e("CastRestoration"), 0xb4},
	{e("CastDisruptionWeb"), 0xb5},
	{e("CastMindControl"), 0xb6},
	{e("DarkArchonMeld"), 0xb7},
	{e("CastFeedback"), 0xb8},
	{e("CastOpticalFlare"), 0xb9},
	{e("CastMaelstrom"), 0xba},
	{e("JunkYardDog"), 0xbb},
	{e("Fatal"), 0xbc},
	{e("None"), 0xbd},
}

// OrderByID returns the Order for a given ID.
// A new Order with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func OrderByID(ID byte) *Order {
	if int(ID) < len(Orders) {
		return Orders[ID]
	}
	return &Order{repcore.UnknownEnum(ID), ID}
}

// Order IDs
const (
	OrderIDStop                 = 0x01
	OrderIDMove                 = 0x06
	OrderIDReaverStop           = 0x07
	OrderIDAttack1              = 0x08
	OrderIDAttack2              = 0x09
	OrderIDAttackUnit           = 0x0a
	OrderIDAttackFixedRange     = 0x0b
	OrderIDAttackTile           = 0x0c
	OrderIDAttackMove           = 0x0e
	OrderIDPlaceProtossBuilding = 0x1f
	OrderIDRallyPointUnit       = 0x27
	OrderIDRallyPointTile       = 0x28
	OrderIDCarrierStop          = 0x34
	OrderIDCarrierAttack        = 0x35
	OrderIDCarrierHoldPosition  = 0x39
	OrderIDReaverHoldPosition   = 0x3e
	OrderIDReaverAttack         = 0x3b
	OrderIDHoldPosition         = 0x6b
	OrderIDQueenHoldPosition    = 0x6c
	OrderIDUnload               = 0x6f
	OrderIDMoveUnload           = 0x70
	OrderIDNukeLaunch           = 0x7d
	OrderIDCastRecall           = 0x89
	OrderIDCastScannerSweep     = 0x8b
	OrderIDMedicHoldPosition    = 0xb2
)

// IsOrderIDKindStop tells if the given order ID is one of the stop orders.
func IsOrderIDKindStop(orderID byte) bool {
	switch orderID {
	case OrderIDStop, OrderIDReaverStop, OrderIDCarrierStop:
		return true
	}
	return false
}

// IsOrderIDKindHold tells if the given order ID is one of the hold orders.
func IsOrderIDKindHold(orderID byte) bool {
	switch orderID {
	case OrderIDHoldPosition, OrderIDCarrierHoldPosition, OrderIDReaverHoldPosition,
		OrderIDQueenHoldPosition, OrderIDMedicHoldPosition:
		return true
	}
	return false
}

// IsOrderIDKindAttack tells if the given order ID is one of the attack orders.
func IsOrderIDKindAttack(orderID byte) bool {
	switch orderID {
	case OrderIDAttack1, OrderIDAttack2, OrderIDAttackUnit, OrderIDAttackFixedRange,
		OrderIDAttackMove, OrderIDCarrierAttack, OrderIDReaverAttack:
		return true
	}
	return false
}
