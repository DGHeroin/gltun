package gltun

import (
	"net"
)

type Server interface {
	ListenAndServe(Addr string) error
	HandleDataFunc(func(Session, []byte))
	HandleConnectedFunc(func(Session))
	HandleClosedFunc(func(Session))
}

type server struct {
	listener net.Listener
	dataFunc func(Session, []byte)
	connectedFunc func(Session)
	closedFunc func(Session)
}

func NewServer() (Server) {
	return &server{}
}

// Listen and serve
//
func (s *server) ListenAndServe(Addr string) (error) {
	if l, err := net.Listen("tcp", Addr); err != nil {
		return err
	} else {
		s.listener = l
		s.serve()
	}
	return nil
}

//
//
func (s *server) serve() {
	for {
		if conn, err := s.listener.Accept(); err != nil {

		} else {
			go s.handleSession(conn)
		}
	}
}

// handle
//
func (s *server) handleSession(conn net.Conn)  {
	session := NewSession(conn)
	buf     := make([]byte, 4096)
	defer func() {
		if s.closedFunc != nil {
			s.closedFunc(session)
		}
		session.Close()
	}()
	if s.connectedFunc != nil {
		s.connectedFunc(session)
	}
	for {
		if n, err := conn.Read(buf); err != nil {
			break
		} else {
			data := buf[:n]
			if pkg := session.read(data); pkg != nil {
				if s.dataFunc != nil {
					s.dataFunc(session, pkg.GetPayload())
				}
			}
		}
	}
}

// on session data
//
func (s *server) HandleDataFunc(cb func(Session, []byte)) { s.dataFunc = cb }

// on new session connected
//
func (s *server) HandleConnectedFunc(cb func(Session)) { s.connectedFunc = cb}

// on session closed
//
func (s *server) HandleClosedFunc(cb func(Session)) { s.closedFunc = cb }