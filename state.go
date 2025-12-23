package atem

import "github.com/mraerino/atem-go/models"

// Characteristics describes the capabilities of the ATEM switcher.
type Characteristics struct {
	ProductName   string
	Topology      models.Topology
	MixEffect     models.MixEffectConfig
	MediaPlayer   models.MediaPlayerConfig
	MultiViews    int
	SuperSources  int // num of boxes
	TallyChannels int
	MacroBanks    int
}

// SwitcherState represents the complete state of the ATEM switcher.
type SwitcherState struct {
	Version struct {
		Major int
		Minor int
	}
	Warning   string
	Power     models.PowerStatus
	VideoMode models.VideoMode

	Inputs map[models.VideoSource]models.InputProperties
	Config Characteristics

	Program map[int]models.VideoSource
	Preview map[int]models.VideoSource
	Aux     map[int]models.VideoSource

	TallyByIndex  map[int]models.TallyState
	TallyBySource map[models.VideoSource]models.TallyState

	MediaPlayer map[int]*models.MediaPlayer
	MediaFiles  map[int]models.MediaStillFrame

	TimeCodeLastChange models.Timecode
}

// NewSwitcherState creates a new empty switcher state.
func NewSwitcherState() SwitcherState {
	return SwitcherState{
		Inputs: make(map[models.VideoSource]models.InputProperties),
		Config: Characteristics{
			MixEffect: make(models.MixEffectConfig),
		},
		Program: make(map[int]models.VideoSource),
		Preview: make(map[int]models.VideoSource),
		Aux:     make(map[int]models.VideoSource),

		MediaPlayer: make(map[int]*models.MediaPlayer),
		MediaFiles:  make(map[int]models.MediaStillFrame),
	}
}
