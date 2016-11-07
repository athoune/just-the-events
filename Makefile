build:
	go build .

linux:
	docker run -it --rm -v `pwd`:/go golang go build -o just-the-events_linux .
