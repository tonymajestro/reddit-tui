package utils

import (
	"fmt"
	"log/slog"
	"time"
)

type Timer struct {
	Context string
	Start   time.Time
	End     time.Time
	time.Duration
}

func NewTimer(context string) Timer {
	return Timer{
		Context: context,
		Start:   time.Now(),
	}
}

func (t *Timer) Stop() {
	t.End = time.Now()
	t.Duration = t.End.Sub(t.Start)
}

func (t *Timer) StopAndLog(args ...any) {
	t.Stop()

	slogArgs := []any{
		"duration",
		fmt.Sprintf("%d ms", t.Milliseconds()),
	}

	for _, arg := range args {
		slogArgs = append(slogArgs, any(arg))
	}

	slog.Debug(t.Context, slogArgs...)
}
