package parser

import (
	"testing"

	"github.com/spf13/afero"
)

func TestPackName(t *testing.T) {
	smDir := "/home/user/stepmania/pack/song/song.sm"
	packName := PackName(smDir)
	if packName != "pack" {
		t.Error("Pack name not parsed properly.")
	}
}

func TestReadSM(t *testing.T) {
	var tests = []struct {
		smPath string
	}{
		{"../testdata/sharpnelstreamz/bluearmy/bluearmy.sm"},
		{"../testdata/README.md"},
	}

	for _, test := range tests {
		data, err := ReadSM(test.smPath)
		if (data != nil) && (err != nil) {
			t.Error("ReadSM read file but returned error.")
		}
		if (data == nil) && (err.Error() != "Extension Error: File is not of type .sm") {
			t.Error("ReadSM read incompatible file extension but returned improper error.")
		}
	}
}

func TestWriteJson(t *testing.T) {
	var Fs = afero.NewOsFs()
	sim := Simfile{}
	outputDir, _ := afero.TempDir(Fs, "/tmp", "_")
	err := WriteJSON(sim, outputDir)
	if err != nil {
		t.Error("JSON write failed.")
	}
}
