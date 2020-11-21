package sound

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Sound struct {
	ID     uuid.UUID
	Name   string
	Data   [][]byte
	UserID string

	CreatedAt time.Time
	Archived  bool
}

func DataRead(reader io.Reader) (ret [][]byte, err error) {
	for {
		var opusLen int16
		err = binary.Read(reader, binary.LittleEndian, &opusLen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return ret, nil
		} else if err != nil {
			return ret, fmt.Errorf("sound data read: binary read length: %w", err)
		}

		buf := make([]byte, opusLen)
		err = binary.Read(reader, binary.LittleEndian, &buf)
		if err != nil {
			return ret, fmt.Errorf("sound data read: binary read data: %w", err)
		}

		ret = append(ret, buf)
	}
}

func DataWrite(data [][]byte, w io.Writer) (err error) {
	for _, v := range data {
		err = binary.Write(w, binary.LittleEndian, int16(len(v)))
		if err != nil {
			return fmt.Errorf("sound data write: binary write length: %w", err)
		}

		err = binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return fmt.Errorf("sound data write: binary write data: %w", err)
		}
	}

	return nil
}
