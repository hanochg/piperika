package report

type Report interface {
	fetchReport()
	toString() string
}

func FetchReport(r Report) {
	r.fetchReport()
}

func ToString(r Report) string {
	return r.toString()
}
