# go-beacon

Generic beacon for js trackers. It is developed with boomerang.js (https://github.com/lognormal/boomerang/) in mind, 
but can be used by any other javascript tracker. It answers to a request to an URI with a transparent GIF.

The default behaviour is to log a JSON of the query string to a file, but it can be plugged into any backend.

I've implemented the boomerang async method. Just include /js/tracker.js into your web page like this:

       <script src="/js/tracker.js"></script>

Look below for the boomerang build section for more information.

## Download and Build

    $ git clone https://github.com/gleicon/go-build
    $ cd go-build
    $ make

## Build SSL

If you want to use SSL, just set it right into the config file. In case you want a self-signed SSL:
    
    $ cd SSL
    $ make

Answer the questions and you are ready to go.

## Build Boomerang

The best way is to serve boomerang from a CDN, but in case you need a specific setup and want go-beacon to serve boomerang for you:

    $ cd boomerang_build
    $ make

Edit Makefile to pick a different minifier or configure the plugins that you want. 
There's a plugin called async_loading.js that implements async loading as per boomerang request, but that can be removed in the Makefile.

## Install

Edit the config file and run the server, check the beacon_uri parameter and execute the server:

	vi go-beacon.conf
	./go-beacon

Install, uninstall. Edit Makefile and set PREFIX to the target directory:

	sudo make install
	sudo make uninstall

Allow non-root process to listen on low ports:

	/sbin/setcap 'cap_net_bind_service=+ep' /opt/go-beacon/server


Gleicon 2014 - MIT License.
