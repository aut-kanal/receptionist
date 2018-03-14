package receptionist

import (
	"os"
	"time"
)

var (
	Version       string
	Commit        string
	BuildTime     string
	RootDirectory string
	Title         string
	StartTime     time.Time
)

func init() {
	if Version == "" {
		Version = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if RootDirectory == "" {
		RootDirectory = os.Getenv("HOME") + "/receptionist"
	}
	if Title == "" {
		Title = "receptionist"
	}
	StartTime = time.Now()
}
