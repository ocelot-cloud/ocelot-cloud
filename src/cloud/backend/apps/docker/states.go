package docker

type AppState int

const (
	Uninitialized AppState = iota
	Running
	Starting
	Available
	Downloading
	Stopping
)

func (s *AppState) ToString() string {
	return [...]string{"Uninitialized", "Running", "Starting", "Available", "Downloading", "Stopping"}[*s]
}
