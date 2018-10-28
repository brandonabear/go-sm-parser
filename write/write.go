package write

import (
	"encoding/json"
	"go-sm-parser/simfile"
	"io/ioutil"
)

// AsJSON serializes parsed Simfile data as JSON.
func AsJSON(sim simfile.Simfile) error {
	simJSON, err := json.Marshal(sim)
	outputName := sim.Header.Title + ".json"
	ioutil.WriteFile(outputName, simJSON, 0644)
	return err
}
