package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

type jsonMessage struct {
	Method string `json:"method"`
}

func main() {
	var localAddr string
	var upstreamAddr string
	var delay int

	flag.StringVar(&localAddr, "local", "", "The local address to listen on (required)")
	flag.StringVar(&upstreamAddr, "upstream", "", "The upstream address to proxy to (required)")
	flag.IntVar(&delay, "delay", 0, "Delay for eth_submitWork method in milliseconds")

	flag.Parse()

	if localAddr == "" || upstreamAddr == "" {
		flag.PrintDefaults()
		return
	}

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Println("Error starting listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on", localAddr)
	fmt.Println("Proxying to", upstreamAddr)
	fmt.Println("Delay for eth_submitWork method:", delay, "ms")

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(client, upstreamAddr, delay)
	}
}

func handleConnection(client net.Conn, upstreamAddr string, delay int) {
	target, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		fmt.Println("Error dialing upstream:", err)
		return
	}

	// Create a goroutine to proxy data from the client to the target
	go func() {
		reader := bufio.NewReader(client)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("[Client]: %v\n", err)
				break
			}
			var jsonLine jsonMessage
			if json.Unmarshal([]byte(line), &jsonLine) == nil && (jsonLine.Method == "eth_submitWork" || jsonLine.Method == "mining.submit") {
				go func(line string) {
					time.Sleep(time.Duration(delay) * time.Millisecond)
					fmt.Printf("[Client -> Target, delayed]: %s", line)
					target.Write([]byte(line))
				}(line)
			} else {
				fmt.Printf("[Client -> Target]: %s", line)
				target.Write([]byte(line))
			}
		}
		client.Close()
		target.Close()
	}()

	// Create a goroutine to proxy data from the target to the client
	go func() {
		reader := bufio.NewReader(target)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("[Target]: %v\n", err)
				break
			}
			fmt.Printf("[Target -> Client]: %s", line)
			client.Write([]byte(line))
		}
		target.Close()
		client.Close()
	}()
}
