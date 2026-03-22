package cmds

import (
	"bytes"
	"encoding/binary"

	"github.com/joshxtra/atem-go/models"
)

type statusStringCommand string

func (s *statusStringCommand) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 44)
	copy(pl, []byte(*s))
	return pl, nil
}

func (s *statusStringCommand) UnmarshalBinary(data []byte) error {
	*s = statusStringCommand(decodeString(data))
	return nil
}

func (s *statusStringCommand) Value() string {
	return string(*s)
}

// PinCmd represents the product name command.
type PinCmd struct {
	statusStringCommand
}

// Slug returns the command slug.
func (PinCmd) Slug() string {
	return "_pin"
}

// WarnCmd represents a warning message command.
type WarnCmd struct {
	statusStringCommand
}

// Slug returns the command slug.
func (WarnCmd) Slug() string {
	return "Warn"
}

// TopCmd represents the topology command.
type TopCmd models.Topology

// Slug returns the command slug.
func (TopCmd) Slug() string {
	return "_top"
}

// MarshalBinary serializes the command to binary format.
func (t *TopCmd) MarshalBinary() ([]byte, error) {
	fields := []uint8{
		t.MEs, t.Sources, t.ColorGenerators,
		t.AUXBusses, t.DownstreamKeyers, t.Stingers,
		t.DVEs, t.SuperSources,
		1, 0, 1, 0, // unknown
	}
	if t.SDOutput {
		fields[9] = 1
	}

	pl := make([]byte, 0, 12)
	buf := bytes.NewBuffer(pl)
	if err := binary.Write(buf, binary.BigEndian, fields); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary deserializes the command from binary format.
func (t *TopCmd) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	var hasSDOut uint8
	fields := []interface{}{
		&t.MEs, &t.Sources, &t.ColorGenerators,
		&t.AUXBusses, &t.DownstreamKeyers, &t.Stingers,
		&t.DVEs, &t.SuperSources,
		new(uint8), &hasSDOut,
	}

	err := decode(buf, fields)
	if err != nil {
		return err
	}

	if hasSDOut > 0 {
		t.SDOutput = true
	}

	return nil
}

// MecCmd represents the mix effect config command.
type MecCmd struct {
	ME     uint8
	Keyers uint8
}

// Slug returns the command slug.
func (MecCmd) Slug() string {
	return "_MeC"
}

// MarshalBinary serializes the command to binary format.
func (m *MecCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = m.ME
	pl[1] = m.Keyers
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MecCmd) UnmarshalBinary(data []byte) error {
	m.ME = data[0]
	m.Keyers = data[1]
	return nil
}

// MplCmd represents the media player config command.
type MplCmd models.MediaPlayerConfig

// Slug returns the command slug.
func (MplCmd) Slug() string {
	return "_mpl"
}

// MarshalBinary serializes the command to binary format.
func (m *MplCmd) MarshalBinary() ([]byte, error) {
	fields := []uint8{
		m.StillBanks, m.ClipBanks,
		0, 0, // unknown
	}
	pl := make([]byte, 0, 4)
	buf := bytes.NewBuffer(pl)
	if err := binary.Write(buf, binary.BigEndian, fields); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MplCmd) UnmarshalBinary(data []byte) error {
	fields := []interface{}{
		&m.StillBanks, &m.ClipBanks,
	}
	buf := bytes.NewBuffer(data)
	return decode(buf, fields)
}

// MvcCmd represents the multiview config command.
type MvcCmd uint8

// Slug returns the command slug.
func (MvcCmd) Slug() string {
	return "_MvC"
}

// MarshalBinary serializes the command to binary format.
func (m *MvcCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = uint8(*m)
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MvcCmd) UnmarshalBinary(data []byte) error {
	*m = MvcCmd(data[0])
	return nil
}

// SscCmd represents the super source config command.
type SscCmd uint8

// Slug returns the command slug.
func (SscCmd) Slug() string {
	return "_SSC"
}

// MarshalBinary serializes the command to binary format.
func (m *SscCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = uint8(*m)
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *SscCmd) UnmarshalBinary(data []byte) error {
	*m = SscCmd(data[0])
	return nil
}

// TlcCmd represents the tally channel config command.
type TlcCmd uint8

// Slug returns the command slug.
func (TlcCmd) Slug() string {
	return "_TlC"
}

// MarshalBinary serializes the command to binary format.
func (t *TlcCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 8)
	pl[4] = uint8(*t)
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (t *TlcCmd) UnmarshalBinary(data []byte) error {
	*t = TlcCmd(data[4])
	return nil
}

// MacCmd represents the macro config command.
type MacCmd uint8

// Slug returns the command slug.
func (MacCmd) Slug() string {
	return "_MAC"
}

// MarshalBinary serializes the command to binary format.
func (m *MacCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = uint8(*m)
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *MacCmd) UnmarshalBinary(data []byte) error {
	*m = MacCmd(data[0])
	return nil
}

// PowrCmd represents the power status command.
type PowrCmd models.PowerStatus

// Slug returns the command slug.
func (PowrCmd) Slug() string {
	return "Powr"
}

// MarshalBinary serializes the command to binary format.
func (p *PowrCmd) MarshalBinary() ([]byte, error) {
	var status uint8
	if p.Main {
		status |= 1
	}
	if p.Backup {
		status |= 2
	}

	pl := make([]byte, 4)
	pl[0] = status

	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (p *PowrCmd) UnmarshalBinary(data []byte) error {
	status := data[0]
	if status&1 != 0 {
		p.Main = true
	}
	if status&2 != 0 {
		p.Backup = true
	}
	return nil
}

// VidmCmd represents the video mode command.
type VidmCmd models.VideoMode

// Slug returns the command slug.
func (VidmCmd) Slug() string {
	return "VidM"
}

// MarshalBinary serializes the command to binary format.
func (m *VidmCmd) MarshalBinary() ([]byte, error) {
	pl := make([]byte, 4)
	pl[0] = uint8(*m)
	return pl, nil
}

// UnmarshalBinary deserializes the command from binary format.
func (m *VidmCmd) UnmarshalBinary(data []byte) error {
	*m = VidmCmd(data[0])
	return nil
}
