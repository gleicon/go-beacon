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
	var b []byte
	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(query)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
