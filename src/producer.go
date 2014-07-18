package main

import (
	"bitbucket.org/gdamore/mangos/protocol/req"
	"bitbucket.org/gdamore/mangos/transport/all"
	"container/list"
	"errors"
	"github.com/ugorji/go/codec"
	"log"
	"net/url"
	"time"
)

// TODO: implement buffer size counter
// TODO: implement retry counter
// TODO: look into creating a conn pool

type Producer struct {
	mh            codec.MsgpackHandle
	buffer        *list.List
	backendUrl    string
	flushInterval int
}

func newProducer(u string, flushInt int) *Producer {
	producer := new(Producer)
	producer.buffer = list.New()
	producer.backendUrl = u
	producer.flushInterval = flushInt
	go func() {
		log.Println("Buffer flush started")
		for {
			time.Sleep(time.Duration(producer.flushInterval) * time.Second)
			producer.flushBuffer()
			log.Println("Buffer flushed ")
		}
	}()
	return producer
}

func (p *Producer) flushBuffer() {
	for i := p.buffer.Front(); i != nil; i = i.Next() {
		wat := i.Value.([]byte)
		err := p.sendMessage(&wat)
		if err != nil {
			p.buffer.PushBack(wat)
		}
	}
}

func (p *Producer) sendMessage(message *[]byte) error {
	requestSocket, err := req.NewSocket()
	if err != nil {
		return err
	}
	defer requestSocket.Close()
	all.AddTransports(requestSocket)

	if err = requestSocket.Dial(p.backendUrl); err != nil {
		log.Println(err)
		return err
	}

	if err = requestSocket.Send(*message); err != nil {
		log.Println(err)
		return err
	}

	var clientMsg []byte

	if clientMsg, err = requestSocket.Recv(); err != nil {
		log.Println(err)
		return err
	}

	if string(clientMsg) != "OK" {
		return errors.New("Response not OK, requeued")
	}

	return nil
}

func (p *Producer) Send(query url.Values) error {
	err, b := p.encode(query)
	if err != nil {
		return err
	}
	err = p.sendMessage(&b)
	if err != nil {
		p.buffer.PushBack(b)
	}
	return nil
}

func (p *Producer) encode(query url.Values) (error, []byte) {
	var b []byte
	enc := codec.NewEncoderBytes(&b, &p.mh)
	err := enc.Encode(query)
	if err != nil {
		return err, nil
	}
	return nil, b
}
