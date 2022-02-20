package main

import "testing"

func TestPingDB(t *testing.T) {
	pool = newPool(RedisURI)
	defer pool.Close()
	err := redisPing()
	if err != nil {
		t.Errorf("Error pinging DB: %v", err)
	}
}
