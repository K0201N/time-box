# time-box

Time-box is a simple command-line Pomodoro timer written in Go.
It helps you focus by cycling through work and break periods with
a small notification at the end of each phase.

## Requirements

* Go 1.24 or later

## Building and Running

Use `go run` for quick execution or build/install the binary with
`go build` or `go install`:

```shell
# Run directly
$ go run ./cmd/time-box start

# Build and run
$ go build -o time-box ./cmd/time-box
$ ./time-box start -w 25 -b 5 -c 4

# Install to $GOBIN for easier access
$ go install ./cmd/time-box
```

The example above starts a 25‑minute work session followed by a
5‑minute break for one cycle. Adjust the `-w`, `-b` and `-c` flags for
your preferred durations and number of cycles.

## Directory Structure

```shell
cmd/time-box/     # Cobra CLI entry point and commands
internal/notify/  # Cross‑platform notification helpers
internal/timer/   # Core timer logic and tests
```

Main CLI behavior lives in `cmd/time-box/root.go` and `start.go`.
Timer logic is in `internal/timer/timer.go` with tests alongside in
`timer_test.go`. Notification implementations are in `internal/notify`.

The CLI command `time-box start` drives the timer using the
`internal/timer` package. Notifications with a small sound are
implemented in `internal/notify`.

## Running Tests

Unit tests live under `internal/` and can be executed with:

```shell
$ go test ./...

# Run static analysis
$ go vet ./...

# Run tests for the timer package only
$ go test ./internal/timer
```