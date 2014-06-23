ASSETS=../assets/public_html/js

all: download unzip build install

download:
	wget https://github.com/lognormal/boomerang/archive/master.zip
	wget https://github.com/yui/yuicompressor/releases/download/v2.4.8/yuicompressor-2.4.8.jar

unzip:
	unzip master.zip

build:
	cp async_loading.js boomerang-master
	(cd boomerang-master && make PLUGINS="plugins/rt.js plugins/bw.js plugins/dns.js plugins/navtiming.js plugins/restiming.js async_loading.js" MINIFIER="java -jar ../yuicompressor-2.4.8.jar --type js")

install:
	cp `(ls -Art boomerang-master/boomerang-*.js | grep -v debug | tail -n 1)` $(ASSETS)
	cp `(ls -Art boomerang-master/boomerang-*.js | grep -v debug | tail -n 1)` $(ASSETS)/boomerang.min.js
	
clean:
	rm -rf boomerang-master master.zip yuicompressor-2.4.8.jar $(ASSETS)/boomerang-*.js

