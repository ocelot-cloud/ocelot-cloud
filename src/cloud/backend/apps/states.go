package apps

type stackState int

const (
	Uninitialized stackState = iota
	Running
	Starting
	Available
	Downloading
	Stopping
)

func (s *stackState) toString() string {
	return [...]string{"Uninitialized", "Running", "Starting", "Available", "Downloading", "Stopping"}[*s]
}
