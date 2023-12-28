package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func extractCommand() string {
	bufferReader := bufio.NewReader(os.Stdin)
	var command string
	command, _ = bufferReader.ReadString('\n')
	command = strings.TrimSuffix(command, "\n")
	return strings.Split(command, " ")[0]

}

func userInfoMessages() {
	fmt.Println("<---------------------------------------------------------------------------------------------->")
	fmt.Println("Enter the command:")
	fmt.Println("options: put_async <message>| put_sync <message>| get | create_topic <topic_name>| publish <topic> <message> | exit")
}

func commandProcessor() {
	userInfoMessages()
	var command = extractCommand()

	switch command {
	case get:
		break

	case put_sync:
		break

	case put_async:
		break

	case publish:
		break

	case create_topic:
		break

	case exit:
		break

	default:
		break
	}

}
