package server

import "net/http"

type Server struct {
	http.Server
	Port uint
}

func (*Server) Listen() err {

}

func New() *Server {
	return &Server{}
}
