package sockets

type EventType string

const (
	Ack EventType = "ack"
)

func (self EventType) Valid() bool {
	switch self {
	case Ack:
		return true
	}

	return false
}
