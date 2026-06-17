package models

import (
	"encoding/json"
	"fmt"
)

type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusProcessing
	StatusInvalid
	StatusProcessed
)

func (s OrderStatus) String() string {
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

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *OrderStatus) Scan(value any) error {
	switch v := value.(type) {
	case int64:
		*s = OrderStatus(v)
	case int32:
		*s = OrderStatus(v)
	case int:
		*s = OrderStatus(v)
	default:
		return fmt.Errorf("cannot scan %T into OrderStatus", value)
	}
	return nil
}
