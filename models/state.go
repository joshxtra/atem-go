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
