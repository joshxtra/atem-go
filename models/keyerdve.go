// Package models provides data models for ATEM switcher state and configuration.
package models

// KeyerDVE represents the DVE (Digital Video Effects) properties for a keyer.
type KeyerDVE struct {
	ME                  uint8
	Keyer               uint8
	Rate                uint8
	SizeX               int16
	SizeY               int16
	PositionX           int16
	PositionY           int16
	Rotation            int16
	BorderEnabled       bool
	BorderOuterWidth    int16
	BorderInnerWidth    int16
	BorderOuterSoftness uint8
	BorderInnerSoftness uint8
	BorderBevelSoftness uint8
	BorderBevelPosition uint8
	BorderBevel         bool
	BorderHue           uint16
	BorderSaturation    uint16
	BorderLuminance     uint16
	LightSourceDirection uint16
	LightSourceAltitude  uint16
	MaskEnabled          bool
	MaskTop              int16
	MaskBottom           int16
	MaskLeft             int16
	MaskRight            int16
	KeyFrame             bool
}

