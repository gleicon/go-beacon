package main

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"net/url"
)

var (
	mh codec.MsgpackHandle
)

func send(query url.Values) error {
	err, b := encode(query)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
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
