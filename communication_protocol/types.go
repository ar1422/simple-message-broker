package communication_protocol

type PutMessageArgs struct {
	Message string
	IsAsync bool
}

type PutMessageReply struct {
	IsBufferOverflow bool
}

type PutBackMessageArgs struct {
	Message string
}

type PutBackMessageReply struct {
	IsBufferOverflow bool
}

type GetMessageArgs struct {
}

type GetMessageReply struct {
	Message string
}

type GetBackMessageArgs struct {
}

type GetBackMessageReply struct {
	Message string
}

type CreateTopicArgs struct {
	TopicName string
}

type CreateTopicReply struct {
}

type PublishArgs struct {
	TopicName string
	Message   string
}

type PublishReply struct {
}
