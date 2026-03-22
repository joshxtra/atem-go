package cmds

import (
	"encoding/binary"
	"errors"

	"github.com/joshxtra/atem-go/models"
)

// MpceCmd represents a media player source command.
type MpceCmd struct {
	PlayerIndex uint8
	Type        models.MediaPlayerType
	StillIndex  int
	ClipIndex   int
}

// Slug returns the command slug.
func (MpceCmd) Slug() string {
	return "MPCE"
}

// MarshalBinary serializes the command to binary format.
func (m *MpceCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = m.PlayerIndex
	pl[1] = uint8(m.Type)
	pl[2] = uint8(m.StillIndex) //nolint:gosec // StillIndex is always within uint8 range
	pl[3] = uint8(m.ClipIndex)  //nolint:gosec // ClipIndex is always within uint8 range
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MpceCmd) UnmarshalBinary(data []byte) error {
	m.PlayerIndex = data[0]
	m.Type = models.MediaPlayerType(data[1])
	m.StillIndex = int(data[2])
	m.ClipIndex = int(data[3])
	return nil
}

// MpfeCmd represents a media player file command.
type MpfeCmd struct {
	Index    uint16
	Used     bool
	Hash     []byte
	Filename string
}

// Slug returns the command slug.
func (MpfeCmd) Slug() string {
	return "MPfe"
}

// MarshalBinary serializes the command to binary format.
func (m *MpfeCmd) MarshalBinary() ([]byte, error) {
	filenameBytes := []byte(m.Filename)
	pl := make([]byte, 24+len(filenameBytes))

	pl[0] = 0 // still type
	binary.BigEndian.PutUint16(pl[2:], m.Index)
	if m.Used {
		pl[4] = 1
	}

	copy(pl[5:], m.Hash[:16])

	pl[23] = uint8(len(filenameBytes)) //nolint:gosec // filename length is always within uint8 range
	copy(pl[24:], filenameBytes)

	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MpfeCmd) UnmarshalBinary(data []byte) error {
	if data[0] != 0 {
		return errors.New("unknown file type")
	}

	m.Index = binary.BigEndian.Uint16(data[2:])
	m.Used = data[4] == 1
	m.Hash = data[5:21]

	filenameLen := data[23]
	m.Filename = string(data[24 : 24+filenameLen])

	return nil
}
