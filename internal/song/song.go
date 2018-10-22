package song

// BeatChange structure
type BeatChange struct {
	Beat, Value float64
}

// Header structure
type Header struct {
	Title            string
	Subtitle         string
	Artist           string
	TitleTranslit    string
	SubtitleTranslit string
	ArtistTranslit   string
	Genre            string
	Credit           string
	Banner           string
	Background       string
	LyricsPath       string
	CDTitle          string
	Music            string
	Offset           float64
	SampleStart      float64
	SampleLength     float64
	Selectable       string
	DisplayBPM       float64
	BPMs             []BeatChange
	Stops            []BeatChange
	BGChanges        []BeatChange
	KeySounds        []BeatChange
}

// Chart structure
type Chart struct {
	RawData     string
	Type        string
	Description string
	Difficulty  string
	Meter       int64
	GrooveRadar Radar
	Notes       []Measure
}

// Radar structure
type Radar struct {
	Stream  float64
	Voltage float64
	Air     float64
	Freeze  float64
	Chaos   float64
}

// Step structure
type Step struct {
	Beat  float64
	Left  int
	Down  int
	Up    int
	Right int
}

// Measure structure
type Measure struct {
	Index        int
	Quantization int
	Steps        []Step
}

// Simfile structure
type Simfile struct {
	Header
	Charts []Chart
}
