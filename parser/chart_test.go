package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestTableCalcQuantization(t *testing.T) {
	var tests = []struct {
		measure      string
		quantization int
	}{
		{"0000", 1},
		{"00000000", 2},
		{"0000000000000000", 4},
		{"00000000000000000", 4},
		{"00000000000000000000000000000000", 8},
	}

	for _, test := range tests {
		if output := calcQuantization(test.measure); output != test.quantization {
			errorMsg := fmt.Sprintf("Expected %d, received: %d", test.quantization, output)
			t.Error(errorMsg)
		}
	}
}

func TestTableCalcBeat(t *testing.T) {
	var tests = []struct {
		measureNumber int
		beatPart      float64
		measureIndex  int
		beat          float64
	}{
		{0, 0.250, 0, 0.000},
		{1, 0.250, 0, 4.000},
		{1, 0.250, 4, 5.000},
		{1, 0.500, 4, 6.000},
	}

	for _, test := range tests {
		output := calcBeat(test.measureNumber, test.beatPart, test.measureIndex)
		if output != test.beat {
			errorMsg := fmt.Sprintf("Expected %f, received: %f", test.beat, output)
			t.Error(errorMsg)
		}

	}
}

func TestTableSplitSteps(t *testing.T) {
	var tests = []struct {
		measure       string
		measureNumber int
		quantization  int
		steps         []Step
	}{
		{"0000", 0, 4, []Step{Step{Beat: 0, L: "0", D: "0", U: "0", R: "0"}}},
		{"0011", 0, 4, []Step{Step{Beat: 0, L: "0", D: "0", U: "1", R: "1"}}},
		{"100M", 1, 4, []Step{Step{Beat: 4, L: "1", D: "0", U: "0", R: "M"}}},
	}

	for _, test := range tests {
		output := splitSteps(test.measure, test.measureNumber, test.quantization)
		for i, value := range output {
			if value.Beat != test.steps[i].Beat {
				t.Error("Parsed step incorrectly.")
			}
			if value.L != test.steps[i].L {
				t.Error("Parsed step incorrectly.")
			}
			if value.D != test.steps[i].D {
				t.Error("Parsed step incorrectly.")
			}
			if value.U != test.steps[i].U {
				t.Error("Parsed step incorrectly.")
			}
			if value.R != test.steps[i].R {
				t.Error("Parsed step incorrectly.")
			}
		}
	}
}

func TestRadarCategory(t *testing.T) {
	radar := "1.000,1.000,0.116,0.571,1.000"
	gr := radarCategory(radar)

	if gr.Stream != 1.000 {
		t.Error("Stream parsed incorrectly.")
	}
	if gr.Voltage != 1.000 {
		t.Error("Voltage parsed incorrectly.")
	}
	if gr.Air != 0.116 {
		t.Error("Air parsed incorrectly.")
	}
	if gr.Freeze != 0.571 {
		t.Error("Freeze parsed incorrectly.")
	}
	if gr.Chaos != 1.000 {
		t.Error("Chaos parsed incorrectly.")
	}
}

func TestTableRawNoteValue(t *testing.T) {
	var tests = []struct {
		tag    string
		values []string
	}{
		{"dance-single:16:0000,0000", []string{"dance-single", "16", "0000,0000"}},
		{"Challenge:1.000,1.000,0.116,0.571,1.000:00000000", []string{"Challenge", "1.000,1.000,0.116,0.571,1.000", "00000000"}},
	}

	for _, test := range tests {
		noteValue := RawNoteValue(test.tag)
		for j, value := range test.values {
			if noteValue[j] != value {
				errorMsg := fmt.Sprintf("Expected %s, received: %s", value, noteValue[j])
				t.Error(errorMsg)
			}
		}
	}
}

func TestTableNoteData(t *testing.T) {
	var tests = []struct {
		notes    string
		measures []Measure
	}{
		{"0000000000000000", []Measure{Measure{MeasureNumber: 0, Quantization: 4}}},
	}

	for _, test := range tests {
		output := noteData(test.notes)
		for i, measure := range test.measures {
			if measure.MeasureNumber != output[i].MeasureNumber {
				t.Error("Failed to parse note data.")
			}
			if measure.Quantization != output[i].Quantization {
				t.Error("Failed to parse note data.")
			}
		}
	}
}

func TestExtractCharts(t *testing.T) {
	sim := Simfile{}

	data := `
	#TITLE:Blue Army;
	#SUBTITLE:;
	#ARTIST:DJ Sharpnel;
	#TITLETRANSLIT:;
	#SUBTITLETRANSLIT:;
	#ARTISTTRANSLIT:;
	#GENRE:;
	#CREDIT:;
	#BANNER:bluearmybn.png;
	#BACKGROUND:bluearmybg.png;
	#LYRICSPATH:;
	#CDTITLE:;
	#MUSIC:bluearmy.ogg;
	#OFFSET:-0.701;
	#SAMPLESTART:200.271;
	#SAMPLELENGTH:21.073;
	#SELECTABLE:YES;
	#DISPLAYBPM:182.000;
	#BPMS:0.000=182.200,136.000=91.100,168.000=182.200,304.000=91.100,376.000=182.200,1056.000=91.100,1060.000=182.200,1156.000=91.100,1188.000=182.200;
	#STOPS:;
	#BGCHANGES:;
	#KEYSOUNDS:;

	//---------------dance-single - I. Pyles (v2 SH16)----------------
	#NOTES:
     dance-single:
     16 (Archi):
     Challenge:
     16:
     1.000,1.000,0.116,0.571,1.000:
	0000
	0000
	0000
	0000
	,
	0000
	0000
	0000
	0000
	;`

	// Parse the header tags.
	tags := strings.Split(string(data), ";")
	for i := range tags {
		tag := tags[i]
		sim = ExtractHeader(tag, sim)
	}

	// Parse the notes tag (chart data).
	for i := range sim.Charts {
		notes := RawNoteValue(sim.Charts[i].RawData)
		sim = ExtractCharts(i, notes, sim)
	}

	for _, chart := range sim.Charts {
		if chart.Type != "dance-single" {
			t.Error("Fuck.")
		}
		if chart.Meter != 16 {
			t.Error("Fuck.")
		}
	}
}
