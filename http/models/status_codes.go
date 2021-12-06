package models

import "strconv"

type StatusCode int

const (
	Queued     StatusCode = 4000
	Processing StatusCode = 4001
	Success    StatusCode = 4002
	Failure    StatusCode = 4003
	Error      StatusCode = 4004
	Waiting    StatusCode = 4005
	Canceled   StatusCode = 4006
	Unstable   StatusCode = 4007
	Skipped    StatusCode = 4008
	TimedOut   StatusCode = 4009
	TimingOut  StatusCode = 4014
	Creating   StatusCode = 4015
)

func (sc StatusCode) String() string {
	return strconv.Itoa(int(sc))
}

func (sc StatusCode) StatusCodeName() string {
	switch sc {
	case Queued:
		return "Queued"
	case Processing:
		return "Processing"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Error:
		return "Error"
	case Waiting:
		return "Waiting"
	case Canceled:
		return "Canceled"
	case Unstable:
		return "Unstable"
	case Skipped:
		return "Skipped"
	case TimedOut:
		return "TimedOut"
	case TimingOut:
		return "TimingOut"
	case Creating:
		return "Creating"
	}
	return ""
}

var StatusCodeNamesMap = map[StatusCode]string{
	Queued:     "Queued",
	Processing: "Processing",
	Success:    "Success",
	Failure:    "Failure",
	Error:      "Error",
	Waiting:    "Waiting",
	Canceled:   "Canceled",
	Unstable:   "Unstable",
	Skipped:    "Skipped",
	TimedOut:   "TimedOut",
	TimingOut:  "TimingOut",
	Creating:   "Creating",
}
