package queue

type Message struct {
	Data string
}

type Queue struct {
	Messages []Message
}
