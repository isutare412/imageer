package valkeystream

type Message struct {
	EntryID string
	Data    []byte
	Ack     func() error
}
