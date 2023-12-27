package broker

import (
	"bufio"
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

func commandProcessor() {
	var command = extractCommand()

	switch command {
	case display_queue_get:
		break

	case display_queue_put:
		break

	case display_topics:
		break

	case clear_queue_get:
		break

	case clear_queue_put:
		break

	case exit:
		break

	default:
		break
	}
}
