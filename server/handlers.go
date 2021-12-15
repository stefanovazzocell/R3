package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/stefanovazzocell/R3/shared"
)

func handleEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req shared.EditRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error while decoding EditRequest: %v\n", err)
		errorResponse(w, shared.ErrorRequest)
		return
	}
	if !req.Verify() {
		errorResponse(w, shared.ErrorRequest)
		return
	}
	dataDecoded, err := req.GetDataBytes()
	if err != nil {
		log.Printf("Error getting bytes: %v\n", err)
		errorResponse(w, shared.ErrorRequest)
		return
	}

	err = redisEditShare(r, req, dataDecoded)
	if err != nil {
		errorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(shared.APIResponse{Success: true, Err: ""})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func handleGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req shared.ViewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error while decoding ViewRequest: %v\n", err)
		errorResponseData(w, shared.ErrorRequest)
		return
	}
	if !req.Verify() {
		errorResponseData(w, shared.ErrorRequest)
		return
	}

	resp, err := redisGetShare(r, req)
	if err != nil {
		errorResponseData(w, err)
		return
	}

	json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func errorResponse(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(shared.APIResponse{Success: false, Err: err.Error()})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func errorResponseData(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(shared.APIResponseData{Success: false, Err: err.Error(), Data: "", Hits: 0, TTL: 0})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
