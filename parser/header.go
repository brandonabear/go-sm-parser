package parser

import (
	"strconv"
	"strings"
)

// Header contains shared information across charts in a simfile.
type Header struct {
	Title            string       `json:"title"`
	Subtitle         string       `json:"subtitle"`
	Artist           string       `json:"artist"`
	TitleTranslit    string       `json:"title_translit"`
	SubtitleTranslit string       `json:"subtitle_translit"`
	ArtistTranslit   string       `json:"artist_translit"`
	Genre            string       `json:"genre"`
	Credit           string       `json:"credit"`
	Banner           string       `json:"banner"`
	Background       string       `json:"background"`
	LyricsPath       string       `json:"lyrics_path"`
	CDTitle          string       `json:"cd_title"`
	Music            string       `json:"music"`
	Offset           float64      `json:"offset"`
	SampleStart      float64      `json:"sample_start"`
	SampleLength     float64      `json:"sample_length"`
	Selectable       string       `json:"selectable"`
	DisplayBPM       []float64    `json:"display_bpm"`
	BPMs             []BeatChange `json:"bpms"`
	Stops            []BeatChange `json:"stops"`
	BGChanges        []BeatChange `json:"bg_changes"`
	KeySounds        []BeatChange `json:"keysounds"`
}

// BeatChange is a Beat/Value pair representing a change in a song (ex. Stops).
type BeatChange struct {
	Beat  float64 `json:"beat"`
	Value float64 `json:"value"`
}

// ExtractHeader parses the Header tags.
func ExtractHeader(tag string, sim Simfile) Simfile {
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
		displayBPMText := tagValue(tag)
		sim.Header.DisplayBPM = displayBPM(displayBPMText)
	case lineContains(tag, "#BPMS:"):
		bpmString := tagValue(tag)
		sim.Header.BPMs = extractBeatChanges(bpmString)
	case lineContains(tag, "#STOPS:"):
		stopString := tagValue(tag)
		sim.Header.Stops = extractBeatChanges(stopString)
	case lineContains(tag, "#BGCHANGES:"):
		bgChangeString := tagValue(tag)
		sim.Header.BGChanges = extractBeatChanges(bgChangeString)
	case lineContains(tag, "#KEYSOUNDS:"):
		keySoundString := tagValue(tag)
		sim.Header.KeySounds = extractBeatChanges(keySoundString)
	case lineContains(tag, "#NOTES:"):
		chart := Chart{}
		chart.RawData = strings.TrimSpace(tag)
		sim.Charts = append(sim.Charts, chart)
	}
	return sim
}

// extractBeatChanges parses header tags with changes like Stops or BPMs.
//
// Raw => "0.000=179.000,920.000=117.073"
// Parsed => [{0 179} {920 117.073}]
func extractBeatChanges(changes string) []BeatChange {
	if changes != "" {
		changeArray := strings.Split(changes, ",")
		parsedChanges := make([]BeatChange, len(changeArray))
		for i, change := range changeArray {
			beat, _ := strconv.ParseFloat(strings.Split(change, "=")[0], 64)
			value, _ := strconv.ParseFloat(strings.Split(change, "=")[1], 64)
			changeStruct := BeatChange{Beat: beat, Value: value}
			parsedChanges[i] = changeStruct
		}
		return parsedChanges
	}
	return nil
}

// lineContains checks for a simfile Header tag
func lineContains(line string, match string) bool {
	return strings.Contains(line, match)
}

// tagValue retrieves the value from a Header Tag
//
// Raw => "#TITLE:Song Title;"
// Parsed => "Song Title"
func tagValue(tag string) string {
	removedNewline := strings.Replace(tag, "\n", "", -1)
	removedSemicolon := strings.Replace(removedNewline, ";", "", -1)
	trimmedWhitespace := strings.TrimSpace(removedSemicolon)
	smValue := strings.Split(trimmedWhitespace, ":")[1]
	return smValue
}

func displayBPM(displayBPMText string) []float64 {
	switch {
	case lineContains(displayBPMText, "*"):
		return []float64{0.00}
	case lineContains(displayBPMText, ":"):
		bpms := strings.Split(displayBPMText, ":")
		floor, _ := strconv.ParseFloat(bpms[0], 64)
		ceil, _ := strconv.ParseFloat(bpms[1], 64)
		return []float64{floor, ceil}
	default:
		bpm, _ := strconv.ParseFloat(displayBPMText, 64)
		return []float64{bpm}
	}
}
