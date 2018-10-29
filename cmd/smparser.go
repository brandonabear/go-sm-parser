// Package main implements a Stepmania Simfile parser.
// It currently supports the following format: sm
package main

import (
	"go-sm-parser/parser"
	"os"
	"strings"
)

func main() {
	// Read the simfile data into memory.
	sim := parser.Simfile{}
	smPath := os.Args[1]
	data, err := parser.ReadSM(smPath)
	parser.CheckError(err)

	// Parse the header tags.
	tags := strings.Split(string(data), ";")
	sim.SongPack = parser.PackName(smPath)
	for i := range tags {
		tag := tags[i]
		sim = parser.ExtractHeader(tag, sim)
	}

	// Parse the notes tag (chart data).
	for i := range sim.Charts {
		notes := parser.RawNoteValue(sim.Charts[i].RawData)
		sim = parser.ExtractCharts(i, notes, sim)
	}

	// Write to disk as JSON.
	outDir := "../testdata"
	err = parser.WriteJSON(sim, outDir)
	parser.CheckError(err)
}
