package src

import "github.com/mattn/go-shellwords"

func ParseCommand(command string) ([]string, error) {
	parser := shellwords.NewParser()
	return parser.Parse(command)
}
