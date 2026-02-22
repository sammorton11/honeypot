package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Attempt struct {
	Address string `json:"address"`
	Network string `json:"network"`
	Message string `json:"message"`
}

var DOCKER_URL = fmt.Sprintf("%s%s", os.Getenv("API_URL"), "/api/attempt")
const DEV_URL = "http://localhost:8080/api/attempt"

func main() {
	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("honeypot listening on :2222")
	fmt.Println("try: nc localhost 2222")

	// Forever loop that creates tcp connections with clients
	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
		}

		clientAddress := c.RemoteAddr().String()
		clientNetworkType := c.RemoteAddr().Network()

		// fake SSH banner to bait attackers
		n, err := c.Write([]byte("SSH-2.0-OpenSSH_7.4\r\n"))
		if err != nil {
			log.Println(fmt.Errorf("Error writing banner to connection: %w\n", err))
		}
		fmt.Printf("wrote %d bytes\n", n)

		// launch connections in a go routine -- so we can have more than one connection
		go handleConnection(c, clientAddress, clientNetworkType)
	}
}

func tarpit(c net.Conn) {
	for i := range 10 {
		_, err := c.Write([]byte("Loading...\r\n"))
		if err != nil {
			log.Println("Tarpit Error:", err)
			break
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
}

func sanitizeMessage(message []byte) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return '_' // replace with underscore
		}
		return r
	}, string(message))
}

func handleConnection(c net.Conn, userAddr, userNetwork string) {
	// buffer to read whatever they send (login attempts, etc)
	buf := make([]byte, 1024)
	c.SetDeadline(time.Now().Add(5 * time.Second)) // set a dead line for reads/writes

	for {
		// read data from the connection into the buffer
		n, err := c.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection: ", err)
			break
		}

		// Create message slice from buffer
		message := buf[:n]

		// Remove any '\n' or '\r' to prevent log injections
		safeMessage := sanitizeMessage(message)

		// TODO: convertAddressToCountry -- or could do this on the frontend

		// Package data into an object
		attempt := Attempt{
			Address: userAddr,
			Network: userNetwork,
			Message: safeMessage,
		}

		// Send to proxy server to store in db
		err = sendAttempt(attempt, DOCKER_URL)
		if err != nil {
			log.Println(fmt.Errorf("Error closing response body %w\n", err))
			return
		}

		// Trap the attacker in a tarpit
		tarpit(c)
	}

	// Close the TCP connection
	err := c.Close()
	if err != nil {
		fmt.Println("Error closing connection", err)
	}
}

func sendAttempt(attempt Attempt, postUrl string) error {
	// Serialize
	jsonData, err := json.Marshal(attempt)
	if err != nil {
		fmt.Println("Error marshaling attempt to rest server", err)
	}

	// Send over the wire to proxy server to store in db
	req, err := http.Post(postUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	return req.Body.Close()
}
