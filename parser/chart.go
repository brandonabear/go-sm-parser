package parser

import (
	"math"
	"strconv"
	"strings"
)

// Chart contains individual chart attributes and note data.
type Chart struct {
	RawData     string    `json:"raw_data"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Meter       int       `json:"meter"`
	GrooveRadar Radar     `json:"grooveradar"`
	Notes       []Measure `json:"notes"`
}

// Radar represents the 5 GrooveRadar attributes.
type Radar struct {
	Stream  float64 `json:"stream"`
	Voltage float64 `json:"voltage"`
	Air     float64 `json:"air"`
	Freeze  float64 `json:"freeze"`
	Chaos   float64 `json:"chaos"`
}

// Measure contains metadata for a single chart measure.
type Measure struct {
	MeasureNumber int    `json:"measure_nbr"`
	Quantization  int    `json:"quantization"`
	Steps         []Step `json:"steps"`
}

// Step indicates the step pattern at a given beat.
type Step struct {
	Beat float64 `json:"beat"`
	L    string  `json:"l"`
	D    string  `json:"d"`
	U    string  `json:"u"`
	R    string  `json:"r"`
}

func calcQuantization(measure string) int {
	return len(measure) / 4
}

func calcBeat(measureNumber int, beatPart float64, measureIndex int) float64 {
	return 4 * (float64(measureNumber) + (beatPart * float64(measureIndex/4)))
}

func splitSteps(measure string, measureNumber int, quantization int) []Step {
	steps := []Step{}
	step := Step{}
	beatPart := 1.00 / float64(quantization)
	for i, m := range measure {
		switch remainder := math.Mod(float64(i+1), 4); remainder {
		case 1:
			step.L = string(m)
		case 2:
			step.D = string(m)
		case 3:
			step.U = string(m)
		case 0:
			step.R = string(m)
			step.Beat = calcBeat(measureNumber, beatPart, i)
			steps = append(steps, step)
			step = Step{}
		}
	}
	return steps
}

// RawNoteValue returns an array of raw Note elements
func RawNoteValue(tag string) []string {
	removedNewline := strings.Replace(tag, "\n", "", -1)
	trimmedWhitespace := strings.TrimSpace(removedNewline)
	noteValue := strings.Split(trimmedWhitespace, ":")

	for i, note := range noteValue {
		noteValue[i] = strings.TrimSpace(note)
	}
	return noteValue
}

// ExtractCharts parses the charts in the Notes tag
func ExtractCharts(i int, notes []string, sim Simfile) Simfile {
	// Only supports parsing singles
	if notes[1] == "dance-single" {
		sim.Charts[i].Type = notes[1]
		sim.Charts[i].Description = notes[2]
		sim.Charts[i].Difficulty = notes[3]
		meter, err := strconv.Atoi(notes[4])
		CheckError(err)
		sim.Charts[i].Meter = meter
		sim.Charts[i].GrooveRadar = radarCategory(notes[5])
		sim.Charts[i].Notes = noteData(notes[6])
	}
	sim.Charts[i].RawData = ""
	return sim
}

// radarCategory parses the radar values
func radarCategory(radar string) Radar {
	categories := strings.Split(radar, ",")

	stream, err := strconv.ParseFloat(categories[0], 64)
	CheckError(err)
	voltage, err := strconv.ParseFloat(categories[1], 64)
	CheckError(err)
	air, err := strconv.ParseFloat(categories[2], 64)
	CheckError(err)
	freeze, err := strconv.ParseFloat(categories[3], 64)
	CheckError(err)
	chaos, err := strconv.ParseFloat(categories[4], 64)
	CheckError(err)

	grooveRadar := Radar{
		Stream:  stream,
		Voltage: voltage,
		Air:     air,
		Freeze:  freeze,
		Chaos:   chaos}
	return grooveRadar
}

// noteData captures beat/measure information
func noteData(notes string) []Measure {
	measureSlices := []Measure{}
	measures := strings.SplitAfter(notes, ",")
	for measureNumber, measureString := range measures {
		measureClean := strings.Replace(strings.Replace(measureString, "\r", "", -1), ",", "", -1)
		quantization := calcQuantization(measureClean)
		steps := splitSteps(measureClean, measureNumber, quantization)
		measure := Measure{MeasureNumber: measureNumber, Quantization: quantization, Steps: steps}
		measureSlices = append(measureSlices, measure)
	}
	return measureSlices
}
