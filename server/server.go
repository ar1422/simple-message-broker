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
	return splitValues[0], splitValues[1]
}

func putMessageHandler(message string, isAsync bool) bool {
	args := communication_protocol.PutMessageArgs{Message: message, IsAsync: isAsync}
	reply := communication_protocol.PutMessageReply{}
	if !communication_protocol.Call("Broker.PutMessage", &args, &reply) {
		os.Exit(0)
	}
	return reply.IsBufferOverflow
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
		break

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
		break

	case create_topic:
		break

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
