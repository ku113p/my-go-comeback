package server

import (
	"fmt"
	"net/http"
)

type server struct {
	ip   [4]byte
	port uint16
}

func (s *server) Run() {
	mux := getRootMux()

	http.ListenAndServe(s.url(), mux)
}

func NewServer(ip [4]byte, port uint16) server {
	return server{ip, port}
}

func (s *server) url() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", s.ip[0], s.ip[1], s.ip[2], s.ip[3], s.port)
}
