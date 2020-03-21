package draft

import (
	"net"
	"net/http"
)

type Server struct {
	http.Server
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := &Server{http.Server{Addr: addr, Handler: handler}}
	return server.ListenAndServe()
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *Server) Serve(ln net.Listener) error {
	return nil
}
