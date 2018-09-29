package gltun

type Packet interface {
	GetType() int8
	GetPayload() []byte
	Encode() []byte
}

type packet struct {
	t int8
	b []byte
}

func NewPacket(pkgType int8, payload []byte) (Packet) {
	return &packet{t:pkgType, b:payload}
}

func (p *packet) GetType() int8 { return p.t }
func (p *packet) GetPayload() []byte { return p.b }
func (p *packet) Encode() []byte {
	header := make([]byte, 4)
	header[0] = byte(p.t)
	size := uint32(len(p.b))

	header[1] = byte(size >> 16)
	header[2] = byte(size >> 8)
	header[3] = byte(size >> 0)

	return append(header, p.b...)
}