package main

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/K0201N/time-box/internal/timer"
)

func TestRunStartContinuesAfterNotificationErrors(t *testing.T) {
	var stdout, stderr bytes.Buffer
	var messages []string
	push := func(title, message string) error {
		messages = append(messages, message)
		return errors.New("notification unavailable")
	}

	err := runStart(context.Background(), []timer.Phase{
		{Label: "Work", Duration: time.Second},
		{Label: "Break", Duration: time.Second},
	}, 1, &stdout, &stderr, push)
	if err != nil {
		t.Fatalf("runStart returned an error: %v", err)
	}
	wantMessages := []string{
		"Timer started!",
		"Work done!",
		"Break done! All cycles completed!",
	}
	if !reflect.DeepEqual(messages, wantMessages) {
		t.Fatalf("notification messages = %q, want %q", messages, wantMessages)
	}
	for _, displayed := range []string{"Work  00:01", "Work  00:00", "Break 00:01", "Break 00:00"} {
		if !strings.Contains(stdout.String(), displayed) {
			t.Errorf("countdown %q missing from %q", displayed, stdout.String())
		}
	}
	if got := strings.Count(stderr.String(), "notification unavailable"); got != len(wantMessages) {
		t.Errorf("notification errors = %d, want %d: %q", got, len(wantMessages), stderr.String())
	}
}

func TestRunStartContinuesAfterSlowNotificationError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	var notificationDelay time.Duration
	var resumedAt time.Time
	push := func(title, message string) error {
		if message == "Work done!" {
			time.Sleep(1050 * time.Millisecond)
			resumedAt = time.Now()
		}
		if message == "Break done! All cycles completed!" {
			notificationDelay = time.Since(resumedAt)
		}
		return errors.New("notification unavailable")
	}

	err := runStart(context.Background(), []timer.Phase{
		{Label: "Work", Duration: time.Second},
		{Label: "Break", Duration: time.Second},
	}, 1, &stdout, &stderr, push)
	if err != nil {
		t.Fatalf("runStart returned an error: %v", err)
	}
	if resumedAt.IsZero() || notificationDelay == 0 {
		t.Fatalf("phase notifications did not complete: %q", stderr.String())
	}
	if notificationDelay > 700*time.Millisecond {
		t.Fatalf("final notification arrived %s after consumer resumed", notificationDelay)
	}
	if !strings.Contains(stdout.String(), "Break 00:00") {
		t.Fatalf("break completion missing from %q", stdout.String())
	}
	if got := strings.Count(stderr.String(), "notification unavailable"); got != 3 {
		t.Errorf("notification errors = %d, want 3: %q", got, stderr.String())
	}
}
