// Package cmds provides command types for the ATEM protocol.
package cmds

import (
	"encoding/binary"

	"github.com/krombel/atem-go/models"
)

// UnknownCommand represents a command that is not yet implemented.
type UnknownCommand struct {
	slug string
	data []byte
}

// NewUnknownCommand creates a new unknown command with the given slug.
func NewUnknownCommand(slug string) *UnknownCommand {
	return &UnknownCommand{slug: slug}
}

// Slug returns the command slug.
func (u *UnknownCommand) Slug() string {
	return u.slug
}

// MarshalBinary serializes the command to binary format.
func (u *UnknownCommand) MarshalBinary() ([]byte, error) {
	return u.data, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (u *UnknownCommand) UnmarshalBinary(data []byte) error {
	u.data = data
	return nil
}

// VerCmd represents the _ver command
type VerCmd struct {
	Major uint16
	Minor uint16
}

// Slug returns the command slug.
func (VerCmd) Slug() string {
	return "_ver"
}

// MarshalBinary serializes the command to binary format.
func (c *VerCmd) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf, c.Major)
	binary.BigEndian.PutUint16(buf[2:], c.Minor)
	return buf, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (c *VerCmd) UnmarshalBinary(data []byte) error {
	c.Major = binary.BigEndian.Uint16(data)
	c.Minor = binary.BigEndian.Uint16(data[2:])
	return nil
}

// IncmCmd represents the initialization complete command.
type IncmCmd struct{}

// Slug returns the command slug.
func (*IncmCmd) Slug() string {
	return "InCm"
}

// MarshalBinary serializes the command to binary format.
func (*IncmCmd) MarshalBinary() ([]byte, error) {
	return make([]byte, 4), nil // todo: handle states
}

// UnmarshalBinary deserializes the command from binary format.
func (*IncmCmd) UnmarshalBinary(_ []byte) error {
	return nil
}

// TimeCmd represents a timecode command.
type TimeCmd models.Timecode

// Slug returns the command slug.
func (TimeCmd) Slug() string {
	return "Time"
}

// MarshalBinary serializes the command to binary format.
func (t *TimeCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 8)
	pl[0] = uint8(t.Hour)   //nolint:gosec // time values are always within uint8 range
	pl[1] = uint8(t.Minute) //nolint:gosec // time values are always within uint8 range
	pl[2] = uint8(t.Second) //nolint:gosec // time values are always within uint8 range
	pl[3] = uint8(t.Frame)  //nolint:gosec // time values are always within uint8 range
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (t *TimeCmd) UnmarshalBinary(data []byte) error {
	t.Hour = int(data[0])
	t.Minute = int(data[1])
	t.Second = int(data[2])
	t.Frame = int(data[3])
	return nil
}

// TiRqCmd represents a timecode request command.
type TiRqCmd struct{}

// Slug returns the command slug.
func (TiRqCmd) Slug() string {
	return "TiRq"
}

// MarshalBinary serializes the command to binary format.
func (TiRqCmd) MarshalBinary() ([]byte, error) {
	return make([]byte, 0), nil
}

// UnmarshalBinary deserializes the command from binary format.
func (TiRqCmd) UnmarshalBinary(_ []byte) error {
	return nil
}
