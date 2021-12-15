package clientlib

import (
	"mime"
	"regexp"
)

var typeToByte map[string]byte = map[string]byte{
	"url":  0,
	"text": 1,
	"file": 2,
}
var byteToType map[byte]string = map[byte]string{
	0: "url",
	1: "text",
	2: "file",
}

func ShareEncode(shareType string, data []byte) []byte {
	tB, ok := typeToByte[shareType]
	if !ok {
		return []byte{}
	}
	out := []byte{tB}
	out = append(out, data...)
	return out
}
func ShareDecode(data []byte) (string, []byte) {
	if len(data) == 0 {
		return "error", []byte{}
	}
	return byteToType[data[0]], data[1:]
}

// A file can be concatenated to store multiple files

type File struct {
	Name     string
	MimeType string
	Data     []byte
}

func (file *File) Load(path string, data []byte) {
	rXfilename := regexp.MustCompile(`[^/\\]*$`)
	rXext := regexp.MustCompile(`.[^\.]*$`)
	file.Data = data
	file.Name = rXfilename.FindString(path)
	fileext := rXext.FindString(file.Name)
	file.MimeType = mime.TypeByExtension(fileext)
}
func (file File) Encode() (encoded []byte) {
	// Encoding:
	// {name}{mimetype}{data}
	// Encoding of {name} and {mimetype}
	// [1 byte length][string data][if prev. length = 255, 1 byte length], repeat until all string data encoded
	encoded = encodeString(file.Name)
	encoded = append(encoded, encodeString(file.MimeType)...)
	// Encoding of {data}
	// [3 byte length][data][if prev. length = 16581375, 1 byte length], repeat until all data encoded
	encoded = append(encoded, encodeData(file.Data)...)
	return encoded
}
func (file *File) Decode(data []byte) (extra []byte, err error) {
	file.Name, data, err = decodeString(data)
	if err != nil {
		return []byte{}, err
	}
	file.MimeType, data, err = decodeString(data)
	if err != nil {
		return []byte{}, err
	}
	file.Data, data, err = decodeData(data)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func decodeString(encoded []byte) (s string, extra []byte, err error) {
	var offset int = 0
	var sBytes []byte = []byte{}
	for {
		if len(encoded) < 1 {
			return "", []byte{}, errorMalformedData
		}
		offset, encoded = int(encoded[0]), encoded[1:]
		if len(encoded) < offset {
			return "", []byte{}, errorMalformedData
		}
		if offset > 0 {
			sBytes, encoded = append(sBytes, encoded[0:offset]...), encoded[offset:]
		}
		if offset < 255 {
			break
		}
	}
	return string(sBytes), encoded, nil
}
func encodeString(s string) (encoded []byte) {
	var offset int = 0
	sBytes := []byte(s)
	for {
		if len(sBytes) == 0 {
			if offset == 255 {
				encoded = append(encoded, byte(0))
			}
			break
		}
		offset = 255
		if len(sBytes) < 255 {
			offset = len(sBytes)
		}
		encoded = append(encoded, byte(offset))
		encoded = append(encoded, sBytes[0:offset]...)
		sBytes = sBytes[offset:]
	}
	return encoded
}
func decodeData(encoded []byte) (data []byte, extra []byte, err error) {
	var offset int = 0
	for {
		if len(encoded) < 3 {
			return []byte{}, []byte{}, errorMalformedData
		}
		offset, encoded = int(encoded[0])+int(encoded[1])<<8&0xff00+int(encoded[2])<<16&0xff0000, encoded[3:]
		if len(encoded) < offset {
			return []byte{}, []byte{}, errorMalformedData
		}
		if offset > 0 {
			data, encoded = append(data, encoded[0:offset]...), encoded[offset:]
		}
		if offset < 16777215 {
			break
		}
	}
	return data, encoded, nil
}
func encodeData(d []byte) (encoded []byte) {
	var offset int = 0
	for {
		if len(d) == 0 {
			if offset == 16777215 {
				encoded = append(encoded, byte(0))
			}
			break
		}
		offset = 16777215
		if len(d) < 16777215 {
			offset = len(d)
		}
		encoded = append(encoded, []byte{byte(offset & 0xff), byte(offset >> 8 & 0xff), byte(offset >> 16 & 0xff)}...)
		encoded = append(encoded, d[0:offset]...)
		d = d[offset:]
	}
	return encoded
}
