package httpclient

import "time"

// RequestTimeouts allows users of this code to tweak all necessary timeouts used during an HTTP request-response
type RequestTimeouts struct {
	DialTimeout           time.Duration
	ExpectContinueTimeout time.Duration
	IdleConnectionTimeout time.Duration
	ResponseHeaderTimeout time.Duration
	TLSHandshakeTimeout   time.Duration
	// GlobalTimeout the maximum allowed duration to complete a single HTTP request and response
	GlobalTimeout time.Duration
}

// NewDefaultTimeouts initializes the timeouts struct using default values
func NewDefaultTimeouts() *RequestTimeouts {
	return &RequestTimeouts{
		DialTimeout:           10 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		IdleConnectionTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		GlobalTimeout:         30 * time.Second,
	}
}
