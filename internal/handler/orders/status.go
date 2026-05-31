package orders

import "fmt"

type Status int

const (
	StatusNew Status = iota
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func (s Status) String() string {
	switch s {
	case StatusNew:
		return "NEW"
	case StatusProcessing:
		return "PROCESSING"
	case StatusInvalid:
		return "INVALID"
	case StatusProcessed:
		return "PROCESSED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}
