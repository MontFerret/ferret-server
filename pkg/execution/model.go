package execution

import (
	"github.com/MontFerret/ferret-server/pkg/scripts"
	"strings"
	"time"
)

type (
	// Job state
	Status uint

	// Job cause
	Cause uint

	// Event represents an event that triggered a script execution
	Event struct {
		ProjectID string
		ScriptID  string
		CausedBy  Cause
	}

	// Job represents a running script
	Job struct {
		ID        string
		ProjectID string
		CausedBy  Cause
		Script    scripts.ScriptEntity
	}

	State struct {
		Job       Job
		Timestamp time.Time
		Status    Status
		Error     error
	}

	Result struct {
		State
		Data []byte
	}
)

const (
	StatusUnknown   Status = 0
	StatusQueued    Status = 1
	StatusRunning   Status = 2
	StatusCompleted Status = 3
	StatusCancelled Status = 4
	StatusErrored   Status = 5

	CauseUnknown  Cause = 0
	CauseManual   Cause = 1
	CauseSchedule Cause = 2
	CauseHook     Cause = 3
)

var statusstr = map[Status]string{
	StatusUnknown:   "unknown",
	StatusQueued:    "queued",
	StatusRunning:   "running",
	StatusCompleted: "completed",
	StatusCancelled: "cancelled",
	StatusErrored:   "errored",
}

var causestr = map[Cause]string{
	CauseUnknown:  "unknown",
	CauseManual:   "manual",
	CauseSchedule: "schedule",
	CauseHook:     "hook",
}

func NewStatus(input string) Status {
	input = strings.ToLower(input)

	for res, str := range statusstr {
		if input == str {
			return res
		}
	}

	return StatusUnknown
}

func (t Status) String() string {
	return statusstr[t]
}

func NewCause(input string) Cause {
	input = strings.ToLower(input)

	for res, str := range causestr {
		if input == str {
			return res
		}
	}

	return CauseUnknown
}

func (t Cause) String() string {
	return causestr[t]
}
