package note

import (
	"encoding/binary"
	"errors"
	"time"
)

type Note struct {
	Id       uint16
	Time     time.Time
	Len      uint32
	Text     string
	Progress byte // 001 - not started, 010 - in progress, 100 - done
}

func NewNote(id uint16, text string) *Note {
	return &Note{Id: id, Time: time.Now(), Text: text, Progress: 1}
}

func (note *Note) MarshalBinary() ([]byte, error) {
	if note.Progress != 1 && note.Progress != 2 && note.Progress != 4 {
		return nil, errors.New("unavaliable progress status")
	}
	textBytes := []byte(note.Text)
	if len(textBytes) > 1<<31 {
		return nil, errors.New("text is too large")
	}
	//id
	buf := make([]byte, 2+8+4+len(textBytes)+1)
	binary.BigEndian.PutUint16(buf[0:2], note.Id)
	//time
	binary.BigEndian.PutUint64(buf[2:10], uint64(note.Time.UnixNano()))
	//len
	binary.BigEndian.PutUint32(buf[10:14], note.Len)
	//text
	buf = append(buf, textBytes...)
	//progress
	buf = append(buf, note.Progress)
	return buf, nil

}

func (note *Note) UnmarshalBinary(buf []byte) error {
	note.Id = binary.BigEndian.Uint16(buf[0:2])
	note.Time = time.Unix(0, int64(binary.BigEndian.Uint64(buf[2:10])))
	note.Len = binary.BigEndian.Uint32(buf[10:14])
	note.Text = string(buf[14 : 14+note.Len])
	note.Progress = buf[len(buf)-1]
	if note.Progress != 1 && note.Progress != 2 && note.Progress != 4 {
		return errors.New("unavaliable progress status")
	}
	return nil
}
