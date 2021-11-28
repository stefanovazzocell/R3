package clientlib

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestGenPass(t *testing.T) {
	pass := GenPass()
	matched, err := regexp.Match("^([a-zA-Z0-9]{0,3}\\.){3}[a-zA-Z0-9]{0,3}$", []byte(pass))

	if err != nil {
		t.Fatalf("Got error (%v) while checking generated password '%s'\n", err, pass)
	}
	if !matched {
		t.Fatalf("Password '%s' didn't match\n", pass)
	}
}

func TestGenHash(t *testing.T) {
	var hc chan string = make(chan string, 1)
	defer close(hc)
	go GenHash("secret", hc)
	output := <-hc
	if output != "+dSrTYRf" {
		t.Fatalf("Expected '+dSrTYRf' for GenHash of 'secret', got '%s'\n", output)
	}

	var random string = string(randomBytes(20))

	start := time.Now()

	go GenHash(random, hc)
	output = <-hc

	elapsed := time.Since(start)

	if len(output) != 8 {
		t.Fatalf("Expected len 8 for GenHash of '%s', got '%s' (len %d)\n", random, output, len(output))
	}

	if elapsed > time.Second*2 {
		t.Fatalf("Error encrypting/decrypting, it took %s to GenHash", elapsed)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := string(randomBytes(20))
	data := randomBytes(100)

	var kc chan []byte = make(chan []byte, 1)
	var sc chan string = make(chan string, 1)
	var dc chan string = make(chan string, 1)
	var pc chan []byte = make(chan []byte, 1)
	defer close(kc)
	defer close(sc)
	defer close(dc)
	defer close(pc)

	var output []byte

	start := time.Now()

	go KeyDerivation(key, kc, sc)
	go Encrypt(data, <-kc, dc)
	go Decrypt((<-sc)+(<-dc), key, pc)
	output = (<-pc)

	elapsed := time.Since(start)

	if !bytes.Equal(data, output) {
		t.Fatalf("Error encrypting/decrypting, key=%v, data=%v, output=%v, (data and output don't match)", key, data, output)
	}

	if elapsed > time.Second*2 {
		t.Fatalf("Error encrypting/decrypting, it took %s to encrypt+decrypt", elapsed)
	}
}

func randomBytes(n uint) []byte {
	rand.Seed(time.Now().UnixNano())

	outBytes := []byte{}
	for i := uint(0); i < n; i++ {
		outBytes = append(outBytes, byte(rand.Intn(256)))
	}
	return outBytes
}
