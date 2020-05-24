package main

// LinkData - Represents a Link
type LinkData struct {
	Data string `json:"data"`
	Hits int    `json:"hits"`
	TTL  int    `json:"ttl"`
	Edit string `json:"edit"`
}

// LinkRequest - Represents a request for a Link
type LinkRequest struct {
	ID       string   `json:"id"`
	Password string   `json:"password"`
	Payload  LinkData `json:"payload"`
}

// ErrorMsg - Represents an error
type ErrorMsg struct {
	Message string `json:"error"`
	Status  int    `json:"status"`
}
