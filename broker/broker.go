package broker

import (
	"bufio"
	"fmt"
	"log"
	"message_broker/communication_protocol"
	"message_broker/configuration"
	"message_broker/queue"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type Broker struct {
	sendQueue               queue.Queue
	recvQueue               queue.Queue
	isSyncMessageTransfered bool
	topics                  map[string]*Topic
}

type Topic struct {
	name        string
	subscribers []net.Conn
}

func printMessages(messages []queue.Message) {
	fmt.Print("[")
	for i := len(messages) - 1; i >= 0; i-- {
		fmt.Printf("%s", messages[i].Data)
		if i > 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println("]")
}

func printTopics(b map[string]*Topic) {
	if len(b) == 0 {
		fmt.Println("No topics")
	} else {
		for _, topic := range b {
			fmt.Printf("%s: %d subscribers\n", topic.name, len(topic.subscribers))
		}
	}
}

func clearGetQueue(b *Broker) {
	b.recvQueue.Clear()
	fmt.Println("send queue cleared")

}

func clearPutQueue(b *Broker) {
	b.sendQueue.Clear()
	fmt.Println("recv queue cleared")
}

func (b *Broker) serve() {
	rpc.Register(b)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", configuration.GetRPCAddress())
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go http.Serve(l, nil)
}

func userInfoMessages() {
	fmt.Println("<---------------------------------------------------------------------------------------------->")
	fmt.Println("Enter the command:")
	fmt.Println("display_queue_get | display_queue_put | display_topics | clear_queue_get | clear_queue_put | exit")
}

func commandProcessor(b *Broker) {
	userInfoMessages()
	var command, _ = extractCommand()

	switch command {
	case display_queue_get:
		printMessages(b.recvQueue.Messages)

	case display_queue_put:
		printMessages(b.sendQueue.Messages)

	case display_topics:
		printTopics(b.topics)

	case clear_queue_get:
		clearGetQueue(b)

	case clear_queue_put:
		clearPutQueue(b)

	case exit:
		os.Exit(0)

	default:
		fmt.Println("Invalid command. Please try again... ")
	}
}

func closeConn(conn net.Conn, b *Broker) {
	for topicName := range b.topics {
		if b.isConnSubscribedToTopic(conn, topicName) {
			b.removeConnFromTopicSubscribers(conn, topicName)
		}
	}
	err := conn.Close()
	if err != nil {
		log.Print("close error:", err)
	}
}

func (b *Broker) removeConnFromTopicSubscribers(conn net.Conn, topicName string) {
	if !b.isConnSubscribedToTopic(conn, topicName) {
		fmt.Println("Not subscribed to the topic.. Please try again.. ")
	}
	for i, c := range b.topics[topicName].subscribers {
		if c == conn {
			b.topics[topicName].subscribers = append(b.topics[topicName].subscribers[:i], b.topics[topicName].subscribers[i+1:]...)
			return
		}
	}
}

func (b *Broker) isConnSubscribedToTopic(conn net.Conn, topicName string) bool {
	for _, c := range b.topics[topicName].subscribers {
		if c == conn {
			return true
		}
	}
	return false
}

func extractCommand() (string, string) {
	bufferReader := bufio.NewReader(os.Stdin)
	var command string
	command, _ = bufferReader.ReadString('\n')
	command = strings.TrimSuffix(command, "\n")
	splitValues := strings.SplitN(command, " ", 2)
	if len(splitValues) == 1 {
		return splitValues[0], ""
	} else {
		return splitValues[0], splitValues[1]
	}
}

func parseClientRequest(clientRequest string) (string, string) {
	splitValues := strings.SplitN(strings.TrimSuffix(clientRequest, "\n"), " ", 2)
	if len(splitValues) == 1 {
		return splitValues[0], ""
	} else {
		return splitValues[0], splitValues[1]
	}
}

func serveClient(con net.Conn, b *Broker) {
	defer closeConn(con, b)
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadString('\n')
		if err != nil {
			continue
		}
		command, arguments := parseClientRequest(clientRequest)
		switch command {
		case "subscribe":
			b.topics[arguments].subscribers = append(b.topics[arguments].subscribers, con)
		case "unsubscribe":
			b.removeConnFromTopicSubscribers(con, arguments)
		}
	}
}

func acceptClient(subListener net.Listener, b *Broker) {
	for {
		con, err := subListener.Accept()
		if err != nil {
			log.Print("accept error:", err)
			continue
		}
		go serveClient(con, b)
	}
}

func (b *Broker) PutMessage(args *communication_protocol.PutMessageArgs, reply *communication_protocol.PutMessageReply) error {
	b.isSyncMessageTransfered = false
	err := b.sendQueue.Put(&queue.Message{Data: args.Message})
	if err != nil {
		if err.Error() == queue.OverflowErrorMsg {
			reply.IsBufferOverflow = true
			return nil
		}
		return err
	}
	if args.IsAsync {
		return nil
	}
	for !b.isSyncMessageTransfered {
		time.Sleep(time.Second)
	}
	reply.IsBufferOverflow = false
	return nil
}

func (b *Broker) PutBackMessage(args *communication_protocol.PutBackMessageArgs, reply *communication_protocol.PutBackMessageReply) error {
	err := b.recvQueue.Put(&queue.Message{Data: args.Message})
	if err != nil {
		if err.Error() == queue.OverflowErrorMsg {
			reply.IsBufferOverflow = true
			return nil
		}
		return err
	}
	reply.IsBufferOverflow = false
	return nil
}

func (b *Broker) GetMessage(args *communication_protocol.GetMessageArgs, reply *communication_protocol.GetMessageReply) error {
	msg, err := b.sendQueue.Pop()
	if err != nil {
		return err
	}
	if msg != nil {
		reply.Message = msg.Data
		b.isSyncMessageTransfered = true
	}
	return nil
}

func (b *Broker) GetBackMessage(args *communication_protocol.GetBackMessageArgs, reply *communication_protocol.GetBackMessageReply) error {
	msg, err := b.recvQueue.Pop()
	if err != nil {
		return err
	}
	if msg != nil {
		reply.Message = msg.Data
	}
	return nil
}

func (b *Broker) CreateTopic(args *communication_protocol.CreateTopicArgs, reply *communication_protocol.CreateTopicReply) error {
	if _, ok := b.topics[args.TopicName]; ok {
		return fmt.Errorf("topic already exists")
	}
	b.topics[args.TopicName] = &Topic{args.TopicName, []net.Conn{}}
	return nil
}

func (b *Broker) Publish(args *communication_protocol.PublishArgs, reply *communication_protocol.PublishReply) error {
	if _, ok := b.topics[args.TopicName]; !ok {
		return fmt.Errorf("topic doesn't exist")
	}
	for _, c := range b.topics[args.TopicName].subscribers {
		c.Write([]byte(args.Message + "\n"))
	}
	return nil
}

func Run() {
	broker := &Broker{queue.Queue{}, queue.Queue{}, false, make(map[string]*Topic)}
	broker.serve()

	subListener, err := net.Listen("tcp", configuration.GetPubSubAddress())
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go acceptClient(subListener, broker)

	for {
		commandProcessor(broker)
	}
}
