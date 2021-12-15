package clientlib

import (
	"bytes"
	"testing"
)

func TestFile(t *testing.T) {
	var testFile File
	testFile.Load("hello.txt", []byte("Hello World"))
	var encoded []byte = testFile.Encode()
	var testDecodedFile File
	extra, err := testDecodedFile.Decode(encoded)
	if err != nil {
		t.Fatalf("Decoded file failed with error '%s', expected none\n", err)
	}
	if len(extra) != 0 {
		t.Fatalf("Decoded file had %d extra bytes, expected none\n", len(extra))
	}
	if testDecodedFile.Name != "hello.txt" {
		t.Fatalf("Decoded file had name '%s', expected 'hello.txt'\n", testDecodedFile.Name)
	}
	if testDecodedFile.MimeType != "text/plain; charset=utf-8" {
		t.Fatalf("Decoded file had mime '%s', expected 'text/plain; charset=utf-8'\n", testDecodedFile.MimeType)
	}
	if !bytes.Equal(testDecodedFile.Data, []byte("Hello World")) || string(testDecodedFile.Data) != "Hello World" {
		t.Fatalf("Decoded file had mime '%s' (%b), expected 'Hello World'\n", testDecodedFile.Data, testDecodedFile.Data)
	}
}
