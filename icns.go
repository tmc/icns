package icns

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// IconSet represents a group of icons.
type IconSet []*Icon

// IconType is a type of icon.
type IconType [4]byte

// String returns the string representation of an IconType.
func (it IconType) String() string {
	return string(it[:])
}

// Icon is a single icon
type Icon struct {
	Type IconType
	Data []byte
}

type iconSetHeader struct {
	Magic [4]byte
	Len   int32
}

type iconHeader struct {
	Type [4]byte
	Len  int32
}

// Parse attempts to parse a given Reader into an IconSet.
func Parse(r io.Reader) (IconSet, error) {
	h := iconSetHeader{}

	if err := binary.Read(r, binary.LittleEndian, &h.Magic); err != nil {
		return nil, errors.Wrap(err, "parsing magic header")
	}
	if err := binary.Read(r, binary.BigEndian, &h.Len); err != nil {
		return nil, errors.Wrap(err, "parsing header size")
	}
	buf := make([]byte, int(h.Len-8))
	if _, err := io.ReadAtLeast(r, buf, int(h.Len-8)); err != nil {
		return nil, errors.Wrap(err, "reading icon data")
	}

	var result IconSet
	br := bytes.NewReader(buf)
	for {
		i, err := parseIcon(br)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "issue parsing icon")
		}
		result = append(result, i)
	}
	// TODO: ensure all data is parsed?

	return result, nil
}

func parseIcon(r io.Reader) (*Icon, error) {
	ih := iconHeader{}
	if err := binary.Read(r, binary.LittleEndian, &ih.Type); err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, errors.Wrap(err, "parsing icon type")
	}
	if err := binary.Read(r, binary.BigEndian, &ih.Len); err != nil {
		return nil, errors.Wrap(err, "parsing icon size")
	}
	result := &Icon{
		Type: ih.Type,
		Data: make([]byte, ih.Len-8),
	}
	if _, err := io.ReadAtLeast(r, result.Data, int(ih.Len-8)); err != nil {
		return nil, err
	}
	return result, nil
}
