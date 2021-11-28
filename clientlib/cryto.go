package clientlib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func GenPass() string {
	rand.Seed(time.Now().UnixNano())
	var output string = ""
	var validChars []rune = []rune("abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789")

	for i := 0; i < 15; i++ {
		r := rand.Intn(len(validChars) + 1)
		if i == 3 || i == 7 || i == 11 {
			output += "."
		} else if r < len(validChars) {
			output += string(validChars[rand.Intn(len(validChars))])
		}
	}

	return output
}

func GenHash(keyStr string, hc chan string) {
	var HASH_SALT []byte = []byte{82, 242, 11, 190, 119, 15, 58, 152, 115, 230, 184, 149, 107, 12, 5, 37, 184, 242, 159, 111, 72, 180, 65, 53, 104, 78, 252, 123, 188, 17, 71, 187, 216, 128, 141, 148, 126, 110, 15, 113, 175, 70, 216, 37, 211, 247, 93, 216, 210, 197, 189, 100, 37, 81, 113, 113, 173, 8, 184, 97, 225, 223, 24, 69}

	hc <- base64.StdEncoding.EncodeToString(pbkdf2.Key([]byte("hashgen::"+keyStr+"::salty"), HASH_SALT, 1000000, 64, sha512.New)[0:6])
}

func KeyDerivation(keyStr string, kc chan []byte, sc chan string) {
	// Generate Salt
	rand.Seed(time.Now().UnixNano())
	var KEY_SALT []byte = []byte{}
	for i := 0; i < 32; i++ {
		KEY_SALT = append(KEY_SALT, byte(rand.Intn(256)))
	}
	// Key derivation
	kc <- pbkdf2.Key([]byte("keyder::"+keyStr+"::salted"), KEY_SALT, 1000000, 32, sha256.New)
	sc <- base64.StdEncoding.EncodeToString(KEY_SALT)
}

func Encrypt(data []byte, keyData []byte, dc chan string) {
	rand.Seed(time.Now().UnixNano())

	block, err := aes.NewCipher(keyData)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	fmt.Println(gcm.NonceSize())
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	dc <- base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, data, nil))
}

func Decrypt(data string, keyStr string, pc chan []byte) {
	keySalt, err := base64.StdEncoding.DecodeString(data[0:44])
	if err != nil {
		pc <- []byte{}
	}
	encryptedData, err := base64.StdEncoding.DecodeString(data[44:])
	if err != nil {
		pc <- []byte{}
	}

	key := pbkdf2.Key([]byte("keyder::"+keyStr+"::salted"), keySalt, 1000000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		pc <- []byte{}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		pc <- []byte{}
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		pc <- []byte{}
	}

	pc <- plaintext
}
