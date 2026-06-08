package repository

import "fmt"

type StatusOrder int

const (
	StatusNew StatusOrder = iota
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func (s StatusOrder) String() string {
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
