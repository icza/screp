/*
Package repparser implements StarCraft: Brood War replay parsing.

The package is safe for concurrent use.

Information sources:

BWHF replay parser:

https://github.com/icza/bwhf/tree/master/src/hu/belicza/andras/bwhf/control

BWAPI replay parser:

https://github.com/bwapi/bwapi/tree/master/bwapi/libReplayTool

https://github.com/bwapi/bwapi/tree/master/bwapi/include/BWAPI

https://github.com/bwapi/bwapi/tree/master/bwapi/PKLib

Command models:

https://github.com/icza/bwhf/blob/master/src/hu/belicza/andras/bwhf/model/Action.java

https://github.com/bwapi/bwapi/tree/master/bwapi/libReplayTool

jssuh replay parser:

https://github.com/neivv/jssuh

Map Data format:

https://www.starcraftai.com/wiki/CHK_Format

http://www.staredit.net/wiki/index.php/Scenario.chk

http://blog.naver.com/PostView.nhn?blogId=wisdomswrap&logNo=60119755717&parentCategoryNo=&categoryNo=19&viewDate=&isShowPopularPosts=false&from=postView

https://github.com/ShieldBattery/bw-chk/blob/master/index.js
*/
package repparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/rep/repcmd"
	"github.com/icza/screp/rep/repcore"
	"github.com/icza/screp/repparser/repdecoder"
	"golang.org/x/text/encoding/korean"
)

const (
	// Version is a Semver2 compatible version of the parser.
	Version = "v1.12.13"
)

var (
	// ErrNotReplayFile indicates the given file (or reader) is not a valid
	// replay file
	ErrNotReplayFile = errors.New("not a replay file")

	// ErrParsing indicates that an unexpected error occurred, which may be
	// due to corrupt / invalid replay file, or some implementation error.
	ErrParsing = errors.New("parsing")
)

// Config holds parser configuration.
type Config struct {
	// Commands tells if the commands section is to be parsed
	Commands bool

	// MapData tells if the map data section is to be parsed
	MapData bool

	// Debug tells if debug and replay internal binaries is to be retained in the returned Replay.
	Debug bool

	// MapGraphics tells if map data usually required for map image rendering is to be parsed.
	// MapData must be parsed too.
	MapGraphics bool

	// Custom logger to use to report parsing errors.
	// If nil, the default logger of the log package will be used.
	// To suppress logs, use a new logger directed to io.Discard, e.g.:
	// discardLogger := log.New(io.Discard, "", 0)
	Logger *log.Logger

	_ struct{} // To prevent unkeyed literals
}

// ParseFile parses all sections from an SC:BW replay file.
func ParseFile(name string) (r *rep.Replay, err error) {
	return ParseFileConfig(name, Config{Commands: true, MapData: true})
}

// ParseFileSections parses an SC:BW replay file.
// Parsing commands and map data sections depends on the given parameters.
// Replay ID and header sections are always parsed.
//
// Deprecated: Use ParseFileConfig() instead.
func ParseFileSections(name string, commands, mapData bool) (r *rep.Replay, err error) {
	return ParseFileConfig(name, Config{Commands: commands, MapData: mapData})
}

// ParseFileConfig parses an SC:BW replay file based on the given parser configuration.
// Replay ID and header sections are always parsed.
func ParseFileConfig(name string, cfg Config) (r *rep.Replay, err error) {
	dec, err := repdecoder.NewFromFile(name)
	if err != nil {
		return nil, err
	}
	defer dec.Close()

	return parseProtected(dec, cfg)
}

// Parse parses all sections of an SC:BW replay from the given byte slice.
// Map graphics related info is not parsed (see Config.MapGraphics).
func Parse(repData []byte) (*rep.Replay, error) {
	return ParseConfig(repData, Config{Commands: true, MapData: true})
}

// ParseSections parses an SC:BW replay from the given byte slice.
// Parsing commands and map data sections depends on the given parameters.
// Replay ID and header sections are always parsed.
//
// Deprecated: Use ParseConfig() instead.
func ParseSections(repData []byte, commands, mapData bool) (*rep.Replay, error) {
	return ParseConfig(repData, Config{Commands: commands, MapData: mapData})
}

// ParseConfig parses an SC:BW replay from the given byte sice based on the given parser configuration.
// Replay ID and header sections are always parsed.
func ParseConfig(repData []byte, cfg Config) (*rep.Replay, error) {
	dec := repdecoder.New(repData)
	defer dec.Close()

	return parseProtected(dec, cfg)
}

// parseProtected calls parse(), but protects the function call from panics,
// in which case it returns ErrParsing.
func parseProtected(dec repdecoder.Decoder, cfg Config) (r *rep.Replay, err error) {
	// Make sure cfg.Logger is not nil, in one place:
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	// Input is untrusted data, protect the parsing logic.
	// It also protects against implementation bugs.
	defer func() {
		if r := recover(); r != nil {
			cfg.Logger.Printf("Parsing error: %v", r)
			buf := make([]byte, 2000)
			n := runtime.Stack(buf, false)
			cfg.Logger.Printf("Stack: %s", buf[:n])
			err = ErrParsing
		}
	}()

	return parse(dec, cfg)
}

// Section describes a Section of the replay.
type Section struct {
	// ID of the section
	ID int

	// Size of the uncompressed section in bytes;
	// 0 means the Size has to be read as a section of 4 bytes
	Size int32

	// ParseFunc defines the function responsible to process (parse / interpret)
	// the section's data.
	ParseFunc func(data []byte, r *rep.Replay, cfg Config) error

	// Optional section string ID
	StrID string
}

// Sections describes the subsequent Sections of replays
var Sections = []*Section{
	{ID: 0, Size: 0x04, ParseFunc: parseReplayID},
	{ID: 1, Size: 0x279, ParseFunc: parseHeader},
	{ID: 2, Size: 0, ParseFunc: parseCommands},
	{ID: 3, Size: 0, ParseFunc: parseMapData},
	{ID: 4, Size: 0x300, ParseFunc: parsePlayerNames},
}

// ModernSections holds custom sections added in Remastered, and also custom sections
// added by 3rd party vendors.
var ModernSections = map[int32]*Section{
	1313426259: {ID: 5, Size: 0x15e0, ParseFunc: parseSkin, StrID: "SKIN"},
	1398033740: {ID: 6, Size: 0x1c, ParseFunc: parseLmts, StrID: "LMTS"},
	1481197122: {ID: 7, Size: 0x08, ParseFunc: parseBfix, StrID: "BFIX"},
	1380729667: {ID: 8, Size: 0xc0, ParseFunc: parsePlayerColors, StrID: "CCLR"},
	1195787079: {ID: 9, Size: 0x19, ParseFunc: parseGcfg, StrID: "GCFG"},

	// ShieldBattery's custom section
	1952539219: {ID: 10, Size: 0, ParseFunc: parseShieldBatterySection, StrID: "Sbat"},
}

// Named sections
var (
	SectionReplayID    = Sections[0]
	SectionHeader      = Sections[1]
	SectionCommands    = Sections[2]
	SectionMapData     = Sections[3]
	SectionPlayerNames = Sections[4]
)

// parse parses an SC:BW replay using the given Decoder.
func parse(dec repdecoder.Decoder, cfg Config) (*rep.Replay, error) {
	r := new(rep.Replay)
	r.RepFormat = dec.RepFormat()

	// We have to read all sections, some data (e.g. player colors) are positioned after map data.

	// A replay is a sequence of sections:
	for sectionCounter := 0; ; sectionCounter++ {
		if err := dec.NewSection(); err != nil {
			if err == repdecoder.ErrNoMoreSections {
				break
			}
			return nil, fmt.Errorf("Decoder.NewSection() error: %w", err)
		}

		var s *Section
		var size int32
		if sectionCounter < len(Sections) {
			s = Sections[sectionCounter]

			// Determine section size:
			size = s.Size
			if size == 0 {
				sizeData, _, err := dec.Section(4)
				if err != nil {
					return nil, fmt.Errorf("Decoder.Section() error when reading size: %w", err)
				}
				size = int32(binary.LittleEndian.Uint32(sizeData))
			}
		}

		// Read section data
		data, sectionID, err := dec.Section(size)
		if err != nil {
			if s != nil && s.ID == SectionReplayID.ID {
				err = ErrNotReplayFile // In case of Replay ID section return special error
			}
			if err == io.EOF {
				break // New sections with StrID are optional
			}
			if sectionCounter >= len(Sections) {
				// If we got "enough" info, just log the error:
				cfg.Logger.Printf("Warning: Decoder.Section() error: %v", err)
				break
			}
			return nil, fmt.Errorf("Decoder.Section() error: %w", err)
		}

		if s == nil {
			s = ModernSections[sectionID]
			if s == nil {
				// Unknown section, just skip it:
				idBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(idBytes, uint32(sectionID))
				cfg.Logger.Printf("Unknown modern section ID: %s", idBytes)
				continue
			}
		}

		// Need to process?
		switch {
		case s == SectionCommands && !cfg.Commands:
		case s == SectionMapData && !cfg.MapData:
		default:
			// Process section data
			if err = s.ParseFunc(data, r, cfg); err != nil {
				return nil, fmt.Errorf("ParseFunc() error (sectionID: %d): %v", s.ID, err)
			}
		}
	}

	// Modern sections may or may not exist. Remastered's modern sections are in fixed order,
	// but we don't rely on it.

	return r, nil
}

// repIDs is the possible valid content of the Replay ID section
var repIDs = [][]byte{
	[]byte("seRS"), // Starting from 1.21
	[]byte("reRS"), // Up until 1.20.
}

// parseReplayID processes the replay ID data.
func parseReplayID(data []byte, r *rep.Replay, cfg Config) (err error) {
	for _, repID := range repIDs {
		if bytes.Equal(data, repID) {
			return
		}
	}

	return ErrNotReplayFile
}

var headerFields = []*rep.DebugFieldDescriptor{
	{Offset: 0x00, Length: 1, Name: "Engine"},
	{Offset: 0x01, Length: 4, Name: "Frames"},
	{Offset: 0x08, Length: 8, Name: "Start time"},
	{Offset: 0x18, Length: 28, Name: "Title"},
	{Offset: 0x34, Length: 2, Name: "Map width"},
	{Offset: 0x36, Length: 2, Name: "Map height"},
	{Offset: 0x39, Length: 1, Name: "Available slots count"},
	{Offset: 0x3a, Length: 1, Name: "Speed"},
	{Offset: 0x3c, Length: 2, Name: "Type"},
	{Offset: 0x3e, Length: 2, Name: "SubType"},
	{Offset: 0x48, Length: 24, Name: "Host"},
	{Offset: 0x61, Length: 26, Name: "Map"},
	{Offset: 0xa1, Length: 432, Name: "Player structs (12)"},
	{Offset: 0xa1, Length: 36, Name: "Player 1 struct"},
	{Offset: 0xa1, Length: 2, Name: "Player 1 slot ID"},
	{Offset: 0xa1 + 4, Length: 1, Name: "Player 1 ID"},
	{Offset: 0xa1 + 8, Length: 1, Name: "Player 1 type"},
	{Offset: 0xa1 + 9, Length: 1, Name: "Player 1 race"},
	{Offset: 0xa1 + 10, Length: 1, Name: "Player 1 team"},
	{Offset: 0xa1 + 11, Length: 25, Name: "Player 1 name"},
	{Offset: 0xa1 + 36, Length: 36, Name: "Player 2 struct"},
	{Offset: 0x251, Length: 8 * 4, Name: "Player colors (8)"},
	{Offset: 0x251, Length: 4, Name: "Player 1 color"},
	{Offset: 0x251 + 4, Length: 4, Name: "Player 2 color"},
}

// parseHeader processes the replay header data.
func parseHeader(data []byte, r *rep.Replay, cfg Config) error {
	bo := binary.LittleEndian // ByteOrder reader: little-endian

	h := new(rep.Header)
	r.Header = h
	if cfg.Debug {
		h.Debug = &rep.HeaderDebug{
			Data:   data,
			Fields: headerFields,
		}
	}

	// Fill Version:
	switch r.RepFormat {
	case repdecoder.RepFormatModern121:
		r.Header.Version = "1.21+"
	case repdecoder.RepFormatLegacy:
		r.Header.Version = "-1.16"
	case repdecoder.RepFormatModern:
		r.Header.Version = "1.18-1.20"
	}

	h.Engine = repcore.EngineByID(data[0x00])
	h.Frames = repcore.Frame(bo.Uint32(data[0x01:]))
	h.StartTime = time.Unix(int64(bo.Uint32(data[0x08:])), 0) // replay stores seconds since EPOCH
	// SC:R uses UTF-8 always (except the map data section which may come from an external source or from the "past").
	// The game UI allows longer title than what fits into its space in the header. If longer, SC simply "cuts" it,
	// even in the middle of a multi-byte UTF-8 sequence :S
	// This may result in reading invalid UTF-8 title data, even though it was generated using UTF-8,
	// and hence must be decoded as such.
	if r.RepFormat == repdecoder.RepFormatLegacy {
		h.Title, h.RawTitle = cString(data[0x18 : 0x18+28])
	} else {
		h.Title, h.RawTitle = cStringUTF8(data[0x18 : 0x18+28])
	}
	h.MapWidth = bo.Uint16(data[0x34:])
	h.MapHeight = bo.Uint16(data[0x36:])
	h.AvailSlotsCount = data[0x39]
	h.Speed = repcore.SpeedByID(data[0x3a])
	h.Type = repcore.GameTypeByID(bo.Uint16(data[0x3c:]))
	h.SubType = bo.Uint16(data[0x3e:])
	h.Host, h.RawHost = cString(data[0x48 : 0x48+24])
	h.Map, h.RawMap = cString(data[0x61 : 0x61+26])

	// Parse players
	const (
		slotsCount = 12
		maxPlayers = 8
	)
	h.PIDPlayers = make(map[byte]*rep.Player, slotsCount)
	h.Slots = make([]*rep.Player, slotsCount)
	playerStructs := data[0xa1 : 0xa1+432]
	for i := range h.Slots {
		p := new(rep.Player)
		h.Slots[i] = p
		ps := playerStructs[i*36 : i*36+432/slotsCount]
		p.SlotID = bo.Uint16(ps)
		p.ID = ps[4]
		p.Type = repcore.PlayerTypeByID(ps[8])
		p.Race = repcore.RaceByID(ps[9])
		p.Team = ps[10]
		p.Name, p.RawName = cString(ps[11 : 11+25])

		if i < maxPlayers {
			p.Color = repcore.ColorByID(bo.Uint32(data[0x251+i*4:]))
		}

		// Filter real players:
		if p.Name != "" {
			h.OrigPlayers = append(h.OrigPlayers, p)
			h.PIDPlayers[p.ID] = p
		}
	}

	// If game type is melee or OneOnOne, all players' teams may be set to 0 or 1.
	// Heuristic improvements: If 2 players only and their teams are the same, change teams to 1 and 2,
	// and so matchup will be e.g. ZvT instead of ZT,
	// and winner detection can also work (because teams will be different).
	if (h.Type == repcore.GameTypeMelee || h.Type == repcore.GameType1on1) && len(h.OrigPlayers) == 2 &&
		h.OrigPlayers[0].Team == h.OrigPlayers[1].Team {
		h.OrigPlayers[0].Team = 1
		h.OrigPlayers[1].Team = 2
	}
	// Also if game type is FFA, teams are set to 0.
	// Assign teams incrementing from 1.
	if h.Type == repcore.GameTypeFFA {
		for i, p := range h.OrigPlayers {
			p.Team = byte(i + 1)
		}
	}

	// Fill Players in team order:
	h.Players = make([]*rep.Player, len(h.OrigPlayers))
	copy(h.Players, h.OrigPlayers)
	sort.SliceStable(h.Players, func(i int, j int) bool {
		return h.Players[i].Team < h.Players[j].Team
	})

	return nil
}

// parseCommands processes the players' commands data.
func parseCommands(data []byte, r *rep.Replay, cfg Config) error {
	bo := binary.LittleEndian // ByteOrder reader: little-endian

	_ = bo
	cs := new(rep.Commands)
	r.Commands = cs
	if cfg.Debug {
		cs.Debug = &rep.CommandsDebug{Data: data}
	}

	for sr, size := (sliceReader{b: data}), uint32(len(data)); sr.pos < size; {
		frame := sr.getUint32()

		// Command block in this frame
		cmdBlockSize := sr.getByte()                    // cmd block size (remaining)
		cmdBlockEndPos := sr.pos + uint32(cmdBlockSize) // Cmd block end position

		for sr.pos < cmdBlockEndPos {
			parseOk := true

			var cmd repcmd.Cmd
			base := &repcmd.Base{
				Frame: repcore.Frame(frame),
			}
			base.PlayerID = sr.getByte()
			base.Type = repcmd.TypeByID(sr.getByte())

			switch base.Type.ID { // Try to list in frequency order:

			case repcmd.TypeIDRightClick:
				rccmd := &repcmd.RightClickCmd{Base: base}
				rccmd.Pos.X = sr.getUint16()
				rccmd.Pos.Y = sr.getUint16()
				rccmd.UnitTag = repcmd.UnitTag(sr.getUint16())
				rccmd.Unit = repcmd.UnitByID(sr.getUint16())
				rccmd.Queued = sr.getByte() != 0
				cmd = rccmd

			case repcmd.TypeIDSelect, repcmd.TypeIDSelectAdd, repcmd.TypeIDSelectRemove:
				count := sr.getByte()
				selectCmd := &repcmd.SelectCmd{
					Base:     base,
					UnitTags: make([]repcmd.UnitTag, count),
				}
				for i := byte(0); i < count; i++ {
					selectCmd.UnitTags[i] = repcmd.UnitTag(sr.getUint16())
				}
				cmd = selectCmd

			case repcmd.TypeIDHotkey:
				hotkeyCmd := &repcmd.HotkeyCmd{Base: base}
				hotkeyCmd.HotkeyType = repcmd.HotkeyTypeByID(sr.getByte())
				hotkeyCmd.Group = sr.getByte()
				cmd = hotkeyCmd

			case repcmd.TypeIDTrain, repcmd.TypeIDUnitMorph:
				cmd = &repcmd.TrainCmd{
					Base: base,
					Unit: repcmd.UnitByID(sr.getUint16()),
				}

			case repcmd.TypeIDTargetedOrder:
				tocmd := &repcmd.TargetedOrderCmd{Base: base}
				tocmd.Pos.X = sr.getUint16()
				tocmd.Pos.Y = sr.getUint16()
				tocmd.UnitTag = repcmd.UnitTag(sr.getUint16())
				tocmd.Unit = repcmd.UnitByID(sr.getUint16())
				tocmd.Order = repcmd.OrderByID(sr.getByte())
				tocmd.Queued = sr.getByte() != 0
				cmd = tocmd

			case repcmd.TypeIDBuild:
				buildCmd := &repcmd.BuildCmd{Base: base}
				buildCmd.Order = repcmd.OrderByID(sr.getByte())
				buildCmd.Pos.X = sr.getUint16()
				buildCmd.Pos.Y = sr.getUint16()
				buildCmd.Unit = repcmd.UnitByID(sr.getUint16())
				if buildCmd.Order.ID == repcmd.OrderIDBuildingLand {
					// It's actually a Land command:
					landCmd := (*repcmd.LandCmd)(buildCmd) // Fields are identical, we may simply convert it
					landCmd.Base.Type = repcmd.TypeLand
					cmd = landCmd
				} else {
					// It's truly a build command
					cmd = buildCmd
				}

			case repcmd.TypeIDStop, repcmd.TypeIDBurrow, repcmd.TypeIDUnburrow,
				repcmd.TypeIDReturnCargo, repcmd.TypeIDHoldPosition, repcmd.TypeIDUnloadAll,
				repcmd.TypeIDUnsiege, repcmd.TypeIDSiege, repcmd.TypeIDCloack, repcmd.TypeIDDecloack:
				cmd = &repcmd.QueueableCmd{
					Base:   base,
					Queued: sr.getByte() != 0,
				}

			case repcmd.TypeIDLeaveGame:
				cmd = &repcmd.LeaveGameCmd{
					Base:   base,
					Reason: repcmd.LeaveReasonByID(sr.getByte()),
				}

			case repcmd.TypeIDMinimapPing:
				pingCmd := &repcmd.MinimapPingCmd{Base: base}
				pingCmd.Pos.X = sr.getUint16()
				pingCmd.Pos.Y = sr.getUint16()
				cmd = pingCmd

			case repcmd.TypeIDChat:
				chatCmd := &repcmd.ChatCmd{Base: base}
				chatCmd.SenderSlotID = sr.getByte()
				chatCmd.Message, _ = cString(sr.readSlice(80))
				cmd = chatCmd

			case repcmd.TypeIDVision:
				data := sr.getUint16()
				visionCmd := &repcmd.VisionCmd{
					Base: base,
				}
				// There is 1 bit for each slot, 0x01: shared vision for that slot
				for i := byte(0); i < 12; i++ {
					if data&0x01 != 0 {
						visionCmd.SlotIDs = append(visionCmd.SlotIDs, i)
					}
					data >>= 1
				}
				cmd = visionCmd

			case repcmd.TypeIDAlliance:
				data := sr.getUint32()
				allianceCmd := &repcmd.AllianceCmd{
					Base: base,
				}
				// There are 2 bits for each slot, 0x00: not allied, 0x1: allied, 0x02: allied victory
				for i := byte(0); i < 11; i++ { // only 11 slots, 12th is always 0x01 or 0x02
					if x := data & 0x03; x != 0 {
						allianceCmd.SlotIDs = append(allianceCmd.SlotIDs, i)
						if x == 2 {
							allianceCmd.AlliedVictory = true
						}
					}
					data >>= 2
				}
				cmd = allianceCmd

			case repcmd.TypeIDGameSpeed:
				cmd = &repcmd.GameSpeedCmd{
					Base:  base,
					Speed: repcore.SpeedByID(sr.getByte()),
				}

			case repcmd.TypeIDCancelTrain:
				cmd = &repcmd.CancelTrainCmd{
					Base:    base,
					UnitTag: repcmd.UnitTag(sr.getUint16()),
				}

			case repcmd.TypeIDUnload:
				cmd = &repcmd.UnloadCmd{
					Base:    base,
					UnitTag: repcmd.UnitTag(sr.getUint16()),
				}

			case repcmd.TypeIDLiftOff:
				liftOffCmd := &repcmd.LiftOffCmd{Base: base}
				liftOffCmd.Pos.X = sr.getUint16()
				liftOffCmd.Pos.Y = sr.getUint16()
				cmd = liftOffCmd

			case repcmd.TypeIDTech:
				cmd = &repcmd.TechCmd{
					Base: base,
					Tech: repcmd.TechByID(sr.getByte()),
				}

			case repcmd.TypeIDUpgrade:
				cmd = &repcmd.UpgradeCmd{
					Base:    base,
					Upgrade: repcmd.UpgradeByID(sr.getByte()),
				}

			case repcmd.TypeIDBuildingMorph:
				cmd = &repcmd.BuildingMorphCmd{
					Base: base,
					Unit: repcmd.UnitByID(sr.getUint16()),
				}

			case repcmd.TypeIDLatency:
				cmd = &repcmd.LatencyCmd{
					Base:    base,
					Latency: repcmd.LatencyTypeByID(sr.getByte()),
				}

			case repcmd.TypeIDCheat:
				cmd = &repcmd.GeneralCmd{
					Base: base,
					Data: sr.readSlice(4),
				}

			case repcmd.TypeIDSaveGame, repcmd.TypeIDLoadGame:
				count := sr.getUint32()
				sr.pos += count

			// NO ADDITIONAL DATA:

			case repcmd.TypeIDKeepAlive:
			case repcmd.TypeIDRestartGame:
			case repcmd.TypeIDPause:
			case repcmd.TypeIDResume:
			case repcmd.TypeIDCancelBuild:
			case repcmd.TypeIDCancelMorph:
			case repcmd.TypeIDCarrierStop:
			case repcmd.TypeIDReaverStop:
			case repcmd.TypeIDOrderNothing:
			case repcmd.TypeIDTrainFighter:
			case repcmd.TypeIDMergeArchon:
			case repcmd.TypeIDCancelNuke:
			case repcmd.TypeIDCancelTech:
			case repcmd.TypeIDCancelUpgrade:
			case repcmd.TypeIDCancelAddon:
			case repcmd.TypeIDStim:
			case repcmd.TypeIDVoiceEnable:
			case repcmd.TypeIDVoiceDisable:
			case repcmd.TypeIDStartGame:
			case repcmd.TypeIDBriefingStart:
			case repcmd.TypeIDMergeDarkArchon:
			case repcmd.TypeIDMakeGamePublic:

			// DON'T CARE COMMANDS:

			case repcmd.TypeIDSync:
				sr.pos += 6
			case repcmd.TypeIDVoiceSquelch:
				sr.pos++
			case repcmd.TypeIDVoiceUnsquelch:
				sr.pos++
			case repcmd.TypeIDDownloadPercentage:
				sr.pos++
			case repcmd.TypeIDChangeGameSlot:
				sr.pos += 5
			case repcmd.TypeIDNewNetPlayer:
				sr.pos += 7
			case repcmd.TypeIDJoinedGame:
				sr.pos += 17
			case repcmd.TypeIDChangeRace:
				sr.pos += 2
			case repcmd.TypeIDTeamGameTeam:
				sr.pos++
			case repcmd.TypeIDUMSTeam:
				sr.pos++
			case repcmd.TypeIDMeleeTeam:
				sr.pos += 2
			case repcmd.TypeIDSwapPlayers:
				sr.pos += 2
			case repcmd.TypeIDSavedData:
				sr.pos += 12
			case repcmd.TypeIDReplaySpeed:
				sr.pos += 9

			// New commands introduced in 1.21

			case repcmd.TypeIDRightClick121:
				rccmd := &repcmd.RightClickCmd{Base: base}
				rccmd.Pos.X = sr.getUint16()
				rccmd.Pos.Y = sr.getUint16()
				rccmd.UnitTag = repcmd.UnitTag(sr.getUint16())
				sr.getUint16() // Unknown, always 0?
				rccmd.Unit = repcmd.UnitByID(sr.getUint16())
				rccmd.Queued = sr.getByte() != 0
				cmd = rccmd

			case repcmd.TypeIDTargetedOrder121:
				tocmd := &repcmd.TargetedOrderCmd{Base: base}
				tocmd.Pos.X = sr.getUint16()
				tocmd.Pos.Y = sr.getUint16()
				tocmd.UnitTag = repcmd.UnitTag(sr.getUint16())
				sr.getUint16() // Unknown, always 0?
				tocmd.Unit = repcmd.UnitByID(sr.getUint16())
				tocmd.Order = repcmd.OrderByID(sr.getByte())
				tocmd.Queued = sr.getByte() != 0
				cmd = tocmd

			case repcmd.TypeIDUnload121:
				ucmd := &repcmd.UnloadCmd{Base: base}
				ucmd.UnitTag = repcmd.UnitTag(sr.getUint16())
				sr.getUint16() // Unknown, always 0?
				cmd = ucmd

			case repcmd.TypeIDSelect121, repcmd.TypeIDSelectAdd121, repcmd.TypeIDSelectRemove121:
				count := sr.getByte()
				selectCmd := &repcmd.SelectCmd{
					Base:     base,
					UnitTags: make([]repcmd.UnitTag, count),
				}
				for i := byte(0); i < count; i++ {
					selectCmd.UnitTags[i] = repcmd.UnitTag(sr.getUint16())
					sr.getUint16() // Unknown, always 0?
				}
				cmd = selectCmd

			default:
				// We don't know how to parse this command, we have to skip
				// to the end of the command block
				// (potentially skipping additional commands...)
				var remBytes []byte
				if sr.pos <= cmdBlockEndPos && cmdBlockEndPos <= uint32(len(sr.b)) { // Due to "bad" parsing these must be checked...
					remBytes = sr.b[sr.pos:cmdBlockEndPos]
				}
				cfg.Logger.Printf("skipping typeID: %#v, frame: %d, playerID: %d, remaining bytes: %d [% x]\n", base.Type.ID, base.Frame, base.PlayerID, cmdBlockEndPos-sr.pos, remBytes)
				pec := &repcmd.ParseErrCmd{Base: base}
				if len(cs.Cmds) > 0 {
					pec.PrevCmd = cs.Cmds[len(cs.Cmds)-1]
				}
				cs.ParseErrCmds = append(cs.ParseErrCmds, pec)
				sr.pos = cmdBlockEndPos
				parseOk = false
			}

			if parseOk {
				if cmd == nil {
					cs.Cmds = append(cs.Cmds, base)
				} else {
					cs.Cmds = append(cs.Cmds, cmd)
				}
			}
		}

		sr.pos = cmdBlockEndPos
	}

	return nil
}

// parseMapData processes the map data data.
func parseMapData(data []byte, r *rep.Replay, cfg Config) error {
	md := new(rep.MapData)
	r.MapData = md
	if cfg.Debug {
		md.Debug = &rep.MapDataDebug{Data: data}
	}
	if cfg.MapGraphics {
		md.MapGraphics = &rep.MapGraphics{}
	}

	// Even though "ERA " section is mandatory, I've seen reps where it was missing.
	// TileSet may be cruitial for some apps, let's ensure it doesn't remain nil.
	// Somewhat arbitrary default:
	md.TileSet = repcore.TileSetTwilight
	md.TileSetMissing = true

	var (
		scenarioNameIdx        uint16 // String index
		scenarioDescriptionIdx uint16 // String index
		stringsData            []byte
		extendedStringsData    bool
	)

	// Map data section is a sequence of sub-sections:
	for sr, size := (sliceReader{b: data}), uint32(len(data)); sr.pos < size; {
		id := sr.getString(4)
		// Seen examples where a "final" UPUS section following UPRP section had only 1 byte hereon, so check:
		if sr.pos+4 >= size {
			break
		}
		ssSize := sr.getUint32()    // sub-section size (remaining)
		ssEndPos := sr.pos + ssSize // sub-section end position

		switch id {
		case "VER ":
			md.Version = sr.getUint16()
		case "ERA ": // Tile set sub-section
			md.TileSet = repcore.TileSetByID(sr.getUint16() & 0x07)
			md.TileSetMissing = false
		case "DIM ": // Dimension sub-section
			// If map has a non-standard size, the replay header contains
			// invalid map size, this is the correct one.
			width := sr.getUint16()
			height := sr.getUint16()
			if width <= 256 && height <= 256 {
				if width > r.Header.MapWidth {
					r.Header.MapWidth = width
				}
				if height > r.Header.MapHeight {
					r.Header.MapHeight = height
				}
			}
		case "OWNR": // StarCraft Player Types
			count := uint32(12) // 12 bytes, 1 for each player
			if count > ssSize {
				count = ssSize
			}
			owners := sr.readSlice(count)
			md.PlayerOwners = make([]*repcore.PlayerOwner, len(owners))
			for i, id := range owners {
				md.PlayerOwners[i] = repcore.PlayerOwnerByID(id)
			}
		case "SIDE": // Player races
			count := uint32(12) // 12 bytes, 1 for each player
			if count > ssSize {
				count = ssSize
			}
			sides := sr.readSlice(count)
			md.PlayerSides = make([]*repcore.PlayerSide, len(sides))
			for i, id := range sides {
				md.PlayerSides[i] = repcore.PlayerSideByID(id)
			}
		case "MTXM": // Tile sub-section
			// map_width*map_height (a tile is an uint16 value)
			maxI := ssSize / 2
			// Note: Sometimes map is broken into multiple sections.
			// The first one is the biggest (whole map size),
			// but the beginning of map is empty. The subsequent MTXM
			// sub-sections will fill the whole at the beginning.
			// An example was found when the first MTXM section was only
			// 8 elements, and the next was the whole map, beginning also filled.
			// Therefore if currently allocated Tile is small, a new one is allocated.
			if len(md.Tiles) < int(maxI) {
				md.Tiles = make([]uint16, maxI)
			}
			for i := uint32(0); i < maxI; i++ {
				md.Tiles[i] = sr.getUint16()
			}
		case "UNIT": // Placed units
			for sr.pos+36 <= ssEndPos { // Loop while we have a complete unit
				unitEndPos := sr.pos + 36 // 36 bytes for each unit

				sr.pos += 4 // uint32 unit class instance ("serial number")
				x := sr.getUint16()
				y := sr.getUint16()
				unitID := sr.getUint16()
				sr.pos += 2                 // uint16 Type of relation to another building (i.e. add-on, nydus link)
				sr.pos += 2                 // uint16 Flags of special properties (e.g. cloacked, burrowed etc.)
				sr.pos += 2                 // uint16 valid elements flag
				ownerID := sr.getByte()     // 0-based SlotID
				sr.pos++                    // Hit points % (1-100)
				sr.pos++                    // Shield points % (1-100)
				sr.pos++                    // Energy points % (1-100)
				resAmount := sr.getUint32() // Resource amount

				switch unitID {
				case repcmd.UnitIDMineralField1, repcmd.UnitIDMineralField2, repcmd.UnitIDMineralField3:
					md.MineralFields = append(md.MineralFields, rep.Resource{Point: repcore.Point{X: x, Y: y}, Amount: resAmount})
				case repcmd.UnitIDVespeneGeyser:
					md.Geysers = append(md.Geysers, rep.Resource{Point: repcore.Point{X: x, Y: y}, Amount: resAmount})
				case repcmd.UnitIDStartLocation:
					md.StartLocations = append(md.StartLocations,
						rep.StartLocation{Point: repcore.Point{X: x, Y: y}, SlotID: ownerID},
					)
				}

				if cfg.MapGraphics {
					md.MapGraphics.PlacedUnits = append(md.MapGraphics.PlacedUnits, &rep.PlacedUnit{
						Point:          repcore.Point{X: x, Y: y},
						UnitID:         unitID,
						SlotID:         ownerID,
						ResourceAmount: resAmount,
					})
				}

				// Skip unprocessed unit data:
				sr.pos = unitEndPos
			}
		case "THG2": // StarCraft Sprites
			if cfg.MapGraphics {
				for sr.pos+10 <= ssEndPos { // Loop while we have a complete sprite
					spriteEndPos := sr.pos + 10 // 10 bytes for each sprite

					spriteID := sr.getUint16()
					x := sr.getUint16()
					y := sr.getUint16()
					ownerID := sr.getByte() // 0-based SlotID
					sr.pos++                // Unused
					flags := sr.getUint16()
					if flags&0x1000 == 0 {
						// It's actually a unit
						md.MapGraphics.PlacedUnits = append(md.MapGraphics.PlacedUnits, &rep.PlacedUnit{
							Point:  repcore.Point{X: x, Y: y},
							UnitID: spriteID,
							SlotID: ownerID,
							Sprite: true,
						})
					} else {
						// It really is a sprite
						md.MapGraphics.Sprites = append(md.MapGraphics.Sprites, &rep.Sprite{
							Point:    repcore.Point{X: x, Y: y},
							SpriteID: spriteID,
						})
					}

					// Skip unprocessed sprite data:
					sr.pos = spriteEndPos
				}
			}
		case "SPRP": // Scenario properties
			// Strings section might be after this, so we just record the string indices for now:
			scenarioNameIdx = sr.getUint16()
			scenarioDescriptionIdx = sr.getUint16()
		case "STR ": // String data
			// There might be multiple "STR " sections, subsequent sections overwrite the
			// beginning of earlier sections.
			stringsStart := int(sr.pos)
			// count := sr.getUint16() // Number of following offsets (uint16 values)
			if len(stringsData) < int(ssEndPos)-stringsStart {
				stringsData = make([]byte, int(ssEndPos)-stringsStart)
			}
			copy(stringsData, data[stringsStart:ssEndPos])
		case "STRx": // Extended String data
			// This section is identical to "STR " except that all uint16 values are uint32 values.
			stringsStart := int(sr.pos)
			// count := sr.getUint32() // Number of following offsets (uint32 values)
			if len(stringsData) < int(ssEndPos)-stringsStart {
				stringsData = make([]byte, int(ssEndPos)-stringsStart)
			}
			copy(stringsData, data[stringsStart:ssEndPos])
			extendedStringsData = true
		}

		// Part or all of the sub-section might be unprocessed, skip the unprocessed bytes
		sr.pos = ssEndPos
	}

	// Get a string from the strings identified by its index.
	getString := func(idx uint16) string {
		if idx == 0 {
			return ""
		}
		var offsetSize uint32
		if extendedStringsData {
			offsetSize = 4
		} else {
			offsetSize = 2
		}
		pos := uint32(idx) * offsetSize // idx is 1-based (0th offset is not included), but stringsData contains the offsets count too
		if int(pos+offsetSize-1) >= len(stringsData) {
			cfg.Logger.Printf("Invalid strings index: %d, map: %s", idx, r.Header.Map)
			return ""
		}
		var offset uint32
		if extendedStringsData {
			offset = (&sliceReader{b: stringsData, pos: pos}).getUint32()
		} else {
			offset = uint32((&sliceReader{b: stringsData, pos: pos}).getUint16())
		}
		if int(offset) >= len(stringsData) {
			cfg.Logger.Printf("Invalid strings offset: %d, strings index: %d, map: %s", offset, idx, r.Header.Map)
			return ""
		}
		s, _ := cString(stringsData[offset:])
		return s
	}

	md.Name = getString(scenarioNameIdx)
	md.Description = getString(scenarioDescriptionIdx)

	return nil
}

// parsePlayerNames processes the player names data.
func parsePlayerNames(data []byte, r *rep.Replay, cfg Config) error {
	// Note: these player names parse well even when decoding is unknown in header
	// (are these always UTF-8?)
	for i, p := range r.Header.Slots {
		pos := i * 96
		if pos+96 > len(data) {
			break
		}

		if p.Type != repcore.PlayerTypeInactive {
			name, orig := cString(data[pos : pos+96])
			if name != "" {
				p.Name, p.RawName = name, orig
			}
		}
	}

	return nil
}

// parseSkin processes the skin data.
func parseSkin(data []byte, r *rep.Replay, cfg Config) error {
	// TODO 0x15e0 bytes of data
	return nil
}

// parseLmts processes the lmts data.
func parseLmts(data []byte, r *rep.Replay, cfg Config) error {
	// TODO 0x1c bytes of data

	// bo := binary.LittleEndian // ByteOrder reader: little-endian
	// bo.Uint32(data[0x0:])     // Images limit
	// bo.Uint32(data[0x4:])     // Sprites limit
	// bo.Uint32(data[0x8:])     // Lone limit
	// bo.Uint32(data[0x0c:])    // Units limit
	// bo.Uint32(data[0x10:])    // Bullets limit
	// bo.Uint32(data[0x14:])    // Orders limit
	// bo.Uint32(data[0x18:])    // Fog sprites limit

	return nil
}

// parseBfix processes the bfix data.
func parseBfix(data []byte, r *rep.Replay, cfg Config) error {
	// TODO 0x08 bytes of data
	return nil
}

// parseGcfg processes the gcfg data.
func parseGcfg(data []byte, r *rep.Replay, cfg Config) error {
	// TODO 0x19 bytes of data
	return nil
}

// parsePlayerColors processes the player colors data.
func parsePlayerColors(data []byte, r *rep.Replay, cfg Config) error {
	// 16 bytes footprint for all colors.
	for i, p := range r.Header.Slots {
		pos := i * 16
		if pos+16 > len(data) {
			break
		}
		if c := repcore.ColorByFootprint(data[pos : pos+16]); c != nil {
			p.Color = c
		}
	}

	return nil
}

// parseShieldBatterySection processes the ShieldBattery data.
func parseShieldBatterySection(data []byte, r *rep.Replay, cfg Config) error {
	// info source:
	// https://github.com/ShieldBattery/ShieldBattery/blob/master/game/src/replay.rs#L62-L80
	// https://github.com/ShieldBattery/ShieldBattery/blob/master/app/replays/parse-shieldbattery-replay.ts

	if len(data) < 0x56 {
		// 0x56 bytes is the size of SB's first version of the section.
		return nil // Unknown format
	}

	bo := binary.LittleEndian // ByteOrder reader: little-endian

	sb := new(rep.ShieldBattery)
	r.ShieldBattery = sb

	formatVersion := bo.Uint16(data)

	sb.StarCraftExeBuild = bo.Uint32(data[0x01:])
	sb.ShieldBatteryVersion, _ = cString(data[0x06:0x16])

	// 0x16 - 0x1a: team_game_main_players
	// 0x1a - 0x26: starting_races

	gameID := data[0x26:0x36]
	sb.GameID = fmt.Sprintf("%x-%x-%x-%x-%x", gameID[:4], gameID[4:6], gameID[6:8], gameID[8:10], gameID[10:])

	if formatVersion >= 0x01 {
		// 0x56 - 0x58: game_logic_version
	}

	return nil
}

var koreanDecoder = korean.EUCKR.NewDecoder()

// cString returns a 0x00 byte terminated string from the given buffer.
// If the string is not valid UTF-8, tries to decode it as EUC-KR (also known as Code Page 949).
// Returns both the decoded and the original string.
func cString(data []byte) (s string, orig string) {
	// Find 0x00 byte:
	for i, ch := range data {
		if ch == 0 {
			data = data[:i] // excludes terminating 0x00

			if !utf8.Valid(data) {
				// Try korean
				if krdata, err := koreanDecoder.Bytes(data); err == nil {
					return string(krdata), string(data)
				}
			}
			break // Either UTF-8 or custom decoding failed
		}
	}

	// Return data as string.
	// We end up here if:
	//   - no terminating 0 char found,
	//   - or string is valid UTF-8,
	//   - or it is invalid but custom decoding failed
	// Either way:
	s = string(data)
	return s, s
}

// cStringUTF8 returns a 0x00 byte terminated string from the given buffer,
// always using UTF-8 encoding.
// If the data is invalid UTF-8, invalid sequences will be removed from it.
//
// Returns both the decoded and the original string.
func cStringUTF8(data []byte) (s string, orig string) {
	// Find 0x00 byte:
	for i, ch := range data {
		if ch == 0 {
			data = data[:i] // excludes terminating 0x00
			break
		}
	}

	if !utf8.Valid(data) {
		return string(bytes.ToValidUTF8(data, nil)), string(data)
	}

	s = string(data)
	return s, s
}
