package parse

import (
	"go-sm-parser/internal/song"
	"math"
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

func calcBeat(measureNumber int, beatPart float64, measureIndex int) float64 {
	return 4 * (float64(measureNumber) + (beatPart * float64(measureIndex/4)))
}

func splitSteps(measure string, measureNumber int, quantization int) []song.Step {
	steps := []song.Step{}
	step := song.Step{}
	beatPart := 1.00 / float64(quantization)
	for idx, m := range measure {
		switch remainder := math.Mod(float64(idx+1), 4); remainder {
		case 1:
			step.Left = string(m)
		case 2:
			step.Down = string(m)
		case 3:
			step.Up = string(m)
		case 0:
			step.Right = string(m)
			step.Beat = calcBeat(measureNumber, beatPart, idx)
			steps = append(steps, step)
			step = song.Step{}
		}
	}
	return steps
}

// NoteData captures beat/measure information
func NoteData(notes string) []song.Measure {
	measureSlices := []song.Measure{}
	measures := strings.SplitAfter(notes, ",")
	for measureNumber, measureString := range measures {
		measureClean := strings.Replace(strings.Replace(measureString, "\r", "", -1), ",", "", -1)
		quantization := calcQuantization(measureClean)
		steps := splitSteps(measureClean, measureNumber, quantization)
		measure := song.Measure{Index: measureNumber, Quantization: quantization, Steps: steps}
		measureSlices = append(measureSlices, measure)
	}
	return measureSlices
}
