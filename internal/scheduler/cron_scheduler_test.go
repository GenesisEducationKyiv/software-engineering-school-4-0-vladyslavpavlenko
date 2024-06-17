package scheduler_test

import (
	"testing"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/scheduler"
)

func TestNewCronScheduler(t *testing.T) {
	s := scheduler.NewCronScheduler()
	if s.Cron == nil {
		t.Fatal("Expected non-nil cron scheduler")
	}
}

func TestScheduleTask(t *testing.T) {
	s := scheduler.NewCronScheduler()
	_, err := s.ScheduleTask("* * * * *", func() { t.Log("Task executed") })
	if err != nil {
		t.Errorf("Failed to schedule task with valid cron schedule: %v", err)
	}

	_, err = s.ScheduleTask("invalid schedule", func() {})
	if err == nil {
		t.Error("Expected error when scheduling task with invalid cron schedule, got none")
	}
}

func TestStart(t *testing.T) {
	s := scheduler.NewCronScheduler()
	done := make(chan bool)
	wasRun := false

	_, err := s.ScheduleTask("@every 1s", func() {
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

func TestStop(t *testing.T) {
	s := scheduler.NewCronScheduler()
	timesRun := 0
	stopTest := make(chan bool)
	taskFinished := make(chan bool)

	_, err := s.ScheduleTask("@every 1s", func() {
		timesRun++
		taskFinished <- true
	})
	if err != nil {
		t.Fatalf("Error scheduling task: %v", err)
	}

	go s.Start()

	select {
	case <-taskFinished:
	case <-time.After(2 * time.Second):
		t.Fatal("Task did not run within expected time")
	}

	s.Stop()

	time.Sleep(2 * time.Second)

	if timesRun != 1 {
		t.Errorf("Task ran %d times; expected to run exactly once after stop", timesRun)
	}

	stopTest <- true
}
