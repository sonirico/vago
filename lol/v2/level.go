package lol

// Level is a severity; lower is more severe.
type Level uint8

const (
	LevelPanic Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)
