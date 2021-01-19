package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"net"
	"os"
	"os/signal"
)

type Transmitter struct {
	conn *net.Conn
	client *Client
	stream *portaudio.Stream
	term chan bool
}

func NewTransmitter(conn *net.Conn) *Transmitter {
	return &Transmitter{
		conn: conn,
	}
}

func (t *Transmitter) RecordAndStream() {
	//Signal to catch Ctrl-C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	//Input buffer
	in := make([]int32, 256)

	//Open Mic Stream
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	t.stream = stream

	//Buffer to write int32 to bytes
	buf := new(bytes.Buffer)
	chk(stream.Start())
	for {
		err := stream.Read()
		if err != nil {
			if !errors.Is(err, portaudio.InputOverflowed) {
				panic(err)
			} else {
				fmt.Println("overflow")
			}
		}


		//Reset buffer and write int32 to output byte to buffer
		buf.Reset()
		err = binary.Write(buf, binary.BigEndian, in)
		chk(err)

		//Write the resulting bytes to socket
		_, err = (*t.conn).Write(buf.Bytes())
		//fmt.Println(buf.Bytes())
		chk(err)
		select {
		case <-sig:
			fmt.Println("Finalising from keyboard interrupt")
			t.finalise()
			return
		case <-t.term:
			fmt.Println("Finalising from stop method")
			t.finalise()
			return
		default:
		}
		//time.Sleep(time.Second*2)
	}
}

func (t *Transmitter) Stop() {
	t.term<-true
}

func (t *Transmitter) finalise() {
	chk(t.stream.Stop())
	chk(t.stream.Close())
	t.client.Stop()
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
