package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn net.Conn
	port int
}

func NewClient(port int) *Client {
	return &Client{port: port}
}

func (c *Client) Connect() bool {
	conn, err := net.Dial("tcp", ":8091")
	if err != nil {
		// handle error
		fmt.Println(err)
		return false
	}
	fmt.Println("connected!")
	c.conn = conn
	return true
}

func (c *Client) Start() {
	transmitter := NewTransmitter(&c.conn)
	receiver := NewReceiver(c.conn)

	go transmitter.RecordAndStream()
	go receiver.ReadAndPlay()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter Text: ")
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()
		switch text {
		case "q":
			fmt.Println("We gotta end")
			c.Stop()
			transmitter.Stop()
			return
		default:
		}
	}

}

func (c *Client) Stop() {
	err := c.conn.Close()

	if err != nil {
		panic(err)
	}
	fmt.Println("Exiting...")
	os.Exit(1)
	fmt.Println("Done.")
}

func (c *Client) Send(s string) {
	fmt.Print(c.conn)
	_, err := c.conn.Write([]byte(s))
	if err != nil {
		fmt.Println("Failed to send data")
	}
}
