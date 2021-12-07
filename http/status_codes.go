package http

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
	Ready      StatusCode = 4016
)

var statusCodeNamesMap = map[StatusCode]string{
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
	Ready:      "Ready",
}

func (sc StatusCode) String() string {
	return strconv.Itoa(int(sc))
}

func (sc StatusCode) StatusCodeName() string {
	return statusCodeNamesMap[sc]
}
