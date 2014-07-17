package main

import (
	"bitbucket.org/gdamore/mangos"
	"bitbucket.org/gdamore/mangos/protocol/rep"
	"bitbucket.org/gdamore/mangos/transport/all"
	"fmt"
	"github.com/ugorji/go/codec"
	"time"
)

var (
	mh codec.MsgpackHandle
	b  []byte
)

func decode(buf []byte) (error, interface{}) {

	doc := map[string]interface{}(nil)
	dec := codec.NewDecoderBytes(buf, &mh)
	err := dec.Decode(&doc)

	if err != nil {
		return err, nil
	}
	return nil, doc
}

func main() {
	url := "tcp://127.0.0.1:8000"

	responseServerReady := make(chan struct{})
	responseServer, err := rep.NewSocket()
	defer responseServer.Close()

	all.AddTransports(responseServer)
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	fmt.Println("Consumer ready")

	go func() {
		var err error
		var serverMsg *mangos.Message

		if err = responseServer.Listen(url); err != nil {
			fmt.Printf("\nServer listen failed: %v", err)
			return
		}

		close(responseServerReady)

		for {
			if serverMsg, err = responseServer.RecvMsg(); err != nil {
				fmt.Printf("\nServer receive failed: %v", err)
			}

			err, d := decode(serverMsg.Body)
			fmt.Println("serverMsg: ", d)

			serverMsg.Body = []byte("OK")

			if err = responseServer.SendMsg(serverMsg); err != nil {
				fmt.Printf("\nServer send failed: %v", err)
				return
			}
		}
		fmt.Println("Listening")
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
