package models

//go:generate go run golang.org/x/tools/cmd/stringer -type=MediaPlayerType -linecomment -output=media_string.go

// MediaPlayerType represents the type of media player.
type MediaPlayerType uint8

// Media player type constants.
const (
	MediaPlayerTypeStill MediaPlayerType = 1 + iota // Still
	MediaPlayerTypeClip                             // Clip
)

// MediaPlayer represents a media player configuration.
type MediaPlayer struct {
	Type       MediaPlayerType
	StillIndex int
	ClipIndex  int
	// missing: clip player state (RCPS)
}

// MediaStillFrame represents a still frame in the media pool.
type MediaStillFrame struct {
	Used     bool
	Hash     []byte
	Filename string
}
