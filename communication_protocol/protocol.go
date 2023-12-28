package communication_protocol

import (
	"fmt"
	"log"
	"message_broker/configuration"
	"net/rpc"
)

func establishConnection() *rpc.Client {
	address := configuration.GetRPCAddress()
	c, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("Error while establishing connection:", err)
	}
	return c
}

func Call(rpcname string, args interface{}, reply interface{}) bool {

	client := establishConnection()
	defer client.Close()

	err := client.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
