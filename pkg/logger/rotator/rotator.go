package rotator

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

// A Rotator represents an active rotating object that uses lumberjack.Logger
// to rotate log files.
type Rotator struct {
	Logger *lumberjack.Logger
}

// New creates a new [Rotator].
func New() *Rotator {
	var r Rotator
	r.Logger = &lumberjack.Logger{
		Filename:   "/var/log/app/app.log",
		MaxSize:    5,  // megabytes
		MaxBackups: 10, // number of backups
		MaxAge:     14, // days
		Compress:   true,
	}

	return &r
}
