ASSETS=../assets/public_html/js

all: download minify install

download:
	wget https://github.com/yui/yuicompressor/releases/download/v2.4.8/yuicompressor-2.4.8.jar

minify:
	cat track.js | java -jar yuicompressor-2.4.8.jar --type js -o track.min.js

install:
	cp track.min.js $(ASSETS)
	
clean:
	rm -rf yuicompressor-2.4.8.jar $(ASSETS)/track.min.js track.min.js

