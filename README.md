# simple-message-broker

* Implemented a simple message broker that supports **topic-based** and **queue-based** messaging patterns using Go. 

* Employed **RPCs** to facilitate for server and client communications with the message broker service.


## Implementation details


### Command details

* For commands supported by the client, please refer to client/commands.go
* For commands supported by the server, please refer to server/commands.go
* For commands supported by the broker, please refer to broker/commands.go

### Connection Configuration

* To check and edit the endpoints, please refer to configs/config.json


### Queue details

* For the queue implementation, please refer to queue/queue.go


### Connection and Communication protocol - 

* For code related to connection and communication protocol details, please refer to communication_protocol package.


## Usage

To start the broker process, please run the command -

```bash
    go run /cmd/broker.go
```

To start the server process, please run the command - 
```bash
    go run /cmd/server.go
```

To start the client process, please run the command -
```bash
    go run /cmd/client.go
```