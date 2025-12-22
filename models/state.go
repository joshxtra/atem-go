// Package models provides data models for ATEM switcher state and configuration.
package models

// PowerStatus represents the power state of the switcher.
type PowerStatus struct {
	Main   bool
	Backup bool
}

// TallyState represents the tally state (program/preview) for a source.
type TallyState struct {
	Program bool
	Preview bool
}

// Bitmask is the on-wire representation of this state
func (t TallyState) Bitmask() uint8 {
	var mask uint8
	if t.Program {
		mask |= 1
	}
	if t.Preview {
		mask |= 2
	}
	return mask
}

// Timecode represents a timecode value with hour, minute, second, and frame.
type Timecode struct {
	Hour   int
	Minute int
	Second int
	Frame  int
}

// KeyerState represents the state of a keyer in a mix effect bus.
type KeyerState struct {
	OnAir        bool // On air in program
	PreviewOnAir bool // On air in preview (if supported)
	FillSource   VideoSource
	KeySource    VideoSource
	DVE          *DVEProperties // DVE properties (nil if not a DVE keyer)
}

// DVEProperties represents Digital Video Effects properties for a keyer
type DVEProperties struct {
	Rate          uint8
	SizeX         int16 // Fixed point format (typically * 1000)
	SizeY         int16 // Fixed point format (typically * 1000)
	PositionX     int16 // Fixed point format (typically * 1000)
	PositionY     int16 // Fixed point format (typically * 1000)
	Rotation      int16 // Fixed point format (typically * 100)
	BorderEnabled bool
	BorderStyle   uint8
	BorderWidth   int16 // Fixed point format (typically * 1000)
	BorderOpacity uint8
	ShadowEnabled bool
	LightSource   uint8
	MaskEnabled   bool
	MaskTop       int16
	MaskBottom    int16
	MaskLeft      int16
	MaskRight     int16
}
