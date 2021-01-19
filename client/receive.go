package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"net"
)

type Receiver struct {
	conn net.Conn
}

func NewReceiver(conn net.Conn) *Receiver {
	return &Receiver{
		conn: conn,
	}
}

func (r *Receiver) ReadAndPlay() {
	tmp := make([]byte, 256*4)
	data := make([]int32, 256)

	//Open Mic Stream
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(data), data)
	chk(err)

	chk(stream.Start())
	for {
		_, err := r.conn.Read(tmp)
		chk(err)

		//Buffer to write bytes to int32
		buf := bytes.NewReader(tmp)

		err = binary.Read(buf, binary.BigEndian, &data)
		chk(err)

		//Write the resulting []int32 to the output Stream
		err = stream.Write()
		if err != nil {
			if errors.Is(err, portaudio.OutputUnderflowed) {
				fmt.Println("Underflow")
			}
		}
	}
}