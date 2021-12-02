package models

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
