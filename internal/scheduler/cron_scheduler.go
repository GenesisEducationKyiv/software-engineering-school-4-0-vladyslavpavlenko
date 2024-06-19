package scheduler

import "github.com/robfig/cron/v3"

// CronScheduler implements TaskScheduler using the robfig/cron package.
type CronScheduler struct {
	Cron *cron.Cron
}

// NewCronScheduler creates a new instance of a CronScheduler.
func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		Cron: cron.New(),
	}
}

// ScheduleTask schedules a given task to run at the specified cron schedule.
func (s *CronScheduler) ScheduleTask(schedule string, task func()) (cron.EntryID, error) {
	id, err := s.Cron.AddFunc(schedule, task)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Start starts the cron scheduler.
func (s *CronScheduler) Start() {
	s.Cron.Start()
}

// Stop stops the cron scheduler.
func (s *CronScheduler) Stop() {
	s.Cron.Stop()
}
