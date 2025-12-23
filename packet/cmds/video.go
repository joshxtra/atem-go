package cmds

import (
	"encoding/binary"

	"github.com/mraerino/atem-go/models"
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

// KedvCmd represents a keyer DVE command.
type KedvCmd struct {
	ME            uint8
	Keyer         uint8
	Rate          uint8
	SizeX         int16
	SizeY         int16
	PositionX     int16
	PositionY     int16
	Rotation      int16
	BorderEnabled bool
	BorderOuterWidth int16
	BorderInnerWidth int16
	BorderOuterSoftness uint8
	BorderInnerSoftness uint8
	BorderBevelSoftness  uint8
	BorderBevelPosition  uint8
	BorderBevel          bool
	BorderHue            uint16
	BorderSaturation     uint16
	BorderLuminance      uint16
	LightSourceDirection uint16
	LightSourceAltitude  uint16
	MaskEnabled          bool
	MaskTop              int16
	MaskBottom           int16
	MaskLeft             int16
	MaskRight            int16
	KeyFrame             bool
}

// Slug returns the command slug.
func (KedvCmd) Slug() string {
	return "KeDV"
}

// MarshalBinary serializes the command to binary format.
func (k *KedvCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 48)
	pl[0] = k.ME
	pl[1] = k.Keyer
	pl[2] = k.Rate
	
	binary.BigEndian.PutUint16(pl[4:], uint16(k.SizeX))
	binary.BigEndian.PutUint16(pl[6:], uint16(k.SizeY))
	binary.BigEndian.PutUint16(pl[8:], uint16(k.PositionX))
	binary.BigEndian.PutUint16(pl[10:], uint16(k.PositionY))
	binary.BigEndian.PutUint16(pl[12:], uint16(k.Rotation))
	
	if k.BorderEnabled {
		pl[14] |= 1
	}
	
	binary.BigEndian.PutUint16(pl[16:], uint16(k.BorderOuterWidth))
	binary.BigEndian.PutUint16(pl[18:], uint16(k.BorderInnerWidth))
	pl[20] = k.BorderOuterSoftness
	pl[21] = k.BorderInnerSoftness
	pl[22] = k.BorderBevelSoftness
	pl[23] = k.BorderBevelPosition
	
	if k.BorderBevel {
		pl[24] |= 1
	}
	
	binary.BigEndian.PutUint16(pl[26:], k.BorderHue)
	binary.BigEndian.PutUint16(pl[28:], k.BorderSaturation)
	binary.BigEndian.PutUint16(pl[30:], k.BorderLuminance)
	binary.BigEndian.PutUint16(pl[32:], k.LightSourceDirection)
	binary.BigEndian.PutUint16(pl[34:], k.LightSourceAltitude)
	
	if k.MaskEnabled {
		pl[36] |= 1
	}
	
	binary.BigEndian.PutUint16(pl[38:], uint16(k.MaskTop))
	binary.BigEndian.PutUint16(pl[40:], uint16(k.MaskBottom))
	binary.BigEndian.PutUint16(pl[42:], uint16(k.MaskLeft))
	binary.BigEndian.PutUint16(pl[44:], uint16(k.MaskRight))
	
	if k.KeyFrame {
		pl[46] |= 1
	}
	
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (k *KedvCmd) UnmarshalBinary(data []byte) error {
	if len(data) < 48 {
		// Pad with zeros if data is shorter
		padded := make([]byte, 48)
		copy(padded, data)
		data = padded
	}
	
	k.ME = data[0]
	k.Keyer = data[1]
	k.Rate = data[2]
	
	k.SizeX = int16(binary.BigEndian.Uint16(data[4:]))
	k.SizeY = int16(binary.BigEndian.Uint16(data[6:]))
	k.PositionX = int16(binary.BigEndian.Uint16(data[8:]))
	k.PositionY = int16(binary.BigEndian.Uint16(data[10:]))
	k.Rotation = int16(binary.BigEndian.Uint16(data[12:]))
	
	k.BorderEnabled = data[14]&1 != 0
	
	k.BorderOuterWidth = int16(binary.BigEndian.Uint16(data[16:]))
	k.BorderInnerWidth = int16(binary.BigEndian.Uint16(data[18:]))
	k.BorderOuterSoftness = data[20]
	k.BorderInnerSoftness = data[21]
	k.BorderBevelSoftness = data[22]
	k.BorderBevelPosition = data[23]
	
	k.BorderBevel = data[24]&1 != 0
	
	k.BorderHue = binary.BigEndian.Uint16(data[26:])
	k.BorderSaturation = binary.BigEndian.Uint16(data[28:])
	k.BorderLuminance = binary.BigEndian.Uint16(data[30:])
	k.LightSourceDirection = binary.BigEndian.Uint16(data[32:])
	k.LightSourceAltitude = binary.BigEndian.Uint16(data[34:])
	
	k.MaskEnabled = data[36]&1 != 0
	
	k.MaskTop = int16(binary.BigEndian.Uint16(data[38:]))
	k.MaskBottom = int16(binary.BigEndian.Uint16(data[40:]))
	k.MaskLeft = int16(binary.BigEndian.Uint16(data[42:]))
	k.MaskRight = int16(binary.BigEndian.Uint16(data[44:]))
	
	k.KeyFrame = data[46]&1 != 0
	
	return nil
}
