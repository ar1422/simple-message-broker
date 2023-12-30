package client

import (
	"bufio"
	"fmt"
	"log"
	"message_broker/communication_protocol"
	"message_broker/configuration"
	"net"
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

func userInfoMessages() {
	fmt.Println("<---------------------------------------------------------------------------------------------->")
	fmt.Println("Enter the command:")
	fmt.Println("options: get | put <message> | subscribe <topic> | unsubscribe <topic> | exit")
}

func subscriptionHandler(connection net.Conn, topic string) {
	subscriptionMsg := "subscribe " + topic + "\n"
	connection.Write([]byte(subscriptionMsg))
	fmt.Println("Subscribed to the topic - ", topic)
}

func unsubscriptionHandler(connection net.Conn, topic string) {
	unsubscriptionMsg := "unsubscribe " + topic + "\n"
	connection.Write([]byte(unsubscriptionMsg))
	fmt.Println("Unsubscribed from the topic - ", topic)
}

func getMessageHandler() string {
	arguments := communication_protocol.GetMessageArgs{}
	reply := communication_protocol.GetMessageReply{}
	if !communication_protocol.Call("Broker.GetMessage", &arguments, &reply) {
		os.Exit(0)
	}
	if reply.Message == "" {
		return "No messages are available"
	}
	return reply.Message
}

func putMessageHandler(message string) bool {
	args := communication_protocol.PutBackMessageArgs{Message: message}
	reply := communication_protocol.PutBackMessageReply{}
	if !communication_protocol.Call("Broker.PutBackMessage", &args, &reply) {
		os.Exit(0)
	}
	return reply.IsBufferOverflow

}

func listenToUpdates(connection net.Conn) {
	reader := bufio.NewReader(connection)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")

		fmt.Println("<---------------------------------------------------------------------------------------------->")
		fmt.Println("Message from the topic: ", message)
		fmt.Println("<---------------------------------------------------------------------------------------------->")
	}
}

func connectToSubscribeService() net.Conn {
	con, err := net.Dial("tcp", configuration.GetRPCAddress())
	if err != nil {
		log.Fatalln(err)
	}
	return con
}

func commandProcessor(connection net.Conn) {
	userInfoMessages()
	var command, arguments = extractCommand()

	switch command {
	case get:
		fmt.Println(getMessageHandler())

	case put:
		status := putMessageHandler(arguments)
		if status {
			fmt.Println("Ran into error while sending the message")
		} else {
			fmt.Println("Message sent successfully.")
		}

	case subscribe:
		subscriptionHandler(connection, arguments)
		go listenToUpdates(connection)
	case unsubscribe:
		unsubscriptionHandler(connection, arguments)

	case exit:
		os.Exit(0)

	default:
		fmt.Println("Invalid command. Please try again... ")
	}
}

func Run() {
	con := connectToSubscribeService()
	defer con.Close()
	for {
		commandProcessor(con)
	}

}
