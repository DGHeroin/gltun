package gltun

import (
	"net"
)

type Session interface {
	read(data []byte) (Packet)

	Send([]byte) (error)
	Close()
}

type session struct {
	conn     net.Conn
	buffer   []byte

	readBody bool
	pkgType  int8
	pkgSize  uint32
}

func NewSession(conn net.Conn) (Session) {
	return &session{conn:conn}
}

func (s *session) read(data[]byte) (Packet) {
	if data != nil {
		s.buffer = append(s.buffer, data...)
	}
	// parse double time
	if p := s.parse(); p != nil {
		return p
	}
	if p := s.parse(); p != nil {
		return p
	}
	return nil
}

func (s *session) parse() (Packet) {
	if s.readBody == false { // read header
		if len(s.buffer) < 4 { return nil}
		s.pkgType = int8(s.buffer[0])
		s.pkgSize = uint32(s.buffer[1]) << 16 | uint32(s.buffer[2]) << 8 | uint32(s.buffer[3])

		s.buffer = s.buffer[4:]
		s.readBody = true
	} else { // read body
		if uint32(len(s.buffer)) < s.pkgSize { return nil }

		t := s.pkgType
		b := s.buffer[:s.pkgSize]

		s.buffer   = s.buffer[s.pkgSize:]
		s.pkgType  = 0
		s.pkgSize  = 0
		s.readBody = false

		return NewPacket(t, b)
	}
	return nil
}
func (s *session) Send(data []byte) (error) {
	bin := NewPacket(2, data).Encode()
	if s.conn != nil {
		if _, err := s.conn.Write(bin); err != nil {
			return err
		}
	}
	return nil
}

func (s *session) Close()  {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}