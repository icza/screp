# screp

[![Build Status](https://travis-ci.org/icza/screp.svg?branch=master)](https://travis-ci.org/icza/screp)
[![GoDoc](https://godoc.org/github.com/icza/screp?status.svg)](https://godoc.org/github.com/icza/screp)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/screp)](https://goreportcard.com/report/github.com/icza/screp)

StarCraft: Brood War replay parser.

Initially this parser was developed as part of the [repMastered](https://github.com/icza/repmastered)
project (the initial history can be found there), but was outsourced here as a separate project.

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

Install the command line application via:

	go get github.com/icza/screp/...

This will place `screp` inside `$GOPATH/bin/`.

## Example projects using this

- https://github.com/icza/repmastered
