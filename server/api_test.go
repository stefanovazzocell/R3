package main

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stefanovazzocell/R3/clientlib"
	"github.com/stefanovazzocell/R3/shared"
)

func TestAPI(t *testing.T) {
	// Setup
	go main()
	dataA := randomBytes(uint(shared.SmallQuery - 100))
	dataB := randomBytes(uint(shared.SmallQuery - 100))
	key := "secret"
	editPass := "edit_pass"
	time.Sleep(time.Second)
	// Create Link
	log.Println("Creating a link")
	err := editLink(dataA, key, "", 60, 10, editPass, false)
	if err != nil {
		t.Errorf("Error creating a link: %v\n", err)
		return
	}
	// Edit Link, wrong pass!
	log.Println("Editing a link with the wrong password")
	err = editLink(dataB, key, "not"+editPass, 60, 10, editPass+"!", false)
	if err == nil {
		t.Error("Fail to error on edit a link with wrong pass")
		return
	}
	// Get Link
	log.Println("Getting a link")
	data, hits, ttl, err := getLink(key, "")
	if err != nil {
		t.Errorf("Error editing a link: %v\n", err)
		return
	}
	if (hits == 0) && (ttl == 0) && !bytes.Equal(data, dataA) {
		t.Errorf("Error getting a link, data doesn't match '%s', hits=%d, ttl=%d\n", data, hits, ttl)
		return
	}
	// Check Link
	log.Println("Checking a link [1/2]")
	data, hits, ttl, err = getLink(key, editPass)
	if err != nil {
		t.Errorf("Error checking [1/2] a link: %v\n", err)
		return
	}
	if (hits == 9) && (ttl > 0) && (ttl < 60) && !bytes.Equal(data, []byte{}) {
		t.Errorf("Error checking [1/2] a link, data doesn't match '%s', hits=%d, ttl=%d\n", data, hits, ttl)
		return
	}
	// Edit Link
	log.Println("Editing a link")
	err = editLink(dataB, key, editPass, 120, 15, editPass+"!", false)
	if err != nil {
		t.Errorf("Error editing a link: %v\n", err)
		return
	}
	// Failing to check Link
	log.Println("Checking a link with the wrong pass")
	_, _, _, err = getLink(key, editPass)
	if err == nil {
		t.Errorf("Error checking a link with the wrong password (expected fail): %v\n", err)
		return
	}
	// Check Link
	log.Println("Checking a link [2/2]")
	data, hits, ttl, err = getLink(key, editPass+"!")
	if err != nil {
		t.Errorf("Error checking [2/2] a link: %v\n", err)
		return
	}
	if (hits == 15) && (ttl > 60) && (ttl <= 120) && !bytes.Equal(data, []byte{}) {
		t.Errorf("Error checking [2/2] a link, data doesn't match '%s', hits=%d, ttl=%d\n", data, hits, ttl)
		return
	}
	// Delete Link
	log.Println("Deleting a link")
	err = editLink(dataB, key, editPass+"!", 0, 0, "", true)
	if err != nil {
		t.Errorf("Error deleting a link: %v\n", err)
		return
	}
}

func getLink(key string, pass string) ([]byte, int, int, error) {
	var hc chan string = make(chan string, 1)
	var phc chan string = make(chan string, 1)
	var pc chan []byte = make(chan []byte, 1)
	defer close(hc)
	defer close(phc)
	defer close(pc)

	go clientlib.GenHash(key, hc)
	if pass != "" {
		go clientlib.GenPass(pass, phc)
	} else {
		phc <- ""
	}

	ar, err := clientlib.GetLink("http://localhost:8080", hc, phc)
	if err != nil {
		return []byte{}, 0, 0, err
	}

	if !ar.Success {
		return []byte{}, 0, 0, errors.New("view response not ok: '" + ar.Err + "'")
	}
	if len(ar.Data) > 0 {
		go clientlib.Decrypt(ar.Data, key, pc)
		return (<-pc), ar.Hits, ar.TTL, nil
	}
	return []byte{}, ar.Hits, ar.TTL, nil
}

func editLink(data []byte, key string, pass string, ttl int, hits int, edit string, delete bool) error {
	var kc chan []byte = make(chan []byte, 1)
	var sc chan []byte = make(chan []byte, 1)
	var dc chan string = make(chan string, 1)
	var hc chan string = make(chan string, 1)
	var phc chan string = make(chan string, 1)
	var ehc chan string = make(chan string, 1)
	defer close(kc)
	defer close(sc)
	defer close(dc)
	defer close(hc)
	defer close(ehc)
	defer close(phc)

	if pass != "" {
		go clientlib.GenPass(pass, phc)
	} else {
		phc <- ""
	}
	if edit != "" {
		go clientlib.GenPass(edit, ehc)
	} else {
		ehc <- ""
	}
	go clientlib.GenHash(key, hc)
	go clientlib.KeyDerivation(key, kc, sc)
	go clientlib.Encrypt(data, kc, sc, dc)

	ar, err := clientlib.EditLink("http://localhost:8080", hc, phc, delete, dc, ttl, hits, ehc)
	if err != nil {
		return err
	}
	if !ar.Success {
		return errors.New("edit response not ok: '" + ar.Err + "'")
	}
	return nil
}

func randomBytes(n uint) []byte {
	rand.Seed(time.Now().UnixNano())

	outBytes := []byte{}
	for i := uint(0); i < n; i++ {
		outBytes = append(outBytes, byte(rand.Intn(256)))
	}
	return outBytes
}
