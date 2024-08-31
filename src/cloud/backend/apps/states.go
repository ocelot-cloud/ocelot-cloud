package apps

type appState int

const (
	Uninitialized appState = iota
	Running
	Starting
	Available
	Downloading
	Stopping
)

func (s *appState) toString() string {
	return [...]string{"Uninitialized", "Running", "Starting", "Available", "Downloading", "Stopping"}[*s]
}
