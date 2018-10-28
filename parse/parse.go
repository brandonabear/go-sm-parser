package parse

import (
	"go-sm-parser/simfile"
	"math"
	"path"
	"strconv"
	"strings"
)

// CheckError panics at the presence of an error.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// PackName extracts the pack name from the parent directory of the song folder.
func PackName(simDirectory string) string {
	packDirectory := path.Dir(path.Dir(simDirectory))
	_, packName := path.Split(packDirectory)
	return packName
}

// beatChanges parses header tags with changes like Stops or BPMs.
//
// Raw => "0.000=179.000,920.000=117.073"
// Parsed => [{0 179} {920 117.073}]
func beatChanges(changes string) []simfile.BeatChange {
	if changes != "" {
		changeArray := strings.Split(changes, ",")
		parsedChanges := make([]simfile.BeatChange, len(changeArray))
		for idx, change := range changeArray {
			beat, _ := strconv.ParseFloat(strings.Split(change, "=")[0], 64)
			value, _ := strconv.ParseFloat(strings.Split(change, "=")[1], 64)
			changeStruct := simfile.BeatChange{Beat: beat, Value: value}
			parsedChanges[idx] = changeStruct
		}
		return parsedChanges
	}
	return nil
}

// tagValue retrieves the value from a Header Tag
//
// Raw => "#TITLE:Song Title;"
// Parsed => "Song Title"
func tagValue(tag string) string {
	removedNewline := strings.Replace(tag, "\n", "", -1)
	trimmedWhitespace := strings.TrimSpace(removedNewline)
	smValue := strings.Split(trimmedWhitespace, ":")[1]
	return smValue
}

// lineContains checks for a simfile Header tag
func lineContains(line string, match string) bool {
	return strings.Contains(line, match)
}

func calcQuantization(measure string) int {
	return len(measure) / 4
}

func calcBeat(measureNumber int, beatPart float64, measureIndex int) float64 {
	return 4 * (float64(measureNumber) + (beatPart * float64(measureIndex/4)))
}

func splitSteps(measure string, measureNumber int, quantization int) []simfile.Step {
	steps := []simfile.Step{}
	step := simfile.Step{}
	beatPart := 1.00 / float64(quantization)
	for idx, m := range measure {
		switch remainder := math.Mod(float64(idx+1), 4); remainder {
		case 1:
			step.L = string(m)
		case 2:
			step.D = string(m)
		case 3:
			step.U = string(m)
		case 0:
			step.R = string(m)
			step.Beat = calcBeat(measureNumber, beatPart, idx)
			steps = append(steps, step)
			step = simfile.Step{}
		}
	}
	return steps
}

// ExtractHeader parses the Header tags.
func ExtractHeader(tag string, sim simfile.Simfile) simfile.Simfile {
	switch {
	case lineContains(tag, "#TITLE:"):
		sim.Header.Title = tagValue(tag)
	case lineContains(tag, "#SUBTITLE:"):
		sim.Header.Subtitle = tagValue(tag)
	case lineContains(tag, "#ARTIST:"):
		sim.Header.Artist = tagValue(tag)
	case lineContains(tag, "#TITLETRANSLIT:"):
		sim.Header.TitleTranslit = tagValue(tag)
	case lineContains(tag, "#SUBTITLETRANSLIT:"):
		sim.Header.SubtitleTranslit = tagValue(tag)
	case lineContains(tag, "#ARTISTTRANSLIT:"):
		sim.Header.ArtistTranslit = tagValue(tag)
	case lineContains(tag, "#GENRE:"):
		sim.Header.Genre = tagValue(tag)
	case lineContains(tag, "#CREDIT:"):
		sim.Header.Credit = tagValue(tag)
	case lineContains(tag, "#BANNER:"):
		sim.Header.Banner = tagValue(tag)
	case lineContains(tag, "#BACKGROUND:"):
		sim.Header.Background = tagValue(tag)
	case lineContains(tag, "#LYRICSPATH:"):
		sim.Header.LyricsPath = tagValue(tag)
	case lineContains(tag, "#CDTITLE:"):
		sim.Header.CDTitle = tagValue(tag)
	case lineContains(tag, "#MUSIC:"):
		sim.Header.Music = tagValue(tag)
	case lineContains(tag, "#OFFSET:"):
		sim.Header.Offset, _ = strconv.ParseFloat(tagValue(tag), 64)
	case lineContains(tag, "#SAMPLESTART:"):
		sim.Header.SampleStart, _ = strconv.ParseFloat(tagValue(tag), 64)
	case lineContains(tag, "#SAMPLELENGTH:"):
		sim.Header.SampleLength, _ = strconv.ParseFloat(tagValue(tag), 64)
	case lineContains(tag, "#SELECTABLE:"):
		sim.Header.Selectable = tagValue(tag)
	case lineContains(tag, "#DISPLAYBPM:"):
		sim.Header.DisplayBPM, _ = strconv.ParseFloat(tagValue(tag), 64)
	case lineContains(tag, "#BPMS:"):
		bpmString := tagValue(tag)
		sim.Header.BPMs = beatChanges(bpmString)
	case lineContains(tag, "#STOPS:"):
		stopString := tagValue(tag)
		sim.Header.Stops = beatChanges(stopString)
	case lineContains(tag, "#BGCHANGES:"):
		bgChangeString := tagValue(tag)
		sim.Header.BGChanges = beatChanges(bgChangeString)
	case lineContains(tag, "#KEYSOUNDS:"):
		keySoundString := tagValue(tag)
		sim.Header.KeySounds = beatChanges(keySoundString)
	case lineContains(tag, "#NOTES:"):
		chart := simfile.Chart{}
		chart.RawData = strings.TrimSpace(tag)
		sim.Charts = append(sim.Charts, chart)
	}
	return sim
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
func ExtractCharts(i int, notes []string, sim simfile.Simfile) simfile.Simfile {
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
func radarCategory(radar string) simfile.Radar {
	categories := strings.Split(radar, ",")

	stream, _ := strconv.ParseFloat(categories[0], 64)
	voltage, _ := strconv.ParseFloat(categories[1], 64)
	air, _ := strconv.ParseFloat(categories[2], 64)
	freeze, _ := strconv.ParseFloat(categories[3], 64)
	chaos, _ := strconv.ParseFloat(categories[4], 64)

	grooveRadar := simfile.Radar{
		Stream:  stream,
		Voltage: voltage,
		Air:     air,
		Freeze:  freeze,
		Chaos:   chaos}
	return grooveRadar
}

// noteData captures beat/measure information
func noteData(notes string) []simfile.Measure {
	measureSlices := []simfile.Measure{}
	measures := strings.SplitAfter(notes, ",")
	for measureNumber, measureString := range measures {
		measureClean := strings.Replace(strings.Replace(measureString, "\r", "", -1), ",", "", -1)
		quantization := calcQuantization(measureClean)
		steps := splitSteps(measureClean, measureNumber, quantization)
		measure := simfile.Measure{MeasureNumber: measureNumber, Quantization: quantization, Steps: steps}
		measureSlices = append(measureSlices, measure)
	}
	return measureSlices
}
