package server

import (
	"bufio"
	"fmt"
	"message_broker/communication_protocol"
	"os"
	"strings"
)

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

func putMessageHandler(message string, isAsync bool) bool {
	args := communication_protocol.PutMessageArgs{Message: message, IsAsync: isAsync}
	reply := communication_protocol.PutMessageReply{}
	if !communication_protocol.Call("Broker.PutMessage", &args, &reply) {
		os.Exit(0)
	}
	return reply.IsBufferOverflow
}

func CallForCreateTopic(topic string) communication_protocol.CreateTopicReply {
	args := communication_protocol.CreateTopicArgs{TopicName: topic}
	reply := communication_protocol.CreateTopicReply{}
	if !communication_protocol.Call("Broker.CreateTopic", &args, &reply) {
		os.Exit(0)
	}
	return reply
}

func CallForPublish(topic string, message string) communication_protocol.PublishReply {
	args := communication_protocol.PublishArgs{TopicName: topic, Message: message}
	reply := communication_protocol.PublishReply{}
	if !communication_protocol.Call("Broker.Publish", &args, &reply) {
		os.Exit(0)
	}
	return reply
}

func CallForGetBackMessage() communication_protocol.GetBackMessageReply {
	args := communication_protocol.GetBackMessageArgs{}
	reply := communication_protocol.GetBackMessageReply{}
	if !communication_protocol.Call("Broker.GetBackMessage", &args, &reply) {
		os.Exit(0)
	}
	return reply
}

func userInfoMessages() {
	fmt.Println("<---------------------------------------------------------------------------------------------->")
	fmt.Println("Enter the command:")
	fmt.Println("options: put_async <message>| put_sync <message>| get | create_topic <topic_name>| publish <topic> <message> | exit")
}

func commandProcessor() {
	userInfoMessages()
	var command, arguments = extractCommand()

	switch command {
	case get:
		reply := CallForGetBackMessage()
		message := reply.Message
		if message == "" {
			fmt.Println("No message to get back")
		} else {
			fmt.Println("Message: " + message)
		}

	case put_sync:
		status := putMessageHandler(arguments, false)
		if status {
			fmt.Println("Ran into error while sending the message")
		} else {
			fmt.Println("Message sent successfully.")
		}

	case put_async:
		status := putMessageHandler(arguments, true)
		if status {
			fmt.Println("Ran into error while sending the message")
		} else {
			fmt.Println("Message sent successfully.")
		}

	case publish:
		splitArguments := strings.SplitN(arguments, " ", 2)
		CallForPublish(splitArguments[0], splitArguments[1])
		fmt.Println("Message published successfully")

	case create_topic:
		CallForCreateTopic(arguments)
		fmt.Println("Topic created successfully")

	case exit:
		os.Exit(0)

	default:
		fmt.Println("Invalid command. Please try again... ")
	}

}

func Run() {
	for {
		commandProcessor()
	}
}
