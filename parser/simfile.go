package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

// Simfile represents a single Stepmania simfile.
type Simfile struct {
	SongPack string  `json:"song_pack"`
	Header   Header  `json:"header"`
	Charts   []Chart `json:"charts"`
}

// PackName extracts the pack name from the parent directory of the song folder.
func PackName(smDir string) string {
	packDir := path.Dir(path.Dir(smDir))
	_, packName := path.Split(packDir)
	return packName
}

// ReadSM returns a byte array from a .sm file
func ReadSM(smPath string) ([]uint8, error) {
	if path.Ext(smPath) == ".sm" {
		data, _ := ioutil.ReadFile(smPath)
		return data, nil
	}
	return nil, errors.New("Extension Error: File is not of type .sm")
}

// WriteJSON serializes parsed Simfile data as JSON.
func WriteJSON(sim Simfile, jsonPath string) error {
	simJSON, err := json.Marshal(sim)
	CheckError(err)
	outputName := fmt.Sprintf("%s/%s.json", jsonPath, sim.Header.Title)
	ioutil.WriteFile(outputName, simJSON, 0644)
	return err
}
