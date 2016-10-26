package main

import (
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
	os.Remove(path) //Who cares if the file doesn't exist?
	l, err := net.ListenUnix("unix", &net.UnixAddr{path, "unix"})
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(path)

	director := func(req *http.Request) {
		// to make Dial work with unix URL,
		// scheme and host have to be faked
		req.URL.Scheme = "http"
		req.URL.Host = "socket"
	}
	p := &httputil.ReverseProxy{
		Director:      director,
		FlushInterval: 250 * time.Millisecond,
	}
	p.Transport = &http.Transport{
		Dial: fakeDial,
	}

	mux := http.NewServeMux()
	mux.Handle("/v1.24/events", p)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	log.Fatal(http.Serve(l, mux))
}
