package shared

import "encoding/base64"

type Config struct {
	ListeningAddr  string `json:"ListeningAddr"`
	CORSorigin     string `json:"CORSorigin"`
	RedisURI       string `json:"RedisURI"`
	DatabasePrefix string `json:"DatabasePrefix"`
	DataCap        int    `json:"DataCap"`
	QueriesCap     int    `json:"QueriesCap"`
	HasProxy       bool   `json:"HasProxy"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
}

type APIResponseData struct {
	Success bool   `json:"success"`
	Err     string `json:"error"`
	Data    string `json:"data"` // Optional
	Hits    int    `json:"hits"` // Optional
	TTL     int    `json:"ttl"`  // Optional
}

type ViewRequest struct {
	ID       string `json:"id"`
	Password string `json:"pass"` // Optional
}

func (req ViewRequest) Verify() bool {
	return len(req.ID) == 8 && (len(req.Password) == 0 || len(req.Password) == 88)
}

type EditRequest struct {
	ID       string    `json:"id"`
	Delete   bool      `json:"delete"`  // Optional
	Password string    `json:"pass"`    // Optional
	Payload  ShareData `json:"payload"` // Optional
}

func (req EditRequest) Verify() bool {
	// Note: more checks must be done while converting data
	if len(req.ID) != 8 {
		return false
	}
	if len(req.Password) != 0 && len(req.Password) != 88 {
		return false
	}
	if req.Delete && len(req.Password) != 88 {
		return false
	}
	return req.Delete || req.Payload.Verify()
}

func (req EditRequest) GetDataBytes() ([]byte, error) {
	if req.Delete {
		return []byte{}, nil
	}
	data, err := base64.StdEncoding.DecodeString(req.Payload.Data)
	if err != nil {
		return data, err
	}
	if (len(data) > SmallQuery) && (len(data) > LargeQuery || req.Payload.TTL > 60*60) {
		return []byte{}, ErrorRequest
	}
	return data, err
}

type ShareData struct {
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
	Hits int    `json:"hits"`
	Edit string `json:"edit"`
}

func (data ShareData) Verify() bool {
	// Note: more checks must be done while converting data
	if len(data.Data) < 50 || len(data.Data) > LargeQuery*100 {
		return false
	}
	if data.TTL < 10 || data.TTL > 7*24*60*60 {
		return false
	}
	if data.Hits < 1 || data.Hits > 1000000 {
		return false
	}
	return (len(data.Edit) == 0 || len(data.Edit) == 88)
}
