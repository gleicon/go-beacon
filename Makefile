# Copyright 2014 go-beacon authors.  All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

TARGET=$(DESTDIR)/opt/go-beacon
include src/Makefile.defs
include backend/Makefile.defs

all: server consumer

deps:
	make -C src deps
	make -C backend deps

server:
	make -C src
	@cp src/$(NAME) .

consumer:
	make -C backend
	@cp backend/$(CONSUMER_NAME) .

clean:
	make -C src clean
	make -C backend clean
	@rm -f $(NAME)
	@rm -f $(CONSUMER_NAME)

install: server
	install -m 750 -d $(TARGET)
	install -m 750 $(NAME) $(TARGET)
	install -m 640 go-beacon.conf $(TARGET)
	install -m 750 -d $(TARGET)/ssl
	install -m 640 ssl/Makefile $(TARGET)/ssl
	install -m 750 -d $(TARGET)/assets
	rsync -rupE assets $(TARGET)
	find $(TARGET)/assets -type f -exec chmod 640 {} \;
	find $(TARGET)/assets -type d -exec chmod 750 {} \;
	#chown -R www-data: $(TARGET)

uninstall:
	rm -rf $(TARGET)

.PHONY: server
