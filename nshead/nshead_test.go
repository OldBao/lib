package nshead

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	c = []byte{
		0x00, 0x01, //id
		0x00, 0x02, //version
		0x00, 0x00, 0x00, 0x03, //Logid
		'H', 'E', 'L', 'L', 'O', 'T', 'H', 'E', 'W', 'O', 'R', 'L', 'D', '!', '!', 0x00,
		0xfb, 0x70, 0x93, 0x94, //magicnum
		0x00, 0x00, 0x00, 0x04, //reserverd
		0x00, 0x00, 0x00, 0x05, //bodylen
		'Z', 'H', 'A', 'N', 'G',
	}
)

func TestReadHeader(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader(c)
	packet, err := ReadNsheadPacket(r)
	if err != nil {
		t.Errorf("new nshead reader error %s", err.Error)
		return
	}

	assert.Equal(packet.Header.Id, uint16(1))
	assert.Equal(packet.Header.Version, uint16(2))
	assert.Equal(packet.Header.Logid, uint32(3))
	assert.Equal(string(packet.Header.Provider[:]), "HELLOTHEWORLD!!\x00")
	assert.Equal(packet.Header.Reserved, uint32(4))
	assert.Equal(packet.Header.BodyLen, uint32(5))

	assert.Equal(string(packet.Body), "ZHANG")
}

func TestWrite(t *testing.T) {
	assert := assert.New(t)

	packet := NsheadPacket{
		Header: Nshead{
			Id:       1,
			Version:  2,
			Logid:    3,
			Reserved: 4,
		},
	}
	copy(packet.Header.Provider[:], []byte("HELLOTHEWORLD!!\x00"))
	packet.Body = []byte("ZHANG")

	var raw []byte
	buffer := bytes.NewBuffer(raw)
	buffer.Grow(len(c))
	assert.Nil(packet.Write(buffer))

	assert.Equal(buffer.Bytes(), c)
}
