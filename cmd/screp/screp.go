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
	"runtime"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
)

const (
	appName    = "screp"
	appVersion = "v1.5.0"
	appAuthor  = "Andras Belicza"
	appHome    = "https://github.com/icza/screp"
)

// Flag variables
var (
	version = flag.Bool("version", false, "print version info and exit")

	header    = flag.Bool("header", true, "print replay header")
	mapData   = flag.Bool("map", false, "print map data")
	mapTiles  = flag.Bool("maptiles", false, "print map data tiles; valid with 'map'")
	mapResLoc = flag.Bool("mapres", false, "print map data resource locations (minerals and geysers); valid with 'map'")
	cmds      = flag.Bool("cmds", false, "print player commands")
	computed  = flag.Bool("computed", true, "print computed / derived data")
	outFile   = flag.String("outfile", "", "optional output file name")

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
		os.Exit(1)
	}

	r, err := repparser.ParseFile(args[0])
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(2)
	}

	if *computed {
		r.Compute()
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

	var enc *json.Encoder

	if *outFile == "" {
		enc = json.NewEncoder(os.Stdout)
	} else {
		fp, err := os.Create(*outFile)
		if err != nil {
			fmt.Printf("Failed to create output file: %v\n", err)
			os.Exit(3)
		}
		defer func() {
			if err := fp.Close(); err != nil {
				panic(err)
			}
		}()
		enc = json.NewEncoder(fp)
	}

	if *indent {
		enc.SetIndent("", "  ")
	}
	enc.Encode(r)
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
