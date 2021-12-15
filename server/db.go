package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stefanovazzocell/R3/shared"
)

var pool *redis.Pool

func newPool(redisUri string) *redis.Pool {
	log.Println("[DB] Creating new Redis Pool")
	return &redis.Pool{
		MaxIdle:   20,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisUri)
			if err != nil {
				panic(err)
			}
			return c, err
		},
	}
}

func redisGetShare(httpReq *http.Request, req shared.ViewRequest) (shared.APIResponseData, error) {
	var (
		manage   bool                   = (len(req.Password) == 88)
		response shared.APIResponseData = shared.APIResponseData{}
	)

	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Rate Limit
	err := redisHit(conn, httpReq, 0)
	if err != nil {
		return response, err
	}

	// Try to get data
	conn.Send("MULTI")
	if manage {
		conn.Send("GET", DatabasePrefix+":spass:"+req.ID)
		conn.Send("TTL", DatabasePrefix+":sdata:"+req.ID)
		conn.Send("GET", DatabasePrefix+":shits:"+req.ID)
		rrv, err := redis.Values(conn.Do("EXEC"))
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while retriving share data: %v\n", err)
			return response, ErrorGeneric
		}
		if err == redis.ErrNil {
			return response, ErrorShareExpOrPass
		}
		// Check Password
		pass, err := redis.String(rrv[0], err)
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while retriving share data password: %v\n", err)
			return response, ErrorGeneric
		}
		if pass != req.Password {
			return response, ErrorSharePassword
		}
		// Get Metadata
		ttl, err := redis.Int(rrv[1], err)
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while retriving share data ttl: %v\n", err)
			return response, ErrorGeneric
		}
		hits, err := redis.Int(rrv[2], err)
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while retriving share data hits: %v\n", err)
			return response, ErrorGeneric
		}
		response.Success = true
		response.Err = ""
		response.TTL = ttl
		response.Hits = hits
		return response, nil
	} else {
		conn.Send("INCRBY", DatabasePrefix+":shits:"+req.ID, -1)
		conn.Send("GET", DatabasePrefix+":sdata:"+req.ID)
		rrv, err := redis.Values(conn.Do("EXEC"))
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while retriving share data: %v\n", err)
			return response, ErrorGeneric
		}
		if err == redis.ErrNil {
			return response, ErrorShareExpired
		}
		// Check hits
		hits, err := redis.Int(rrv[0], err)
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while checking share data hits: %v\n", err)
			return response, ErrorGeneric
		}
		if hits < 0 {
			return response, ErrorShareExpired
		}
		// Get data
		data, err := redis.Bytes(rrv[1], err)
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while parsing share data: %v\n", err)
			return response, ErrorGeneric
		}
		response.Success = true
		response.Err = ""
		response.Data = base64.StdEncoding.EncodeToString(data)
		// Expire share
		if hits <= 0 {
			conn.Send("MULTI")
			conn.Send("DEL", DatabasePrefix+":sdata:"+req.ID)
			conn.Send("DEL", DatabasePrefix+":spass:"+req.ID)
			conn.Send("DEL", DatabasePrefix+":shits:"+req.ID)
			conn.Do("EXEC")
		}
		return response, nil
	}
}

func redisEditShare(httpReq *http.Request, req shared.EditRequest, dataDecoded []byte) error {
	editing := (req.Password != "")

	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Check if share exists
	data, err := redis.Bytes(conn.Do("GET", DatabasePrefix+":sdata:"+req.ID))
	if err != nil && err != redis.ErrNil {
		log.Printf("Error while checking if share exists: %v\n", err)
		return ErrorGeneric
	}
	shareExists := (len(data) > 0)

	// If editing/deleting it must exist, otherwise it shouldn't
	if editing && !shareExists {
		return ErrorShareExpired
	}
	if len(req.Password) == 0 && shareExists {
		return ErrorShareExists
	}

	// Rate limit check
	if editing && !req.Delete || editing {
		cost := len(dataDecoded) / 1000
		if cost < 1 {
			cost = 1
		}
		err := redisHit(conn, httpReq, cost)
		if err != nil {
			return err
		}
	}

	// If editing, check that the password matches
	if editing {
		data, err := redis.String(conn.Do("GET", DatabasePrefix+":spass:"+req.ID))
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while checking password match: %v\n", err)
			return ErrorGeneric
		}
		if err == redis.ErrNil || data == "" {
			return ErrorShareNotEditable
		} else if data != req.Password {
			return ErrorSharePassword
		}
	}

	// Either delete or save
	if editing && req.Delete {
		conn.Send("MULTI")
		conn.Send("DEL", DatabasePrefix+":sdata:"+req.ID)
		conn.Send("DEL", DatabasePrefix+":spass:"+req.ID)
		conn.Send("DEL", DatabasePrefix+":shits:"+req.ID)
		_, err := conn.Do("EXEC")
		if err != nil && err != redis.ErrNil {
			log.Printf("Error while deleting: %v\n", err)
			return ErrorGeneric
		}
	} else {
		conn.Send("MULTI")
		// Set Data
		conn.Send("SETEX", DatabasePrefix+":sdata:"+req.ID, req.Payload.TTL, dataDecoded)
		// Either set or clear edit password
		if len(req.Payload.Edit) == 88 {
			conn.Send("SETEX", DatabasePrefix+":spass:"+req.ID, req.Payload.TTL, req.Payload.Edit)
		} else {
			conn.Send("DEL", DatabasePrefix+":spass:"+req.ID)
		}
		// Set hits limit
		conn.Send("SETEX", DatabasePrefix+":shits:"+req.ID, req.Payload.TTL, req.Payload.Hits)
		_, err := conn.Do("EXEC")
		if err != nil {
			log.Printf("Error while creating/editing: %v\n", err)
			return ErrorGeneric
		}
	}
	return nil
}

func redisHit(conn redis.Conn, r *http.Request, data int) error {
	// Note: if data == 0, it will consider this a request
	var (
		ip     string = ""
		id     string
		err    error
		count  int
		isData bool = (data != 0)
	)
	// Get IP
	if HasProxy {
		ip = r.Header.Get("CF-Connecting-IP")
		if ip == "" {
			ip = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		}
	}
	if !HasProxy || ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	hash := md5.Sum([]byte(fmt.Sprintf("%s%d", ip, time.Now().Hour())))
	id = base64.URLEncoding.EncodeToString(hash[0:2])
	// Count
	if isData {
		count, err = redis.Int(conn.Do("INCRBY", DatabasePrefix+":datacap:"+id, data))
		conn.Do("EXPIRE", DatabasePrefix+":datacap:"+id, 3600)
	} else {
		count, err = redis.Int(conn.Do("INCRBY", DatabasePrefix+":queriescap:"+id, 1))
		conn.Do("EXPIRE", DatabasePrefix+":queriescap:"+id, 3600)
	}
	if err != nil {
		log.Printf("Error while setting ratelimit: %v\n", err)
		return ErrorGeneric
	}
	// Check
	if isData {
		if count > DataCap {
			return ErrorRateLimit
		}
	} else {
		if count > QueriesCap {
			return ErrorRateLimit
		}
	}
	return nil
}

func redisPing() error {
	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Attempt to ping
	output, err := redis.String(conn.Do("PING"))
	if err == nil && output != "PONG" {
		err = ErrorPing
	}
	if err == nil && output == "PONG" {
		log.Println("[DB] Received PONG")
	}
	return err
}
