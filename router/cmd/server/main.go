package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go TCPReceiveMessage(c)

	wg.Wait()
}

func TCPReceiveMessage(c net.Conn) {
	defer func(c net.Conn) {
		fmt.Println("Kill connection")
		c.Close()
	}(c)

	for {
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
	}
}

type Teste struct {
	Id       string
	Mensagem string
}

func KafkaConn() {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	message, _ := json.Marshal(Teste{Id: "teste", Mensagem: "Mensagem teste"})

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: message},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
