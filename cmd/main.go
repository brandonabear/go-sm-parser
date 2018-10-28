// Package main implements a Stepmania Simfile (.sm) parser.
package main

import (
	"encoding/json"
	"go-sm-parser/internal/parse"
	"go-sm-parser/internal/song"
	"io/ioutil"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	header := song.Header{}
	simfile := song.Simfile{}

	// Read simfile
	dat, err := ioutil.ReadFile("/home/brandon/go_projects/src/go-sm-parser/examples/twothousand.sm")
	check(err)
	tags := strings.Split(string(dat), ";")

	// Parse the simfile tags
	for idx := range tags {
		tag := tags[idx]
		switch {
		case parse.LineContains(tag, "#TITLE:"):
			header.Title = parse.TagValue(tag)
		case parse.LineContains(tag, "#SUBTITLE:"):
			header.Subtitle = parse.TagValue(tag)
		case parse.LineContains(tag, "#ARTIST:"):
			header.Artist = parse.TagValue(tag)
		case parse.LineContains(tag, "#TITLETRANSLIT:"):
			header.TitleTranslit = parse.TagValue(tag)
		case parse.LineContains(tag, "#SUBTITLETRANSLIT:"):
			header.SubtitleTranslit = parse.TagValue(tag)
		case parse.LineContains(tag, "#ARTISTTRANSLIT:"):
			header.ArtistTranslit = parse.TagValue(tag)
		case parse.LineContains(tag, "#GENRE:"):
			header.Genre = parse.TagValue(tag)
		case parse.LineContains(tag, "#CREDIT:"):
			header.Credit = parse.TagValue(tag)
		case parse.LineContains(tag, "#BANNER:"):
			header.Banner = parse.TagValue(tag)
		case parse.LineContains(tag, "#BACKGROUND:"):
			header.Background = parse.TagValue(tag)
		case parse.LineContains(tag, "#LYRICSPATH:"):
			header.LyricsPath = parse.TagValue(tag)
		case parse.LineContains(tag, "#CDTITLE:"):
			header.CDTitle = parse.TagValue(tag)
		case parse.LineContains(tag, "#MUSIC:"):
			header.Music = parse.TagValue(tag)
		case parse.LineContains(tag, "#OFFSET:"):
			header.Offset, err = strconv.ParseFloat(parse.TagValue(tag), 64)
		case parse.LineContains(tag, "#SAMPLESTART:"):
			header.SampleStart, err = strconv.ParseFloat(parse.TagValue(tag), 64)
		case parse.LineContains(tag, "#SAMPLELENGTH:"):
			header.SampleLength, err = strconv.ParseFloat(parse.TagValue(tag), 64)
		case parse.LineContains(tag, "#SELECTABLE:"):
			header.Selectable = parse.TagValue(tag)
		case parse.LineContains(tag, "#DISPLAYBPM:"):
			header.DisplayBPM, err = strconv.ParseFloat(parse.TagValue(tag), 64)
		case parse.LineContains(tag, "#BPMS:"):
			bpmString := parse.TagValue(tag)
			header.BPMs = parse.BeatChanges(bpmString)
		case parse.LineContains(tag, "#STOPS:"):
			stopString := parse.TagValue(tag)
			header.Stops = parse.BeatChanges(stopString)
		case parse.LineContains(tag, "#BGCHANGES:"):
			bgChangeString := parse.TagValue(tag)
			header.BGChanges = parse.BeatChanges(bgChangeString)
		case parse.LineContains(tag, "#KEYSOUNDS:"):
			keySoundString := parse.TagValue(tag)
			header.KeySounds = parse.BeatChanges(keySoundString)
		case parse.LineContains(tag, "#NOTES:"):
			chart := song.Chart{}
			chart.RawData = strings.TrimSpace(tag)
			simfile.Charts = append(simfile.Charts, chart)
		}

	}
	simfile.Header = header

	// Parse the chart note data
	for idx := range simfile.Charts {
		notes := parse.NoteValue(simfile.Charts[idx].RawData)
		simfile.Charts[idx].Type = notes[1]
		simfile.Charts[idx].Description = notes[2]
		simfile.Charts[idx].Difficulty = notes[3]
		meter, _ := strconv.Atoi(notes[4])
		simfile.Charts[idx].Meter = meter
		simfile.Charts[idx].GrooveRadar = parse.RadarCategory(notes[5])
		simfile.Charts[idx].Notes = parse.NoteData(notes[6])
		simfile.Charts[idx].RawData = ""
	}

	// Write as JSON
	simfileJSON, _ := json.Marshal(simfile)
	err = ioutil.WriteFile("test.json", simfileJSON, 0644)
}
