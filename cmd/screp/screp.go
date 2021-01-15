/*

Package main is a simple CLI app to parse and display information about
a StarCraft: Brood War replay passed as a CLI argument.

*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/icza/screp/rep"
	"github.com/icza/screp/repparser"
)

const (
	appName    = "screp"
	appVersion = "v1.3.1"
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	userid := r.FormValue("userid") // x will be "" if parameter is not set
	fmt.Println("userid:", userid)
	if r.Method == "POST" {
		// 1. parse input
		r.ParseMultipartForm(10 << 20)
		// 2. retrieve file
		file, handler, err := r.FormFile("repFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		path := "replays/" + userid
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0700)
		}

		// 3. write temporary file on our server
		tempFile, err := ioutil.TempFile(path, "upload-*.rep")
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		tempFile.Write(fileBytes)
		fmt.Printf(tempFile.Name())
		parseRep(tempFile.Name())
		// 4. return result
		fmt.Fprintf(w, "Successfully Uploaded File\n")
	}
}
func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":443", nil)
}

func parseRep(repFile string) {

	fmt.Println(repFile[:len(repFile)-4])
	r, err := repparser.ParseFile(repFile)
	if err != nil {
		fmt.Printf("Failed to parse replay: %v\n", err)
		os.Exit(2)
	}
	r.Compute()
	var enc *json.Encoder

	fp, err := os.Create(repFile[:len(repFile)-4] + ".json")
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

	enc.SetIndent("", "  ")
	enc.Encode(r)
}

func main() {
	setupRoutes()
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
