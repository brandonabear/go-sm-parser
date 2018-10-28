// Package main implements a Stepmania Simfile parser.
// It currently supports the following format: sm
package main

import (
	"go-sm-parser/parse"
	"go-sm-parser/read"
	"go-sm-parser/simfile"
	"go-sm-parser/write"
	"os"
	"strings"
)

func main() {
	// Read the simfile data into memory.
	sim := simfile.Simfile{}
	smPath := os.Args[1]
	data, err := read.SM(smPath)
	parse.CheckError(err)

	// Parse the header tags.
	tags := strings.Split(string(data), ";")
	sim.SongPack = parse.PackName(smPath)
	for i := range tags {
		tag := tags[i]
		sim = parse.ExtractHeader(tag, sim)
	}

	// Parse the notes tag (chart data).
	for i := range sim.Charts {
		notes := parse.RawNoteValue(sim.Charts[i].RawData)
		sim = parse.ExtractCharts(i, notes, sim)
	}

	// Write to disk as JSON.
	err = write.AsJSON(sim)
	parse.CheckError(err)
}
