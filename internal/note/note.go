package note

import (
	"encoding/binary"
	"errors"
	"time"
)

type Note struct {
	Id   uint64
	Time time.Time
	Len  uint32
	Text string
}

func NewNote(text string) *Note {
	return &Note{Time: time.Now(), Len: uint32(len(text)), Text: text}
}

func (note *Note) MarshalBinary() ([]byte, error) {
	textBytes := []byte(note.Text)
	if len(textBytes) > 1<<31 {
		return nil, errors.New("text is too large")
	}
	buf := make([]byte, 8+8+4+len(textBytes))
	//id
	binary.BigEndian.PutUint64(buf[0:8], note.Id)
	//time
	binary.BigEndian.PutUint64(buf[8:16], uint64(note.Time.UnixNano()))
	//len
	binary.BigEndian.PutUint32(buf[16:20], note.Len)
	//text
	copy(buf[20:], textBytes)
	return buf, nil

}

func (note *Note) UnmarshalBinary(buf []byte) error {
	note.Id = binary.BigEndian.Uint64(buf[:8])
	note.Time = time.Unix(0, int64(binary.BigEndian.Uint64(buf[8:16])))
	note.Len = binary.BigEndian.Uint32(buf[16:20])
	note.Text = string(buf[20 : 20+note.Len])
	return nil
}
