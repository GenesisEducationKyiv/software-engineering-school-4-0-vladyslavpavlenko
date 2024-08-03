package scheduler_test

import (
	"testing"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/scheduler"
)

func TestNewCronScheduler(t *testing.T) {
	s := scheduler.NewCronScheduler()
	if s.Cron == nil {
		t.Fatal("Expected non-nil cron scheduler")
	}
}

func TestScheduleTask(t *testing.T) {
	s := scheduler.NewCronScheduler()
	_, err := s.Schedule("* * * * *", func() { t.Log("Task executed") })
	if err != nil {
		t.Errorf("Failed to schedule task with valid cron schedule: %v", err)
	}

	_, err = s.Schedule("invalid schedule", func() {})
	if err == nil {
		t.Error("Expected error when scheduling task with invalid cron schedule, got none")
	}
}

func TestStart(t *testing.T) {
	s := scheduler.NewCronScheduler()
	done := make(chan bool)
	wasRun := false

	_, err := s.Schedule("@every 1s", func() {
		wasRun = true
		done <- true
	})
	if err != nil {
		t.Log("Error scheduling task:", err)
		return
	}

	go s.Start()
	select {
	case <-done:
		if !wasRun {
			t.Error("Scheduled task was not run")
		}
	case <-time.After(2 * time.Second):
		t.Error("Scheduled task was not run within the expected time")
	}
	s.Cron.Stop()
}
