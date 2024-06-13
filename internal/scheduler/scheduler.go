package scheduler

import (
	"github.com/robfig/cron/v3"
)

// Scheduler defines an interface for scheduling tasks.
type Scheduler interface {
	ScheduleTask(schedule string, task func()) (cron.EntryID, error)
	Start()
}
