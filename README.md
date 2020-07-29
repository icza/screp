# screp

[![Build Status](https://travis-ci.org/icza/screp.svg?branch=master)](https://travis-ci.org/icza/screp)
[![GoDoc](https://godoc.org/github.com/icza/screp?status.svg)](https://godoc.org/github.com/icza/screp)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/screp)](https://goreportcard.com/report/github.com/icza/screp)

StarCraft: Brood War replay parser.

The package is designed to be used by other packages or apps, and is safe for concurrent use.
There is also an example CLI app that can be used standalone.

Parses both "modern" (starting from 1.18) and "legacy" (pre 1.18) replays.

_Check out the sister project to parse StarCraft II replays: [s2prot](https://github.com/icza/s2prot)_

## Using the `screp` CLI app

There is a command line application in the [cmd/screp](https://github.com/icza/screp/tree/master/cmd/screp) folder
which can be used to parse and display information about a single replay file.

The extracted data is displayed using JSON representation.

Usage is as simple as:

	screp [FLAGS] repfile.rep

Run with `-h` to see a list of available flags.

Example to parse a file called `sample.rep`, and display replay header (included by default)
and basic map data info (without tiles and resource location info):

	s2prot -map=true sample.rep

Or simply:

	s2prot -map sample.rep
	
## Installing the `screp` CLI app

The easiest is to download the binary release prepared for your platform from the [Releases](https://github.com/icza/screp/releases) page. Extract the archive and start using `screp`.

If you want to build `screp` from source, then simply clone the project and build the `cmd/screp` app:

	git clone https://github.com/icza/screp
	cd screp/cmd/screp
	go build

This will create an executable binary in the `cmd/screp` folder, ready to run.

## Example projects using this

- [repmastered.app](https://repmastered.app)
