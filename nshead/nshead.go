package nshead

import (
	"encoding/binary"
	"io"
	"math"
)

type Nshead struct {
	Id       uint16
	Version  uint16
	Logid    uint32
	Provider [16]byte
	MagicNum uint32
	Reserved uint32
	BodyLen  uint32
}

type NsheadPacket struct {
	Header Nshead
	Body   []byte
}

const (
	MAGIC_NUM = 0xfb709394
)

type NsheadError struct {
	s string
}

func (err *NsheadError) Error() string {
	return err.s
}

var (
	ErrMalformedHeader = &NsheadError{"malformed header"}
	UnexpectIOError    = &NsheadError{"io error"}
	ErrBodyTooLarge    = &NsheadError{"body too large"}
)

func ReadNsheadPacket(r io.Reader) (packet *NsheadPacket, err *NsheadError) {
	packet = &NsheadPacket{}
	ierr := binary.Read(r, binary.BigEndian, &packet.Header)
	if ierr != nil {
		err = UnexpectIOError
		return
	}
	if packet.Header.MagicNum != MAGIC_NUM {
		err = ErrMalformedHeader
		return
	}

	//zero body is not an error
	if packet.Header.BodyLen == 0 {
		return
	}

	packet.Body = make([]byte, packet.Header.BodyLen)

	_, ierr = io.ReadFull(r, packet.Body)
	if ierr != nil {
		err = UnexpectIOError
		return
	}
	return
}

func (packet *NsheadPacket) Write(w io.Writer) (err *NsheadError) {
	if len(packet.Body) > math.MaxInt32 {
		err = ErrBodyTooLarge
	}
	packet.Header.BodyLen = uint32(len(packet.Body))
	packet.Header.MagicNum = MAGIC_NUM

	ierr := binary.Write(w, binary.BigEndian, packet.Header)
	if ierr != nil {
		err = UnexpectIOError
		return
	}

	_, ierr = w.Write(packet.Body)
	if ierr != nil {
		err = UnexpectIOError
		return
	}

	return
}
