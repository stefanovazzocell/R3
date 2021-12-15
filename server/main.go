package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stefanovazzocell/R3/shared"
)

// Flags
var ListeningAddr, CORSorigin, RedisURI, DatabasePrefix string
var DataCap, QueriesCap int
var HasProxy bool

func init() {
	// Read config
	var config shared.Config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}
	ListeningAddr = config.ListeningAddr
	CORSorigin = config.CORSorigin
	RedisURI = config.RedisURI
	DatabasePrefix = config.DatabasePrefix
	DataCap = config.DataCap
	QueriesCap = config.QueriesCap
	HasProxy = config.HasProxy
}

func main() {
	log.Println("Starting Server...")
	pool = newPool(RedisURI)
	defer pool.Close()
	err := redisPing()
	if err != nil {
		panic(err)
	}
	router := newRouter()
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(ListeningAddr, router))
}
