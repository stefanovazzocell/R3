package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	mux "github.com/julienschmidt/httprouter"
)

// RateLimit - Request Rate Limiter, returns true if banned
func RateLimit(r *http.Request, linkReq LinkRequest) bool {
	// Guess IP
	var ip string = ""
	if rl_proxy == 2 {
		// Trust CloudFlare
		ip = r.Header.Get("CF-Connecting-IP")
	}
	if rl_proxy == 1 || (rl_proxy == 2 && len(ip) < 7) {
		// Trust Proxy
		ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		if len(ips) > 0 {
			ip = ips[0]
		}
	}
	if len(ip) < 7 {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	// Generate hash
	h := md5.New()
	h.Write([]byte(ip))
	ip = string([]rune(base64.URLEncoding.EncodeToString(h.Sum(nil)))[0:2])

	total, err := addHitIP(ip, int((len(linkReq.Payload.Data)-1)/1024)+1)
	if err != nil || total >= rl_max_ipr {
		// Limit
		return true
	}
	// Normal
	return false
}

// CheckPayload - Checks a request for validity
func CheckRequest(linkReq LinkRequest, checkPayload bool) bool {
	// ID must be an hash
	if len(linkReq.ID) != 88 && len(linkReq.ID) != 8 {
		return false
	}

	// Check Payload
	if checkPayload {
		// Data between 1 and rl_max_req
		if (len(linkReq.Payload.Data) < 1) && (len(linkReq.Payload.Data) > rl_max_req) {
			return false
		}
		// TTL between 1s and 365days
		if (linkReq.Payload.TTL <= 0) && (linkReq.Payload.TTL > 1314000) {
			return false
		}
		// Hits between 1 and 1million
		if (linkReq.Payload.Hits < 1) && (linkReq.Payload.Hits > 1000000) {
			return false
		}
		// If Data is big, then TTL must be within 10 min
		if (len(linkReq.Payload.Data) >= rl_bigquery) && (linkReq.Payload.TTL > 600) {
			return false
		}
		// If Editable, it needs to be an hash
		if len(linkReq.Payload.Edit) != 0 && len(linkReq.Payload.Edit) != 88 {
			return false
		}
	}

	return true
}

// RateLimitrror - Sends a rate limit error response
func RateLimitError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	HandleError(json.NewEncoder(w).Encode(ErrorMsg{"Too many requests. Try again later.", http.StatusTooManyRequests}))
}

// ServerError - Sends a server error response
func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	HandleError(json.NewEncoder(w).Encode(ErrorMsg{"Server error. Try again later.", http.StatusInternalServerError}))
}

// ResourceError - Sends a resource error response
func ResourceError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	HandleError(json.NewEncoder(w).Encode(ErrorMsg{msg, http.StatusUnauthorized}))
}

// RequestError - Sends a request error response
func RequestError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	HandleError(json.NewEncoder(w).Encode(ErrorMsg{"Check your request and try again.", http.StatusBadRequest}))
}

// handleIndex - Handles the index page
func handleIndex(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Rocket API Endpoint")
}

// handleStatus - Handles the status page
func handleStatus(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// handleGet - Handles the API requests for GET
func handleGet(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	var linkReq LinkRequest
	if json.NewDecoder(r.Body).Decode(&linkReq) != nil {
		RequestError(w)
		return
	}

	if CheckRequest(linkReq, false) == false {
		RequestError(w)
		return
	}

	linkData, err := FindLink(linkReq)
	if err != nil {
		ResourceError(w, "Link not found")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	HandleError(json.NewEncoder(w).Encode(linkData))
}

// handleCreate - Handles the API requests for POST (Create)
func handleCreate(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	var linkReq LinkRequest
	if json.NewDecoder(r.Body).Decode(&linkReq) != nil {
		RequestError(w)
		return
	}

	if CheckRequest(linkReq, true) == false {
		RequestError(w)
		return
	}

	if RateLimit(r, linkReq) {
		RateLimitError(w)
		return
	}

	res, err := SetLink(linkReq)
	if err != nil {
		ServerError(w)
		return
	}
	if res == false {
		ResourceError(w, "Link taken.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// handleEdit - Handles the API requests for PUT (Edit)
func handleEdit(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	var linkReq LinkRequest
	if json.NewDecoder(r.Body).Decode(&linkReq) != nil {
		RequestError(w)
		return
	}

	if (CheckRequest(linkReq, true) == false) || (len(linkReq.Password) != 88) {
		RequestError(w)
		return
	}

	res, err := editLink(linkReq, false)
	if err != nil {
		ServerError(w)
		return
	}
	if res == false {
		ResourceError(w, "Wrong password or missing link.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// handleDelete - Handles the API requests for DELETE (Delete)
func handleDelete(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	var linkReq LinkRequest
	if json.NewDecoder(r.Body).Decode(&linkReq) != nil {
		RequestError(w)
		return
	}

	res, err := editLink(linkReq, true)
	if err != nil {
		ServerError(w)
		return
	}
	if res == false {
		ResourceError(w, "Wrong password or missing link.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
