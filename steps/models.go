package steps

type Step struct {
	Name string `json:"name"`
}

type StatusCode int

const (
	InProgress StatusCode = 0    // TBD
	Wait       StatusCode = 0    // TBD
	Success    StatusCode = 4002 // TBD
	Failure    StatusCode = 4003
	Error      StatusCode = 4004
	Skipped    StatusCode = 4008
)
