/*

A simple CLI app to parse and display information about
a StarCraft: Brood War replay passed as a CLI argument.

*/
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"os"
	"runtime"
	"strings"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
)

const (
	appName    = "screp"
	appVersion = "v1.7.1"
	appAuthor  = "Andras Belicza"
	appHome    = "https://github.com/icza/screp"
)

const (
	ExitCodeMissingArguments         = 1
	ExitCodeFailedToParseReplay      = 2
	ExitCodeFailedToCreateOutputFile = 3
	ExitCodeInvalidMapDataHash       = 4
)

const validMapDataHashes = "valid values are 'sha1', 'sha256', 'sha512', 'md5'"

// Flag variables
var (
	version = flag.Bool("version", false, "print version info and exit")

	header      = flag.Bool("header", true, "print replay header")
	mapData     = flag.Bool("map", false, "print map data")
	mapTiles    = flag.Bool("maptiles", false, "print map data tiles; valid with 'map'")
	mapResLoc   = flag.Bool("mapres", false, "print map data resource locations (minerals and geysers); valid with 'map'")
	cmds        = flag.Bool("cmds", false, "print player commands")
	computed    = flag.Bool("computed", true, "print computed / derived data")
	mapDataHash = flag.String("mapDataHash", "", "calculate and print the hash of map data section too using the given algorithm;\n"+validMapDataHashes)
	dumpMapData = flag.Bool("dumpMapData", false, "dump the raw map data (CHK) instead of JSON replay info\nuse it with the 'outfile' flag")
	outFile     = flag.String("outfile", "", "optional output file name")

	indent = flag.Bool("indent", true, "use indentation when formatting output")
)

func main() {
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(ExitCodeMissingArguments)
	}

	cfg := repparser.Config{
		Commands: true,
		MapData:  true,
	}

	var mapDataHasher hash.Hash
	if *mapDataHash != "" {
		cfg.Debug = true
		switch strings.ToLower(*mapDataHash) {
		case "md5":
			mapDataHasher = md5.New()
		case "sha1":
			mapDataHasher = sha1.New()
		case "sha256":
			mapDataHasher = sha256.New()
		case "sha512":
			mapDataHasher = sha512.New()
		default:
			fmt.Printf("Invalid mapDataHash: %v\n", *mapDataHash)
			fmt.Println(validMapDataHashes)
			os.Exit(ExitCodeInvalidMapDataHash)
		}
	}

	if *dumpMapData {
		cfg.Debug = true
	}

	r, err := repparser.ParseFileConfig(args[0], cfg)
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(ExitCodeFailedToParseReplay)
	}

	var destination = os.Stdout

	if *outFile != "" {
		foutput, err := os.Create(*outFile)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			os.Exit(ExitCodeFailedToCreateOutputFile)
		}
		defer func() {
			if err := foutput.Close(); err != nil {
				panic(err)
			}
		}()

		destination = foutput
	}

	if *dumpMapData {
		if _, err := destination.Write(r.MapData.Debug.Data); err != nil {
			fmt.Printf("Failed to write map data: %v\n", err)
		}
		return
	}

	// custom holds any custom data we want in the output and is not part of rep.Replay
	custom := map[string]interface{}{}

	if *computed {
		r.Compute()
	}

	if mapDataHasher != nil {
		mapDataHasher.Write(r.MapData.Debug.Data)
		custom["MapDataHash"] = hex.EncodeToString(mapDataHasher.Sum(nil))
	}

	// Zero values in replay the user do not wish to see:
	if !*header {
		r.Header = nil
	}
	if !*mapData {
		r.MapData = nil
	} else {
		if !*mapTiles {
			r.MapData.Tiles = nil
		}
		if !*mapResLoc {
			r.MapData.MineralFields = nil
			r.MapData.Geysers = nil
		}
	}
	if !*cmds {
		r.Commands = nil
	}

	enc := json.NewEncoder(destination)

	if *indent {
		enc.SetIndent("", "  ")
	}

	var valueToEncode interface{} = r

	// If there are custom data, wrap (embed) the replay in a struct that holds the custom data too:
	if len(custom) > 0 {
		valueToEncode = struct {
			*rep.Replay
			Custom map[string]interface{}
		}{r, custom}
	}

	if err := enc.Encode(valueToEncode); err != nil {
		fmt.Printf("Failed to encode output: %v\n", err)
	}
}

func printVersion() {
	fmt.Println(appName, "version:", appVersion)
	fmt.Println("Parser version:", repparser.Version)
	fmt.Println("EAPM algorithm version:", rep.EAPMVersion)
	fmt.Println("Platform:", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Built with:", runtime.Version())
	fmt.Println("Author:", appAuthor)
	fmt.Println("Home page:", appHome)
}

func printUsage() {
	fmt.Println("Usage:")
	name := os.Args[0]
	fmt.Printf("\t%s [FLAGS] repfile.rep\n", name)
	fmt.Println("\tRun with '-h' to see a list of available flags.")
}
