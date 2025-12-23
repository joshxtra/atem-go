package models

// InputProperties describes the properties of a video input source.
type InputProperties struct {
	SourceIndex VideoSource
	LongName    string
	ShortName   string

	ExternalPortType ExternalPortType
	PortType         PortType
}

// Generate String() methods
//go:generate go run golang.org/x/tools/cmd/stringer -type=ExternalPortType,PortType -linecomment -output=inputs_string.go

// ExternalPortType represents the type of external port.
type ExternalPortType uint8

// External port type constants.
const (
	ExternalPortTypeInternal  ExternalPortType = iota // Internal
	ExternalPortTypeSDI                               // SDI
	ExternalPortTypeHDMI                              // HDMI
	ExternalPortTypeComposite                         // Composite
	ExternalPortTypeComponent                         // Component
	ExternalPortTypeSVideo                            // SVideo
)

// PortType represents the type of port.
type PortType uint8

// Port type constants.
const (
	PortTypeExternal        PortType = iota // External
	PortTypeBlack                           // Black
	PortTypeColorBars                       // Color Bars
	PortTypeColorGenerator                  // Color Generator
	PortTypeMediaPlayerFill                 // Media Player Fill
	PortTypeMediaPlayerKey                  // Media Player Key
	PortTypeSuperSource                     // SuperSource
)

// Port type constants for ME outputs and auxiliaries.
const (
	PortTypeMEOutput  PortType = 128 + iota // ME Output
	PortTypeAuxiliary                       // Auxiliary
	PortTypeMask                            // Mask
)
