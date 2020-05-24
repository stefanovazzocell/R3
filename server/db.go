package main

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

// newPool - Create a new DB Pool
func newPool(hostname string) *redis.Pool {
	log.Println("[DB] Creating new Redis Pool")
	return &redis.Pool{
		MaxIdle:   20,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostname)
			HandleError(err)
			return c, err
		},
	}
}

// FindLink - Attempt to get a link
func FindLink(linkReq LinkRequest) (LinkData, error) {
	var linkData LinkData
	var data string
	var hits int

	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Try to get Data
	conn.Send("MULTI")
	conn.Send("HINCRBY", db_prefix+"link:"+linkReq.ID, "hits", "-1")
	conn.Send("HGET", db_prefix+"link:"+linkReq.ID, "data")
	reply, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return linkData, err
	}
	hits, err = redis.Int(reply[0], err)

	// Deal with expired links
	if hits <= 0 {
		// Attempt to delete and cleanup
		if deleteLink(linkReq) != nil {
			log.Println("Error deleting link")
		}
	}
	if hits < 0 {
		// Don't return data
		return linkData, errors.New("Link Expired")
	}
	
	data, err = redis.String(reply[1], err)
	if err != nil {
		return linkData, err
	}

	linkData.Data = data

	return linkData, nil
}

// SetLink - Attempt to create a link
func SetLink(linkReq LinkRequest) (bool, error) {
	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Attempt to create link
	reply, err := redis.Int(conn.Do("HSETNX", db_prefix+"link:"+linkReq.ID, "data", linkReq.Payload.Data))
	if err != nil {
		HandleError(err)
		return false, err
	}
	if reply != 1 {
		// Link Taken
		return false, nil
	}

	// Complete Creation
	conn.Send("MULTI")
	conn.Send("EXPIRE", db_prefix+"link:"+linkReq.ID, linkReq.Payload.TTL)
	conn.Send("HSET", db_prefix+"link:"+linkReq.ID, "hits", linkReq.Payload.Hits, "edit", linkReq.Payload.Edit)
	_, err = conn.Do("EXEC")
	if err != nil {
		// Attempt a cleanup
		_, _ = conn.Do("DEL", db_prefix+"link:"+linkReq.ID)
		HandleError(err)
		return false, err
	}

	return true, nil
}

// deleteLink - Attempt to delete a link
func deleteLink(linkReq LinkRequest) error {
	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Attempt to delete the link
	_, err := conn.Do("DEL", db_prefix+"link:"+linkReq.ID)
	return err
}

// editLink - Attempt to edit or delete a link
func editLink(linkReq LinkRequest, delete bool) (bool, error) {
	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	// Attempt to check the password
	conn.Send("WATCH", db_prefix+"link:"+linkReq.ID)
	reply, err := redis.String(conn.Do("HGET", db_prefix+"link:"+linkReq.ID, "edit"))
	if err != nil {
		HandleError(err)
		return false, err
	}
	if reply == "" || reply != linkReq.Password {
		// Invalid Password or not editable
		return false, nil
	}

	if delete {
		// Complete Creation
		_, err = conn.Do("UNWATCH")
		if err != nil {
			return false, nil
		}
		err = deleteLink(linkReq)
		if err != nil {
			return false, nil
		}
	} else {
		// Complete Creation
		conn.Send("MULTI")
		conn.Send("EXPIRE", db_prefix+"link:"+linkReq.ID, linkReq.Payload.TTL)
		conn.Send("HSET", db_prefix+"link:"+linkReq.ID, "data", linkReq.Payload.Data, "hits", linkReq.Payload.Hits, "edit", linkReq.Payload.Edit)
		_, err = conn.Do("EXEC")
		if err != nil {
			return false, nil
		}
	}

	return true, nil
}

// addHitIP - Adds a hit to an IP
func addHitIP(ip string, hits int) (int, error) {
	// Get a Connection
	conn := pool.Get()
	defer conn.Close()

	var id string = db_prefix + "ip:" + ip + ":" + strconv.Itoa(time.Now().Hour())

	// Attempt to register hit
	reply, err := redis.Int(conn.Do("INCRBY", id, hits))
	_, _ = conn.Do("EXPIRE", id, 3600)

	if err != nil {
		HandleError(err)
		return 0, err
	}

	return reply, nil
}
