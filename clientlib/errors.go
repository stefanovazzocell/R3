package clientlib

import "errors"

var (
	ErrorNetwork       = errors.New("network error, try again later")
	ErrorValidation    = errors.New("data validation error, check your inputs")
	errorMalformedData = errors.New("the data is malformed, cannot read this share")
)
