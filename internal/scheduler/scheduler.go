package scheduler

import (
	"github.com/robfig/cron/v3"
)

// ScheduleTask schedules a given task to run at the specified cron schedule.
func ScheduleTask(schedule string, task func()) (cron.EntryID, error) {
	c := cron.New()

	id, err := c.AddFunc(schedule, task)
	if err != nil {
		return 0, err
	}

	c.Start()
	return id, nil
}
