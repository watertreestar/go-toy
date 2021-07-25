package regsitry

type Registration struct {
	ServiceName      serviceName
	ServiceURL       string
	RequiredService  []serviceName
	ServiceUpdateURL string
}

type serviceName string

const (
	LogService = serviceName("LogService")
)

type PatchEntry struct {
	Name serviceName
	URL  string
}

type patch struct {
	Added   []PatchEntry
	Removed []PatchEntry
}
