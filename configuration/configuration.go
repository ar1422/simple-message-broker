package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Address struct {
	Host string
	Port int
}

type Configuration struct {
	RPCService    Address
	PubSubService Address
}

func readConfig() Configuration {
	f, err := os.Open("configs/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)
	var config Configuration
	json.Unmarshal(byteValue, &config)
	return config
}

func GetRPCHost() string {
	return readConfig().RPCService.Host
}

func GetRPCPort() int {
	return readConfig().RPCService.Port
}

func GetRPCAddress() string {

	return fmt.Sprintf("%s:%d", GetRPCHost(), GetRPCPort())
}

func GetPubSubHost() string {
	return readConfig().PubSubService.Host
}

func GetPubSubPort() int {
	return readConfig().PubSubService.Port
}

func GetPubSubAddress() string {

	return fmt.Sprintf("%s:%d", GetPubSubHost(), GetPubSubPort())
}
