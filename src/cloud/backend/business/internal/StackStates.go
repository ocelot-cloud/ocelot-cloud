package internal

type StackState int

const (
	Uninitialized StackState = iota
	Running
	Starting
	Available
	Downloading
	Stopping
)

func (s *StackState) String() string {
	return [...]string{"Uninitialized", "Running", "Starting", "Available", "Downloading", "Stopping"}[*s]
}
