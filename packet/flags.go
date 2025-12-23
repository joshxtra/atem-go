package packet

const (
	// FlagNeedACK indicates the message requires an acknowledgment.
	FlagNeedACK Flags = 1 << iota
	// FlagInit indicates an initialization message.
	FlagInit
	// FlagRetrans indicates a retransmission.
	FlagRetrans
	// FlagHello indicates a hello message.
	FlagHello
	// FlagACK indicates an acknowledgment message.
	FlagACK
)

// Flags represents message flags in the ATEM protocol.
type Flags uint8

// Has checks if the flag has the given mask set.
func (f Flags) Has(mask Flags) bool {
	return uint8(f)&uint8(mask) != 0
}

// Debug returns a map of flag names to their boolean values for debugging.
func (f Flags) Debug() map[string]interface{} {
	flags := map[string]Flags{
		"needACK":        FlagNeedACK,
		"init":           FlagInit,
		"retransmission": FlagRetrans,
		"hello":          FlagHello,
		"ack":            FlagACK,
	}
	fields := make(map[string]interface{})
	for name, val := range flags {
		fields[name] = f.Has(val)
	}
	return fields
}

// FlagsFrom creates a Flags value from multiple flag masks.
func FlagsFrom(masks ...Flags) Flags {
	var flags Flags
	for _, mask := range masks {
		flags |= mask
	}
	return flags
}
