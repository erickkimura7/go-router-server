package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]

	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	var connMap = &sync.Map{}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		connMap.Store(id, c)

		go handleConn(id, c, connMap)
	}
}

func handleConn(id string, client net.Conn, connMap *sync.Map) {
	defer func(client net.Conn, connMap *sync.Map) {
		fmt.Println("Kill connection")

		client.Close()
		connMap.Delete(id)
	}(client, connMap)

	fmt.Printf("Iniciando nova conexao [%s]\n", id)

	sendingMessage(id, client)
}

func sendingMessage(id string, client net.Conn) {
	retry := 0
	for {
		fmt.Printf("Enviando mensagem [%s]\n", id)
		_, err := client.Write([]byte(fmt.Sprintf("id: [%s] - Mensagem\n", id)))
		if err != nil {
			fmt.Println(err)
			retry++
			if retry > 4 {
				return
			}
		}

		time.Sleep(time.Duration(rand.Intn(30)+2) * time.Second)
	}
}
