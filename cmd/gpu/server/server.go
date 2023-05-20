package gpu_server

import "golang.org/x/crypto/ssh"

type Server struct {
	ssh ssh.Client
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Connect() error {
	return nil
}
