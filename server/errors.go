package main

import "errors"

var (
	ErrorGeneric          = errors.New("we're having some troubles, try again later")
	ErrorPing             = errors.New("did not get pong from the server")
	ErrorShareExpired     = errors.New("this share has expired")
	ErrorShareExpOrPass   = errors.New("this share has expired or the password was incorrect")
	ErrorShareExists      = errors.New("the share id is already taken, try another")
	ErrorSharePassword    = errors.New("the password doesn't match")
	ErrorShareNotEditable = errors.New("this share does not allow editing")
	ErrorRateLimit        = errors.New("your network is making too many requests, try again later")
)
