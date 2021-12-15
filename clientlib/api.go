package clientlib

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stefanovazzocell/R3/shared"
)

func GetLink(apiEndpoint string, hc chan string, phc chan string) (shared.APIResponseData, error) {
	var ar shared.APIResponseData = shared.APIResponseData{}

	var vr shared.ViewRequest = shared.ViewRequest{}
	vr.ID = <-hc
	vr.Password = <-phc
	if !vr.Verify() {
		return ar, ErrorValidation
	}
	jsonData, _ := json.Marshal(vr)

	request, err := http.NewRequest("POST", apiEndpoint+"/v2/get", bytes.NewBuffer(jsonData))
	if err != nil {
		return ar, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return ar, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ar, ErrorNetwork
	}

	err = json.NewDecoder(response.Body).Decode(&ar)
	return ar, err
}

func EditLink(apiEndpoint string, hc chan string, phc chan string, delete bool, dc chan string, ttl int, hits int, ehc chan string) (shared.APIResponse, error) {
	var ar shared.APIResponse = shared.APIResponse{}

	var er shared.EditRequest = shared.EditRequest{}
	er.Delete = delete
	er.Payload.TTL = ttl
	er.Payload.Hits = hits
	er.Password = <-phc
	er.Payload.Edit = <-ehc
	er.ID = <-hc
	er.Payload.Data = <-dc
	if !er.Verify() {
		return ar, ErrorValidation
	}
	jsonData, _ := json.Marshal(er)

	request, err := http.NewRequest("POST", apiEndpoint+"/v2/edit", bytes.NewBuffer(jsonData))
	if err != nil {
		return ar, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return ar, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ar, ErrorNetwork
	}

	err = json.NewDecoder(response.Body).Decode(&ar)
	return ar, err
}
