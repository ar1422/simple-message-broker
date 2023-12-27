package client

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
	case get:
		break

	case put:
		break

	case subscribe:
		break

	case unsubscribe:
		break

	case exit:
		break

	default:
		break
	}

}
