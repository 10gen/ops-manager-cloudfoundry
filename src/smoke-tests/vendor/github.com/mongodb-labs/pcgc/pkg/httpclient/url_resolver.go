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

import (
	"fmt"
	"net/url"
	"path"

	"github.com/mongodb-labs/pcgc/pkg/useful"
)

type baseURL struct {
	base   *url.URL
	prefix string
}

// URLResolver contract for resolving any paths against the given base URL
type URLResolver interface {
	Of(path string, v ...interface{}) string
	OfUnprefixed(path string, v ...interface{}) string
}

// NewURLResolver builds a new API URL which can be used to build any path
func NewURLResolver(base string) URLResolver {
	return parseBaseURL(base)
}

// NewURLResolverWithPrefix builds a new API URL using a prefix for all other paths
func NewURLResolverWithPrefix(base string, prefix string) URLResolver {
	result := parseBaseURL(base)
	result.prefix = prefix
	return result
}

// Of builds a URL by concatenating the base URL and prefix with the specified path, replacing all expansions
func (u baseURL) Of(apiPath string, expansions ...interface{}) string {
	expandedPath := fmt.Sprintf(apiPath, expansions...)
	if u.prefix != "" {
		// add the prefix, if required
		expandedPath = path.Join(u.prefix, expandedPath)
	}

	result, err := u.base.Parse(path.Clean(expandedPath))
	useful.PanicOnUnrecoverableError(err)
	return result.String()
}

// OfUnprefixed builds a URL by concatenating the base URL with the specified path, replacing all expansions
func (u baseURL) OfUnprefixed(apiPath string, expansions ...interface{}) string {
	expandedPath := fmt.Sprintf(apiPath, expansions...)

	result, err := u.base.Parse(path.Clean(expandedPath))
	useful.PanicOnUnrecoverableError(err)
	return result.String()
}

func parseBaseURL(base string) baseURL {
	result := baseURL{}
	var err error
	result.base, err = url.Parse(base)
	useful.PanicOnUnrecoverableError(err)
	return result
}
