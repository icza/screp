/*
Package main is a simple CLI app to parse and display information about
a StarCraft: Brood War replay passed as a CLI argument.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/icza/screp/repparser"
)

// Flag variables
var (
	header    = flag.Bool("header", true, "print replay header")
	mapData   = flag.Bool("map", false, "print map data")
	mapTiles  = flag.Bool("mapTiles", false, "print map data tiles; valid with 'map'")
	mapResLoc = flag.Bool("mapResLoc", false, "print map data resource locations (minerals and geysers); valid with 'map'")
	cmds      = flag.Bool("cmds", false, "print player commands")

	indent = flag.Bool("indent", true, "use indentation when formatting output")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	r, err := repparser.ParseFile(args[0])
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(2)
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

	enc := json.NewEncoder(os.Stdout)
	if *indent {
		enc.SetIndent("", "  ")
	}
	enc.Encode(r)
}

func printUsage() {
	fmt.Println("Usage:")
	name := os.Args[0]
	fmt.Printf("\t%s [FLAGS] repfile.rep\n", name)
	fmt.Println("\tRun with '-h' to see available flags.")
}
