package simfile

// Simfile represents a single Stepmania simfile.
type Simfile struct {
	SongPack string  `json:"song_pack"`
	Header   Header  `json:"header"`
	Charts   []Chart `json:"charts"`
}

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
	DisplayBPM       float64      `json:"display_bpm"`
	BPMs             []BeatChange `json:"bpms"`
	Stops            []BeatChange `json:"stops"`
	BGChanges        []BeatChange `json:"bg_changes"`
	KeySounds        []BeatChange `json:"keysounds"`
}

// BeatChange is a Beat/Value pair representing a change in a song (ex. Stops).
type BeatChange struct {
	Beat  float64 `json:"bear"`
	Value float64 `json:"value"`
}

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
