package model

type MessageType int

const (
	MessageTypeRequest = iota + 1
	MessageTypeResponse
)

type ThreadMessage struct {
	Type MessageType
	Text string
}

type Thread struct {
	Messages []*ThreadMessage
}
