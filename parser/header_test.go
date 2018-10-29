package parser

import (
	"fmt"
	"testing"
)

func TestTableLineContains(t *testing.T) {
	var tests = []struct {
		line   string
		match  string
		result bool
	}{
		{"#TITLE:Song Title;", "#TITLE:", true},
		{"#SUBTITLE:Subtitle;", "#TITLE:", false},
		{"#TITLETRANSLIT:;", "#TITLE:", false},
	}

	for _, test := range tests {
		if output := lineContains(test.line, test.match); output != test.result {
			errorMsg := fmt.Sprintf("Expected %s to contain %s", test.line, test.match)
			t.Error(errorMsg)
		}
	}
}

func TestTableTagValue(t *testing.T) {
	var tests = []struct {
		tag   string
		value string
	}{
		{"#TITLE:Song Title;", "Song Title"},
		{"#SAMPLESTART:200.271;", "200.271"},
		{"#BPMS:0.000=182.200,136.000=91.100,168.000=182.200,304.000=91.100,376.000=182.200,1056.000=91.100,1060.000=182.200,1156.000=91.100,1188.000=182.200;", "0.000=182.200,136.000=91.100,168.000=182.200,304.000=91.100,376.000=182.200,1056.000=91.100,1060.000=182.200,1156.000=91.100,1188.000=182.200"},
		{"#STOPS:;", ""},
	}

	for _, test := range tests {
		if output := tagValue(test.tag); output != test.value {
			errorMsg := fmt.Sprintf("Expected value %s from tag %s, received: %s", test.value, test.tag, output)
			t.Error(errorMsg)
		}
	}
}

func TestTableDisplayBPM(t *testing.T) {
	var tests = []struct {
		bpmText string
		result  []float64
	}{
		{"*", []float64{0.00}},
		{"100:200", []float64{100.000, 200.000}},
		{"200.000", []float64{200.000}},
	}

	for _, test := range tests {
		if output := displayBPM(test.bpmText); output[0] != test.result[0] {
			errorMsg := fmt.Sprintf("Expected %f, received: %f", test.result[0], output[0])
			t.Error(errorMsg)
		}
	}
}

func TestTableExtractBeatChanges(t *testing.T) {
	var tests = []struct {
		changes string
		result  []BeatChange
	}{
		{"", nil},
		{"0.000=182.200", []BeatChange{BeatChange{Beat: 0.000, Value: 182.200}}},
		{"0.000=200.000,196.500=201.000", []BeatChange{BeatChange{Beat: 0.000, Value: 200.000}, BeatChange{Beat: 196.500, Value: 201.000}}},
	}

	for _, test := range tests {
		output := extractBeatChanges(test.changes)
		if test.result == nil {
			if output != nil {
				t.Error("Expected nil result.")
			}
		}
		for i := range output {
			if output[i].Beat != test.result[i].Beat {
				errorMsg := fmt.Sprintf("Expected %f, received: %f", test.result[i].Beat, output[i].Beat)
				t.Error(errorMsg)
			}
		}
	}
}

func TestTableExtractHeader(t *testing.T) {
	var sim = Simfile{}
	var tags = []string{
		"#TITLE:SongTitle;",
		"#SUBTITLE:;",
		"#ARTIST:Artist;",
		"#TITLETRANSLIT:;",
		"#SUBTITLETRANSLIT:;",
		"ARTISTTRANSLIT:;",
		"#GENRE:Metal;",
		"#CREDIT:barndoor;",
		"#BANNER:bn-song.png;",
		"#BACKGROUND:bg-song.png;",
		"#LYRICSPATH:;",
		"#CDTITLE:;",
		"#MUSIC:;",
		"#OFFSET:-0.125;",
		"#SAMPLESTART:0.000;",
		"#SAMPLELENGTH:20.000;",
		"#SELECTABLE:;",
		"#DISPLAYBPM:200.000;",
		"#BPMS:0.000=200.000;",
		"#STOPS:;",
		"#BGCHANGES:;",
		"#KEYSOUNDS:;",
		"#NOTES:;",
	}

	for _, tag := range tags {
		sim = ExtractHeader(tag, sim)
	}

	if sim.Header.Title != "SongTitle" {
		t.Error("Title not parsed correctly.")
	}
	if sim.Header.Artist != "Artist" {
		t.Error("Artist not parsed correctly.")
	}
	if sim.Header.Credit != "barndoor" {
		t.Error("Credit not parsed correctly.")
	}
}
