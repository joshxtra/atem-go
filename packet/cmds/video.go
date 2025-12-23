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
