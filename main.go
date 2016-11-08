package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
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

	patterns := regexp.MustCompile(`/v1\.\d+/(events|version|(containers(/[0-9a-f]+)?/json))`)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" && patterns.MatchString(req.URL.Path) {
			p.ServeHTTP(w, req)
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Forbidden: ", req)
		}
	})

	log.Fatal(http.Serve(l, mux))
}
