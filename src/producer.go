package main

import (
	"bitbucket.org/gdamore/mangos/protocol/req"
	"bitbucket.org/gdamore/mangos/transport/all"
	"container/list"
	"errors"
	"github.com/ugorji/go/codec"
	"net/url"
	"time"
)

// TODO: fetch url from configuration file
// TODO: implement buffer size counter
// TODO: push flush time to config
// TODO: move it to struct
// TODO: look into creating a conn pool

var (
	mh            codec.MsgpackHandle
	buffer        *list.List
	backendUrl    string
	flushInterval int
)

func initProducer(u string, flushInt int) {
	buffer = list.New()
	backendUrl = u
	flushInterval = flushInt
	go func() {
		time.Sleep(time.Duration(flushInterval) * time.Second)
		flushBuffer()
	}()
}

func flushBuffer() {
	for i := buffer.Front(); i != nil; i = i.Next() {
		wat := i.Value.([]byte)
		err := sendMessage(&wat)
		if err != nil {
			// TODO: retry counter
			buffer.PushBack(wat)
		}
	}
}

func sendMessage(message *[]byte) error {
	requestSocket, err := req.NewSocket()
	if err != nil {
		return err
	}
	defer requestSocket.Close()
	all.AddTransports(requestSocket)
	if err = requestSocket.Dial(backendUrl); err != nil {
		return err
	}

	if err = requestSocket.Send(*message); err != nil {
		return err
	}

	var clientMsg []byte

	if clientMsg, err = requestSocket.Recv(); err != nil {
		return err
	}

	if string(clientMsg) != "OK" {
		return errors.New("Response not OK, requeued")
	}

	return nil
}

func send(query url.Values) error {
	err, b := encode(query)
	if err != nil {
		return err
	}
	err = sendMessage(&b)
	if err != nil {
		buffer.PushBack(b)
	}
	return nil
}

func encode(query url.Values) (error, []byte) {
	var b []byte
	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(query)
	if err != nil {
		return err, nil
	}
	return nil, b
}
