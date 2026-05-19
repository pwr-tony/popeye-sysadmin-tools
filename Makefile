.PHONY: build dev clean install

build:
	wails build -tags webkit2_41

dev:
	wails dev -tags webkit2_41

clean:
	rm -rf build/bin/*

install: build
	cp build/bin/popeye /usr/local/bin/

run: build
	./build/bin/popeye
