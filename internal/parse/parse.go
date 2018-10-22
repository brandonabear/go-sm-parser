package parse

import (
	"fmt"
	"go-sm-parser/internal/song"
	"strconv"
	"strings"
)

// BeatChanges parses header tags with changes like Stops or BPMs.
//
// Raw => "0.000=179.000,920.000=117.073"
// Parsed => [{0 179} {920 117.073}]
func BeatChanges(changes string) []song.BeatChange {
	if changes != "" {
		changeArray := strings.Split(changes, ",")
		parsedChanges := make([]song.BeatChange, len(changeArray))
		for idx, change := range changeArray {
			beat, _ := strconv.ParseFloat(strings.Split(change, "=")[0], 64)
			value, _ := strconv.ParseFloat(strings.Split(change, "=")[1], 64)
			changeStruct := song.BeatChange{Beat: beat, Value: value}
			parsedChanges[idx] = changeStruct
		}
		return parsedChanges
	}
	return nil
}

// TagValue retrieves the value from a Header Tag
//
// Raw => "#TITLE:Song Title;"
// Parsed => "Song Title"
func TagValue(tag string) string {
	removedNewline := strings.Replace(tag, "\n", "", -1)
	trimmedWhitespace := strings.TrimSpace(removedNewline)
	smValue := strings.Split(trimmedWhitespace, ":")[1]
	return smValue
}

// NoteValue returns an array of Note elements
func NoteValue(tag string) []string {
	removedNewline := strings.Replace(tag, "\n", "", -1)
	trimmedWhitespace := strings.TrimSpace(removedNewline)
	noteValue := strings.Split(trimmedWhitespace, ":")

	for idx, note := range noteValue {
		noteValue[idx] = strings.TrimSpace(note)
	}
	return noteValue
}

// LineContains checks for a simfile Header tag
func LineContains(line string, match string) bool {
	return strings.Contains(line, match)
}

// RadarCategory parses the radar values
func RadarCategory(radar string) song.Radar {
	categories := strings.Split(radar, ",")

	stream, _ := strconv.ParseFloat(categories[0], 64)
	voltage, _ := strconv.ParseFloat(categories[1], 64)
	air, _ := strconv.ParseFloat(categories[2], 64)
	freeze, _ := strconv.ParseFloat(categories[3], 64)
	chaos, _ := strconv.ParseFloat(categories[4], 64)

	grooveRadar := song.Radar{
		Stream:  stream,
		Voltage: voltage,
		Air:     air,
		Freeze:  freeze,
		Chaos:   chaos}
	return grooveRadar
}

func calcQuantization(measure string) int {
	return len(measure) / 4
}

// func splitSteps(measure string) []song.Step {

// }

// NoteData captures beat/measure information
func NoteData(notes string) []song.Measure {
	measureSlices := []song.Measure{}
	measures := strings.SplitAfter(notes, ",")
	for idx, measureString := range measures {
		measureClean := strings.Replace(strings.Replace(measureString, "\r", "", -1), ",", "", -1)
		quantization := calcQuantization(measureClean)
		measure := song.Measure{Index: idx, Quantization: quantization}
		measureSlices = append(measureSlices, measure)
		fmt.Println(measure)
	}
	return measureSlices
}

// func main() {
// https://kodejava.org/how-to-split-a-string-by-a-number-of-characters/
// 	index := 0
// 	measure := "0010000010000000"
// 	quantization := len(measure) / 4
// 	beatPart := 0
// 	step := Step{}
// 	measures := []Step
// 	for idx, arrow := range measure {
// 		if math.Mod(idx, 4) != 0 {
// 			right, _ := strconv.ParseInt(arrow, 64)
// 			step.Right = right
// 			step.Beat = beatPart
// 			measures = append(measures, step)
// 			beatPart += 1 / quantization
// 			step := Step{}
// 		} else {

// 		}
// 	}

// 	fmt.Println(quantization)
// }
