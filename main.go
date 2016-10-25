package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	sock := "/var/run/docker.sock"
	return net.Dial("unix", sock)
}

func main() {

	path := "/tmp/just-the-events.sock"
	os.Remove(path)
	l, err := net.ListenUnix("unix", &net.UnixAddr{path, "unix"})
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(path)

	director := func(req *http.Request) {
		// to make Dial work with unix URL,
		// scheme and host have to be faked
		fmt.Println(req)
		req.URL.Scheme = "http"
		req.URL.Host = "socket"
	}
	p := &httputil.ReverseProxy{
		Director:      director,
		FlushInterval: time.Second,
	}
	p.Transport = &http.Transport{
		Dial: fakeDial,
	}

	log.Fatal(http.Serve(l, p))
}
