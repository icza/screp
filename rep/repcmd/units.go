// This file contains the units.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Unit describes the unit.
type Unit struct {
	repcore.Enum

	// ID as it appears in replays
	ID uint16
}

// Units is an enumeration of the possible units
var Units = []*Unit{
	{e("Marine"), 0x00},
	{e("Ghost"), 0x01},
	{e("Vulture"), 0x02},
	{e("Goliath"), 0x03},
	{e("Goliath Turret"), 0x04},
	{e("Siege Tank (Tank Mode)"), 0x05},
	{e("Siege Tank Turret (Tank Mode)"), 0x06},
	{e("SCV"), 0x07},
	{e("Wraith"), 0x08},
	{e("Science Vessel"), 0x09},
	{e("Gui Motang (Firebat)"), 0x0A},
	{e("Dropship"), 0x0B},
	{e("Battlecruiser"), 0x0C},
	{e("Spider Mine"), 0x0D},
	{e("Nuclear Missile"), 0x0E},
	{e("Terran Civilian"), 0x0F},
	{e("Sarah Kerrigan (Ghost)"), 0x10},
	{e("Alan Schezar (Goliath)"), 0x11},
	{e("Alan Schezar Turret"), 0x12},
	{e("Jim Raynor (Vulture)"), 0x13},
	{e("Jim Raynor (Marine)"), 0x14},
	{e("Tom Kazansky (Wraith)"), 0x15},
	{e("Magellan (Science Vessel)"), 0x16},
	{e("Edmund Duke (Tank Mode)"), 0x17},
	{e("Edmund Duke Turret (Tank Mode)"), 0x18},
	{e("Edmund Duke (Siege Mode)"), 0x19},
	{e("Edmund Duke Turret (Siege Mode)"), 0x1A},
	{e("Arcturus Mengsk (Battlecruiser)"), 0x1B},
	{e("Hyperion (Battlecruiser)"), 0x1C},
	{e("Norad II (Battlecruiser)"), 0x1D},
	{e("Terran Siege Tank (Siege Mode)"), 0x1E},
	{e("Siege Tank Turret (Siege Mode)"), 0x1F},
	{e("Firebat"), 0x20},
	{e("Scanner Sweep"), 0x21},
	{e("Medic"), 0x22},
	{e("Larva"), 0x23},
	{e("Egg"), 0x24},
	{e("Zergling"), 0x25},
	{e("Hydralisk"), 0x26},
	{e("Ultralisk"), 0x27},
	{e("Drone"), 0x29},
	{e("Overlord"), 0x2A},
	{e("Mutalisk"), 0x2B},
	{e("Guardian"), 0x2C},
	{e("Queen"), 0x2D},
	{e("Defiler"), 0x2E},
	{e("Scourge"), 0x2F},
	{e("Torrasque (Ultralisk)"), 0x30},
	{e("Matriarch (Queen)"), 0x31},
	{e("Infested Terran"), 0x32},
	{e("Infested Kerrigan (Infested Terran)"), 0x33},
	{e("Unclean One (Defiler)"), 0x34},
	{e("Hunter Killer (Hydralisk)"), 0x35},
	{e("Devouring One (Zergling)"), 0x36},
	{e("Kukulza (Mutalisk)"), 0x37},
	{e("Kukulza (Guardian)"), 0x38},
	{e("Yggdrasill (Overlord)"), 0x39},
	{e("Valkyrie"), 0x3A},
	{e("Mutalisk Cocoon"), 0x3B},
	{e("Corsair"), 0x3C},
	{e("Dark Templar"), 0x3D},
	{e("Devourer"), 0x3E},
	{e("Dark Archon"), 0x3F},
	{e("Probe"), 0x40},
	{e("Zealot"), 0x41},
	{e("Dragoon"), 0x42},
	{e("High Templar"), 0x43},
	{e("Archon"), 0x44},
	{e("Shuttle"), 0x45},
	{e("Scout"), 0x46},
	{e("Arbiter"), 0x47},
	{e("Carrier"), 0x48},
	{e("Interceptor"), 0x49},
	{e("Protoss Dark Templar (Hero)"), 0x4A},
	{e("Zeratul (Dark Templar)"), 0x4B},
	{e("Tassadar/Zeratul (Archon)"), 0x4C},
	{e("Fenix (Zealot)"), 0x4D},
	{e("Fenix (Dragoon)"), 0x4E},
	{e("Tassadar (Templar)"), 0x4F},
	{e("Mojo (Scout)"), 0x50},
	{e("Warbringer (Reaver)"), 0x51},
	{e("Gantrithor (Carrier)"), 0x52},
	{e("Reaver"), 0x53},
	{e("Observer"), 0x54},
	{e("Scarab"), 0x55},
	{e("Danimoth (Arbiter)"), 0x56},
	{e("Aldaris (Templar)"), 0x57},
	{e("Artanis (Scout)"), 0x58},
	{e("Rhynadon (Badlands Critter)"), 0x59},
	{e("Bengalaas (Jungle Critter)"), 0x5A},
	{e("Cargo Ship (Unused)"), 0x5B},
	{e("Mercenary Gunship (Unused)"), 0x5C},
	{e("Scantid (Desert Critter)"), 0x5D},
	{e("Kakaru (Twilight Critter)"), 0x5E},
	{e("Ragnasaur (Ashworld Critter)"), 0x5F},
	{e("Ursadon (Ice World Critter)"), 0x60},
	{e("Lurker Egg"), 0x61},
	{e("Raszagal (Corsair)"), 0x62},
	{e("Samir Duran (Ghost)"), 0x63},
	{e("Alexei Stukov (Ghost)"), 0x64},
	{e("Map Revealer"), 0x65},
	{e("Gerard DuGalle (BattleCruiser)"), 0x66},
	{e("Lurker"), 0x67},
	{e("Infested Duran (Infested Terran)"), 0x68},
	{e("Disruption Web"), 0x69},
	{e("Command Center"), 0x6A},
	{e("ComSat"), 0x6B},
	{e("Nuclear Silo"), 0x6C},
	{e("Supply Depot"), 0x6D},
	{e("Refinery"), 0x6E},
	{e("Barracks"), 0x6F},
	{e("Academy"), 0x70},
	{e("Factory"), 0x71},
	{e("Starport"), 0x72},
	{e("Control Tower"), 0x73},
	{e("Science Facility"), 0x74},
	{e("Covert Ops"), 0x75},
	{e("Physics Lab"), 0x76},
	{e("Machine Shop"), 0x78},
	{e("Repair Bay (Unused)"), 0x79},
	{e("Engineering Bay"), 0x7A},
	{e("Armory"), 0x7B},
	{e("Missile Turret"), 0x7C},
	{e("Bunker"), 0x7D},
	{e("Norad II (Crashed)"), 0x7E},
	{e("Ion Cannon"), 0x7F},
	{e("Uraj Crystal"), 0x80},
	{e("Khalis Crystal"), 0x81},
	{e("Infested CC"), 0x82},
	{e("Hatchery"), 0x83},
	{e("Lair"), 0x84},
	{e("Hive"), 0x85},
	{e("Nydus Canal"), 0x86},
	{e("Hydralisk Den"), 0x87},
	{e("Defiler Mound"), 0x88},
	{e("Greater Spire"), 0x89},
	{e("Queens Nest"), 0x8A},
	{e("Evolution Chamber"), 0x8B},
	{e("Ultralisk Cavern"), 0x8C},
	{e("Spire"), 0x8D},
	{e("Spawning Pool"), 0x8E},
	{e("Creep Colony"), 0x8F},
	{e("Spore Colony"), 0x90},
	{e("Unused Zerg Building1"), 0x91},
	{e("Sunken Colony"), 0x92},
	{e("Zerg Overmind (With Shell)"), 0x93},
	{e("Overmind"), 0x94},
	{e("Extractor"), 0x95},
	{e("Mature Chrysalis"), 0x96},
	{e("Cerebrate"), 0x97},
	{e("Cerebrate Daggoth"), 0x98},
	{e("Unused Zerg Building2"), 0x99},
	{e("Nexus"), 0x9A},
	{e("Robotics Facility"), 0x9B},
	{e("Pylon"), 0x9C},
	{e("Assimilator"), 0x9D},
	{e("Unused Protoss Building1"), 0x9E},
	{e("Observatory"), 0x9F},
	{e("Gateway"), 0xA0},
	{e("Unused Protoss Building2"), 0xA1},
	{e("Photon Cannon"), 0xA2},
	{e("Citadel of Adun"), 0xA3},
	{e("Cybernetics Core"), 0xA4},
	{e("Templar Archives"), 0xA5},
	{e("Forge"), 0xA6},
	{e("Stargate"), 0xA7},
	{e("Stasis Cell/Prison"), 0xA8},
	{e("Fleet Beacon"), 0xA9},
	{e("Arbiter Tribunal"), 0xAA},
	{e("Robotics Support Bay"), 0xAB},
	{e("Shield Battery"), 0xAC},
	{e("Khaydarin Crystal Formation"), 0xAD},
	{e("Protoss Temple"), 0xAE},
	{e("Xel'Naga Temple"), 0xAF},
	{e("Mineral Field (Type 1)"), 0xB0},
	{e("Mineral Field (Type 2)"), 0xB1},
	{e("Mineral Field (Type 3)"), 0xB2},
	{e("Cave (Unused)"), 0xB3},
	{e("Cave-in (Unused)"), 0xB4},
	{e("Cantina (Unused)"), 0xB5},
	{e("Mining Platform (Unused)"), 0xB6},
	{e("Independent Command Center (Unused)"), 0xB7},
	{e("Independent Starport (Unused)"), 0xB8},
	{e("Independent Jump Gate (Unused)"), 0xB9},
	{e("Ruins (Unused)"), 0xBA},
	{e("Khaydarin Crystal Formation (Unused)"), 0xBB},
	{e("Vespene Geyser"), 0xBC},
	{e("Warp Gate"), 0xBD},
	{e("Psi Disrupter"), 0xBE},
	{e("Zerg Marker"), 0xBF},
	{e("Terran Marker"), 0xC0},
	{e("Protoss Marker"), 0xC1},
	{e("Zerg Beacon"), 0xC2},
	{e("Terran Beacon"), 0xC3},
	{e("Protoss Beacon"), 0xC4},
	{e("Zerg Flag Beacon"), 0xC5},
	{e("Terran Flag Beacon"), 0xC6},
	{e("Protoss Flag Beacon"), 0xC7},
	{e("Power Generator"), 0xC8},
	{e("Overmind Cocoon"), 0xC9},
	{e("Dark Swarm"), 0xCA},
	{e("Floor Missile Trap"), 0xCB},
	{e("Floor Hatch (Unused)"), 0xCC},
	{e("Left Upper Level Door"), 0xCD},
	{e("Right Upper Level Door"), 0xCE},
	{e("Left Pit Door"), 0xCF},
	{e("Right Pit Door"), 0xD0},
	{e("Floor Gun Trap"), 0xD1},
	{e("Left Wall Missile Trap"), 0xD2},
	{e("Left Wall Flame Trap"), 0xD3},
	{e("Right Wall Missile Trap"), 0xD4},
	{e("Right Wall Flame Trap"), 0xD5},
	{e("Start Location"), 0xD6},
	{e("Flag"), 0xD7},
	{e("Young Chrysalis"), 0xD8},
	{e("Psi Emitter"), 0xD9},
	{e("Data Disc"), 0xDA},
	{e("Khaydarin Crystal"), 0xDB},
	{e("Mineral Cluster Type 1"), 0xDC},
	{e("Mineral Cluster Type 2"), 0xDD},
	{e("Protoss Vespene Gas Orb Type 1"), 0xDE},
	{e("Protoss Vespene Gas Orb Type 2"), 0xDF},
	{e("Zerg Vespene Gas Sac Type 1"), 0xE0},
	{e("Zerg Vespene Gas Sac Type 2"), 0xE1},
	{e("Terran Vespene Gas Tank Type 1"), 0xE2},
	{e("Terran Vespene Gas Tank Type 2"), 0xE3},
	{e("None"), 0xE4},
}

// unitIDUnit maps from unit ID to unit.
var unitIDUnit = map[uint16]*Unit{}

func init() {
	for _, u := range Units {
		unitIDUnit[u.ID] = u
	}
}

// Unit IDs
const (
	// Critters
	UnitIDRhynadon  = 0x59
	UnitIDBengalaas = 0x5a
	UnitIDScantid   = 0x5d
	UnitIDKakaru    = 0x5e
	UnitIDRagnasaur = 0x5f
	UnitIDUrsadon   = 0x60

	UnitIDCommandCenter   = 0x6A
	UnitIDComSat          = 0x6B
	UnitIDNuclearSilo     = 0x6C
	UnitIDSupplyDepot     = 0x6D
	UnitIDRefinery        = 0x6E
	UnitIDBarracks        = 0x6F
	UnitIDAcademy         = 0x70
	UnitIDFactory         = 0x71
	UnitIDStarport        = 0x72
	UnitIDControlTower    = 0x73
	UnitIDScienceFacility = 0x74
	UnitIDCovertOps       = 0x75
	UnitIDPhysicsLab      = 0x76
	UnitIDMachineShop     = 0x78
	UnitIDEngineeringBay  = 0x7A
	UnitIDArmory          = 0x7B
	UnitIDMissileTurret   = 0x7C
	UnitIDBunker          = 0x7D

	UnitIDInfestedCC       = 0x82
	UnitIDHatchery         = 0x83
	UnitIDLair             = 0x84
	UnitIDHive             = 0x85
	UnitIDNydusCanal       = 0x86
	UnitIDHydraliskDen     = 0x87
	UnitIDDefilerMound     = 0x88
	UnitIDGreaterSpire     = 0x89
	UnitIDQueensNest       = 0x8A
	UnitIDEvolutionChamber = 0x8B
	UnitIDUltraliskCavern  = 0x8C
	UnitIDSpire            = 0x8D
	UnitIDSpawningPool     = 0x8E
	UnitIDCreepColony      = 0x8F
	UnitIDSporeColony      = 0x90
	UnitIDSunkenColony     = 0x92
	UnitIDExtractor        = 0x95

	UnitIDNexus              = 0x9A
	UnitIDRoboticsFacility   = 0x9B
	UnitIDPylon              = 0x9C
	UnitIDAssimilator        = 0x9D
	UnitIDObservatory        = 0x9F
	UnitIDGateway            = 0xA0
	UnitIDPhotonCannon       = 0xA2
	UnitIDCitadelOfAdun      = 0xA3
	UnitIDCyberneticsCore    = 0xA4
	UnitIDTemplarArchives    = 0xA5
	UnitIDForge              = 0xA6
	UnitIDStargate           = 0xA7
	UnitIDFleetBeacon        = 0xA9
	UnitIDArbiterTribunal    = 0xAA
	UnitIDRoboticsSupportBay = 0xAB
	UnitIDShieldBattery      = 0xAC

	UnitIDMineralField1 = 0xB0
	UnitIDMineralField2 = 0xB1
	UnitIDMineralField3 = 0xB2
	UnitIDVespeneGeyser = 0xBC
	UnitIDStartLocation = 0xD6

	UnitIDNone = 0xE4
)

// UnitByID returns the Unit for a given ID.
// A new Unit with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func UnitByID(ID uint16) *Unit {
	if u := unitIDUnit[ID]; u != nil {
		return u
	}
	return &Unit{repcore.UnknownEnum(ID), ID}
}

// unitIDRace maps from unit ID to owner race.
var unitIDRace = map[uint16]*repcore.Race{
	UnitIDCommandCenter:   repcore.RaceTerran,
	UnitIDComSat:          repcore.RaceTerran,
	UnitIDNuclearSilo:     repcore.RaceTerran,
	UnitIDSupplyDepot:     repcore.RaceTerran,
	UnitIDRefinery:        repcore.RaceTerran,
	UnitIDBarracks:        repcore.RaceTerran,
	UnitIDAcademy:         repcore.RaceTerran,
	UnitIDFactory:         repcore.RaceTerran,
	UnitIDStarport:        repcore.RaceTerran,
	UnitIDControlTower:    repcore.RaceTerran,
	UnitIDScienceFacility: repcore.RaceTerran,
	UnitIDCovertOps:       repcore.RaceTerran,
	UnitIDPhysicsLab:      repcore.RaceTerran,
	UnitIDMachineShop:     repcore.RaceTerran,
	UnitIDEngineeringBay:  repcore.RaceTerran,
	UnitIDArmory:          repcore.RaceTerran,
	UnitIDMissileTurret:   repcore.RaceTerran,
	UnitIDBunker:          repcore.RaceTerran,

	UnitIDInfestedCC:       repcore.RaceZerg,
	UnitIDHatchery:         repcore.RaceZerg,
	UnitIDLair:             repcore.RaceZerg,
	UnitIDHive:             repcore.RaceZerg,
	UnitIDNydusCanal:       repcore.RaceZerg,
	UnitIDHydraliskDen:     repcore.RaceZerg,
	UnitIDDefilerMound:     repcore.RaceZerg,
	UnitIDGreaterSpire:     repcore.RaceZerg,
	UnitIDQueensNest:       repcore.RaceZerg,
	UnitIDEvolutionChamber: repcore.RaceZerg,
	UnitIDUltraliskCavern:  repcore.RaceZerg,
	UnitIDSpire:            repcore.RaceZerg,
	UnitIDSpawningPool:     repcore.RaceZerg,
	UnitIDCreepColony:      repcore.RaceZerg,
	UnitIDSporeColony:      repcore.RaceZerg,
	UnitIDSunkenColony:     repcore.RaceZerg,
	UnitIDExtractor:        repcore.RaceZerg,

	UnitIDNexus:              repcore.RaceProtoss,
	UnitIDRoboticsFacility:   repcore.RaceProtoss,
	UnitIDPylon:              repcore.RaceProtoss,
	UnitIDAssimilator:        repcore.RaceProtoss,
	UnitIDObservatory:        repcore.RaceProtoss,
	UnitIDGateway:            repcore.RaceProtoss,
	UnitIDPhotonCannon:       repcore.RaceProtoss,
	UnitIDCitadelOfAdun:      repcore.RaceProtoss,
	UnitIDCyberneticsCore:    repcore.RaceProtoss,
	UnitIDTemplarArchives:    repcore.RaceProtoss,
	UnitIDForge:              repcore.RaceProtoss,
	UnitIDStargate:           repcore.RaceProtoss,
	UnitIDFleetBeacon:        repcore.RaceProtoss,
	UnitIDArbiterTribunal:    repcore.RaceProtoss,
	UnitIDRoboticsSupportBay: repcore.RaceProtoss,
	UnitIDShieldBattery:      repcore.RaceProtoss,
}

// RaceOfUnitID returns the owner race of the unit given by its ID.
// Returns nil if owner is unknown.
// Currently only building units are recognized.
func RaceOfUnitID(ID uint16) *repcore.Race {
	if r := unitIDRace[ID]; r != nil {
		return r
	}
	return nil
}
