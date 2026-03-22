package cmds

import (
	"encoding/binary"
	"errors"

	"github.com/joshxtra/atem-go/models"
)

type videoSourceState struct {
	Bus    uint8
	Source models.VideoSource
}

func (m *videoSourceState) UnmarshalBinary(data []byte) error {
	m.Bus = data[0]
	m.Source = models.VideoSource(binary.BigEndian.Uint16(data[2:]))
	return nil
}

func (m *videoSourceState) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = m.Bus
	binary.BigEndian.PutUint16(pl[2:], uint16(m.Source))
	return pl, nil
}

// PrgiCmd represents a program input command.
type PrgiCmd struct {
	videoSourceState
}

// Slug returns the command slug.
func (PrgiCmd) Slug() string {
	return "PrgI"
}

// PrviCmd represents a preview input command.
type PrviCmd struct {
	videoSourceState
}

// Slug returns the command slug.
func (PrviCmd) Slug() string {
	return "PrvI"
}

// MarshalBinary serializes the command to binary format.
func (p *PrviCmd) MarshalBinary() ([]byte, error) {
	// for whatever reason this is 8 bytes long
	pl := make([]byte, 8)
	mes, err := p.videoSourceState.MarshalBinary()
	if err != nil {
		return nil, err
	}

	copy(pl, mes)
	return pl, nil
}

// AuxsCmd represents an aux source command.
type AuxsCmd struct {
	videoSourceState
}

// Slug returns the command slug.
func (AuxsCmd) Slug() string {
	return "AuxS"
}

// KeOnCmd represents the Keyer On Air command
type KeOnCmd struct {
	ME    uint8
	Keyer uint8
	OnAir uint8
}

// Slug returns the command slug.
func (KeOnCmd) Slug() string {
	return "KeOn"
}

// UnmarshalBinary deserializes the command from binary format.
func (k *KeOnCmd) UnmarshalBinary(data []byte) error {
	if len(data) < 3 {
		return errors.New("KeOnCmd: insufficient data")
	}
	k.ME = data[0]
	k.Keyer = data[1]
	k.OnAir = data[2]
	return nil
}

// MarshalBinary serializes the command to binary format.
func (k *KeOnCmd) MarshalBinary() ([]byte, error) {
	return []byte{k.ME, k.Keyer, k.OnAir}, nil
}

// KeBPCmd represents the Keyer Base Properties command
type KeBPCmd struct {
	ME         uint8
	Keyer      uint8
	FillSource models.VideoSource
	KeySource  models.VideoSource
}

// Slug returns the command slug.
func (KeBPCmd) Slug() string {
	return "KeBP"
}

// UnmarshalBinary deserializes the command from binary format.
func (k *KeBPCmd) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return errors.New("KeBPCmd: insufficient data")
	}
	k.ME = data[0]
	k.Keyer = data[1]
	k.FillSource = models.VideoSource(binary.BigEndian.Uint16(data[2:4]))
	k.KeySource = models.VideoSource(binary.BigEndian.Uint16(data[4:6]))
	return nil
}

// MarshalBinary serializes the command to binary format.
func (k *KeBPCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 6)
	pl[0] = k.ME
	pl[1] = k.Keyer
	binary.BigEndian.PutUint16(pl[2:4], uint16(k.FillSource))
	binary.BigEndian.PutUint16(pl[4:6], uint16(k.KeySource))
	return pl, nil
}

// KeDVCmd represents the Keyer DVE command
// DVE (Digital Video Effects) properties for a keyer
type KeDVCmd struct {
	ME            uint8
	Keyer         uint8
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

// Slug returns the command slug.
func (KeDVCmd) Slug() string {
	return "KeDV"
}

// UnmarshalBinary deserializes the command from binary format.
func (k *KeDVCmd) UnmarshalBinary(data []byte) error {
	if len(data) < 30 {
		return errors.New("KeDVCmd: insufficient data")
	}
	k.ME = data[0]
	k.Keyer = data[1]
	k.Rate = data[2]
	// Skip byte 3 (padding/reserved)
	k.SizeX = int16(binary.BigEndian.Uint16(data[4:6]))       //nolint:gosec // intentional conversion for binary protocol
	k.SizeY = int16(binary.BigEndian.Uint16(data[6:8]))       //nolint:gosec // intentional conversion for binary protocol
	k.PositionX = int16(binary.BigEndian.Uint16(data[8:10]))  //nolint:gosec // intentional conversion for binary protocol
	k.PositionY = int16(binary.BigEndian.Uint16(data[10:12])) //nolint:gosec // intentional conversion for binary protocol
	k.Rotation = int16(binary.BigEndian.Uint16(data[12:14]))  //nolint:gosec // intentional conversion for binary protocol

	// Border properties
	borderFlags := data[14]
	k.BorderEnabled = (borderFlags & 0x01) != 0
	k.BorderStyle = data[15]
	k.BorderWidth = int16(binary.BigEndian.Uint16(data[16:18])) //nolint:gosec // intentional conversion for binary protocol
	k.BorderOpacity = data[18]

	// Shadow and light source
	shadowFlags := data[19]
	k.ShadowEnabled = (shadowFlags & 0x01) != 0
	k.LightSource = data[20]

	// Mask properties
	maskFlags := data[21]
	k.MaskEnabled = (maskFlags & 0x01) != 0
	k.MaskTop = int16(binary.BigEndian.Uint16(data[22:24]))    //nolint:gosec // intentional conversion for binary protocol
	k.MaskBottom = int16(binary.BigEndian.Uint16(data[24:26])) //nolint:gosec // intentional conversion for binary protocol
	k.MaskLeft = int16(binary.BigEndian.Uint16(data[26:28]))   //nolint:gosec // intentional conversion for binary protocol
	k.MaskRight = int16(binary.BigEndian.Uint16(data[28:30]))  //nolint:gosec // intentional conversion for binary protocol

	return nil
}

// MarshalBinary serializes the command to binary format.
func (k *KeDVCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 30)
	pl[0] = k.ME
	pl[1] = k.Keyer
	pl[2] = k.Rate
	// Byte 3 is padding/reserved
	binary.BigEndian.PutUint16(pl[4:6], uint16(k.SizeX))       //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[6:8], uint16(k.SizeY))       //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[8:10], uint16(k.PositionX))  //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[10:12], uint16(k.PositionY)) //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[12:14], uint16(k.Rotation))  //nolint:gosec // intentional conversion for binary protocol

	// Border properties
	if k.BorderEnabled {
		pl[14] |= 0x01
	}
	pl[15] = k.BorderStyle
	binary.BigEndian.PutUint16(pl[16:18], uint16(k.BorderWidth)) //nolint:gosec // intentional conversion for binary protocol
	pl[18] = k.BorderOpacity

	// Shadow and light source
	if k.ShadowEnabled {
		pl[19] |= 0x01
	}
	pl[20] = k.LightSource

	// Mask properties
	if k.MaskEnabled {
		pl[21] |= 0x01
	}
	binary.BigEndian.PutUint16(pl[22:24], uint16(k.MaskTop))    //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[24:26], uint16(k.MaskBottom)) //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[26:28], uint16(k.MaskLeft))   //nolint:gosec // intentional conversion for binary protocol
	binary.BigEndian.PutUint16(pl[28:30], uint16(k.MaskRight))  //nolint:gosec // intentional conversion for binary protocol

	return pl, nil
}
