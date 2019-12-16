// Copyright 2019 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
