package main

import (
	"flag"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"voip/client"
	"voip/server"
)

func main() {
	runServerPtr := flag.Bool("server", false, "Run server")
	runClientPtr := flag.Bool("client", false, "Run client")

	flag.Parse()

	if *runClientPtr == *runServerPtr {
		flag.Usage()
		//flag.PrintDefaults()
		os.Exit(1)
	}

	if *runClientPtr {
		err := portaudio.Initialize()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer portaudio.Terminate()
		devInfo, err := portaudio.DefaultInputDevice()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(devInfo.Name)
		}
		c := client.NewClient(8080)
		c.Connect()
		c.Start()
	} else if *runServerPtr {
		s := server.NewServer()
		s.Listen()
	}
}