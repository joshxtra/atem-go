package cmds

import (
	"encoding/binary"

	"github.com/krombel/atem-go/models"
)

// TlinCmd represents a tally by index command.
type TlinCmd map[int]models.TallyState

// Slug returns the command slug.
func (TlinCmd) Slug() string {
	return "TlIn"
}

// MarshalBinary serializes the command to binary format.
func (t TlinCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, len(t)+2)
	binary.BigEndian.PutUint16(pl, uint16(len(t))) //nolint:gosec // map length is always within uint16 range

	for idx, state := range t {
		pl[idx] = state.Bitmask()
	}

	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (t TlinCmd) UnmarshalBinary(data []byte) error {
	num := binary.BigEndian.Uint16(data)

	for i := 0; i < int(num); i++ {
		t[i] = models.TallyState{
			Program: data[2+i]&1 != 0,
			Preview: data[2+i]&2 != 0,
		}
	}

	return nil
}

// TlsrCmd represents a tally by source command.
type TlsrCmd map[models.VideoSource]models.TallyState

// Slug returns the command slug.
func (TlsrCmd) Slug() string {
	return "TlSr"
}

// MarshalBinary serializes the command to binary format.
func (t TlsrCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, len(t)*3+2)
	binary.BigEndian.PutUint16(pl, uint16(len(t))) //nolint:gosec // map length is always within uint16 range

	var i int
	for idx, state := range t {
		start := 2 + i*3
		binary.BigEndian.PutUint16(pl[start:], uint16(idx))
		pl[start+2] = state.Bitmask()
		i++
	}

	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (t TlsrCmd) UnmarshalBinary(data []byte) error {
	num := binary.BigEndian.Uint16(data)

	for i := 0; i < int(num); i++ {
		start := 2 + i*3
		src := models.VideoSource(binary.BigEndian.Uint16(data[start:]))
		t[src] = models.TallyState{
			Program: data[start+2]&1 != 0,
			Preview: data[start+2]&2 != 0,
		}
	}

	return nil
}
