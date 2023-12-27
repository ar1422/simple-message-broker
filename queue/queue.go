package queue

import "fmt"

func (q *Queue) Full() bool {
	return len(q.Messages) == MaxBufferSize

}

func (q *Queue) Empty() bool {
	return len(q.Messages) == 0
}

func (q *Queue) Put(msg *Message) error {
	if q.Full() {
		return fmt.Errorf(OverflowErrorMsg)
	} else {
		q.Messages = append(q.Messages, *msg)
		return nil
	}
}

func (q *Queue) Get() (*Message, error) {
	if q.Empty() {
		return nil, nil
	} else {
		return &q.Messages[0], nil
	}
}

func (q *Queue) Pop() (*Message, error) {
	if q.Empty() {
		return nil, nil
	}
	msg := q.Messages[0]
	q.Messages = q.Messages[1:]
	return &msg, nil
}

func (q *Queue) Clear() error {
	q.Messages = []Message{}
	return nil
}
