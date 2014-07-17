# Copyright 2014 go-beacon authors.  All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

include Makefile.defs

all: deps server

deps:
	go get -v

server:
	go build -v -o $(CONSUMER_NAME) -ldflags "-X main.VERSION $(VERSION)"

clean:
	rm -f $(NAME)

.PHONY: server
