package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/stefanovazzocell/R3/clientlib"
)

var typeFromFile = map[string]bool{
	"url":  false,
	"text": false,
	"file": true,
}

func main() {
	APIEndpoint := queryString("What is the API Endpoint [enter for localhost]? ")
	if len(APIEndpoint) == 0 {
		APIEndpoint = "http://localhost:8080"
	}
	for {
		// Clear and Init
		var pc chan []byte = make(chan []byte, 1)
		var kc chan []byte = make(chan []byte, 1)
		var sc chan []byte = make(chan []byte, 1)
		var dc chan string = make(chan string, 1)
		var hc chan string = make(chan string, 1)
		var phc chan string = make(chan string, 1)
		var ehc chan string = make(chan string, 1)
		defer close(pc)
		defer close(kc)
		defer close(sc)
		defer close(dc)
		defer close(hc)
		defer close(phc)
		defer close(ehc)
		var (
			shareID     string = ""
			sharePass   string = ""
			shareEdit   string = ""
			shareTTL    int    = 1
			shareHits   int    = 1
			shareDelete bool   = false
		)
		// Action?
		action := strings.ToLower(queryString("\nWhat could you like to do [view,create,edit,delete,stats,quit]? "))
		if action == "q" || action == "quit" {
			break
		}
		// Gather Info
		shareID = queryString("What is the share ID? ")
		if len(shareID) < 1 {
			continue
		}
		go clientlib.GenHash(shareID, hc)
		go clientlib.KeyDerivation(shareID, kc, sc)
		switch action {
		case "v", "view":
			phc <- ""
		case "c", "create":
			shareEdit = queryString("What is the share edit password [empty disables management]? ")
			if shareEdit == "" {
				ehc <- ""
			} else {
				go clientlib.GenPass(shareEdit, ehc)
			}
			go clientlib.Encrypt(queryShareData(), kc, sc, dc)
			shareHits = queryInt("How many hits to allow on this share? ", 1, 1000000)
			shareTTL = queryInt("What is the share TTL [in seconds, limited to 1h]? ", 1, 3600)
			phc <- ""
		case "e", "edit":
			sharePass = queryString("What is the share password? ")
			go clientlib.GenPass(sharePass, phc)
			shareEdit = queryString("What is the share edit password [empty disables management]? ")
			if shareEdit == "" {
				ehc <- ""
			} else {
				go clientlib.GenPass(shareEdit, ehc)
			}
			go clientlib.Encrypt(queryShareData(), kc, sc, dc)
			shareHits = queryInt("How many hits to allow on this share? ", 1, 1000000)
			shareTTL = queryInt("What is the share TTL [in seconds, limited to 1h]? ", 1, 3600)
		case "d", "delete":
			shareDelete = true
			sharePass = queryString("What is the share password? ")
			go clientlib.GenPass(sharePass, phc)
			dc <- ""
			ehc <- ""
		case "s", "stats":
			sharePass = queryString("What is the share password? ")
			go clientlib.GenPass(sharePass, phc)
		}

		// Call API
		switch action {
		case "v", "view", "s", "stats":
			fmt.Println("Querying server...")
			ar, err := clientlib.GetLink(APIEndpoint, hc, phc)
			if err != nil {
				fmt.Printf("Client Error: %v\n", err)
			} else if !ar.Success {
				fmt.Printf("Server Error: %v\n", ar.Err)
			} else if len(ar.Data) > 0 {
				go clientlib.Decrypt(ar.Data, shareID, pc)
				st, sd := clientlib.ShareDecode(<-pc)
				if typeFromFile[st] {
					fmt.Printf("Share is of type %s", st)
					for {
						var file clientlib.File
						sd, err = file.Decode(sd)
						if err == nil {
							fmt.Printf("[File] name: '%s', type: '%s', data length: %d\n", file.Name, file.MimeType, len(file.Data))
						} else {
							fmt.Printf("[Corrupted File] %s\n", err)
						}
						if len(sd) == 0 {
							break
						}
					}
				} else {
					fmt.Printf("Share is of type %s, data: '%s'\n", st, sd)
				}
			} else {
				fmt.Printf("Share will expire in %s, and can be viewed %d more times.\n", time.Second*time.Duration(ar.TTL), ar.Hits)
			}
		case "c", "create", "e", "edit", "d", "delete":
			fmt.Println("Querying server...")
			ar, err := clientlib.EditLink(APIEndpoint, hc, phc, shareDelete, dc, shareTTL, shareHits, ehc)
			if err != nil {
				fmt.Printf("Client Error: %v\n", err)
			} else if !ar.Success {
				fmt.Printf("Server Error: %v\n", ar.Err)
			} else {
				fmt.Println("Done!")
			}
		}
	}
}

func queryString(query string) string {
	fmt.Print(query)
	var line string
	fmt.Scanln(&line)
	return line
}
func queryInt(query string, min int, max int) int {
	fmt.Print(query)
	var res string
	for {
		fmt.Scanln(&res)
		i, err := strconv.Atoi(res)
		if err != nil || i > max || i < min {
			fmt.Printf("please enter a number between %d and %d: ", min, max)
		} else {
			return i
		}
	}
}

func queryShareData() []byte {
	fmt.Print("What is the type of your share [url,text,file]? ")
	var st string
	var (
		ff bool
		ok bool
	)
	for {
		fmt.Scanln(&st)
		st = strings.ToLower(st)
		ff, ok = typeFromFile[st]
		if ok {
			break
		}
		fmt.Print("please enter url, text, or file: ")
	}
	if !ff {
		return clientlib.ShareEncode(st, []byte(queryString("Please enter the "+st+": ")))
	}
	data := []byte{}
	for {
		path := queryString("Please enter a path to your next " + st + " (blank for done): ")
		if path == "" {
			break
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Client Error: %v\n", err)
		} else {
			var file clientlib.File
			file.Load(path, content)
			data = append(data, file.Encode()...)
		}
	}
	return clientlib.ShareEncode(st, data)
}
