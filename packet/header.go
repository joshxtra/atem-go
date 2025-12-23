package packet

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"

	"github.com/pkg/errors"
)

// Message represents an ATEM protocol message containing commands.
type Message struct {
	Flags     Flags  // max: 5 bit
	Length    uint16 // max: 11 bit
	SessionID uint16
	AckID     uint16
	SeqNum    uint16

	Commands []Command
}

const headerSize = 12
const lengthBitmask = 0xFFFF >> 5 // first 5 bits are preserved for flags

// Serialize writes the message to the buffer.
func (m *Message) Serialize(buf *bytes.Buffer) {
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint16(header[2:], m.SessionID)
	binary.BigEndian.PutUint16(header[4:], m.AckID)
	binary.BigEndian.PutUint16(header[10:], m.SeqNum)

	_, _ = buf.Write(header)

	for _, cmd := range m.Commands {
		pl, _ := MarshalCommand(cmd)
		_, _ = buf.Write(pl)
	}

	// allow custom length
	if m.Length == 0 {
		m.Length = uint16(buf.Len()) //nolint:gosec // buffer length is always within uint16 range
	}

	var cmdAndLen uint16
	cmdAndLen |= uint16(m.Flags) << 11
	cmdAndLen |= m.Length & lengthBitmask
	binary.BigEndian.PutUint16(buf.Bytes(), cmdAndLen)
}

// Deserialize reads a message from the buffer.
func Deserialize(_ *slog.Logger, buf *bytes.Buffer) (Message, error) {
	msg := Message{}
	bytes := buf.Next(headerSize)
	if len(bytes) < headerSize {
		return msg, errors.New("EOF while reading header")
	}

	cmdAndLen := binary.BigEndian.Uint16(bytes)
	flags := Flags(cmdAndLen >> 11) //nolint:gosec // flags are always within uint8 range (5 bits)
	msg.Flags = flags
	msg.Length = cmdAndLen & lengthBitmask

	msg.SessionID = binary.BigEndian.Uint16(bytes[2:])
	msg.AckID = binary.BigEndian.Uint16(bytes[4:])
	msg.SeqNum = binary.BigEndian.Uint16(bytes[10:])

	// no further processing if init packet
	if msg.Flags.Has(FlagInit) {
		return msg, nil
	}

	for {
		var length uint16
		err := binary.Read(buf, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		if err != nil {
			return msg, errors.Wrap(err, "Failed to read command length")
		}
		length -= 2 // remove length field

		pl := buf.Next(int(length))
		if len(pl) < int(length) {
			return msg, errors.New("EOF while reading command")
		}

		cmd, err := UnmarshalCommand(pl)
		if err != nil {
			return msg, err
		}
		msg.Commands = append(msg.Commands, cmd)
	}

	return msg, nil
}
