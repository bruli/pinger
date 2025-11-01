package domain

import "errors"

const (
	ReadyStatus Status = iota + 1
	DegradedStatus
	FailStatus
)

var (
	stringToStatusMap = map[string]Status{
		"ready":    ReadyStatus,
		"degraded": DegradedStatus,
		"fail":     FailStatus,
	}
	statusToStringMap = map[Status]string{
		ReadyStatus:    "ready",
		DegradedStatus: "degraded",
		FailStatus:     "fail",
	}

	ErrInvalidStatus = errors.New("invalid status")
)

type Status int

func (s Status) String() string {
	return statusToStringMap[s]
}

func ParseStatus(s string) (Status, error) {
	st, ok := stringToStatusMap[s]
	if !ok {
		return 0, ErrInvalidStatus
	}
	return st, nil
}
