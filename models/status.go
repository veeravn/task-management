package models

type Status int

const (
	Pending    Status = iota // 0
	InProgress               // 1
	Completed                // 2
)

// String provides a string representation for the Status type
func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
	}
}
