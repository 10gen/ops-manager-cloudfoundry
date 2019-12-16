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

package useful

import "log"

// PanicOnUnrecoverableError if the passed error is not nil, panic as the system cannot directly recover from it
// this is a helper for avoiding the _if err!=nil then return err_ pattern, which in the case of certain classes
// of errors, would only make the code harder to read
func PanicOnUnrecoverableError(err error) {
	if err != nil {
		log.Panicf("[ERROR] did not expect a failure, but got: %v", err)
	}
}

// LogError if the passed error is not nil, log it and move on
// this function is useful for example in the case of deferred executions, where we want to let the user know
// something happened, but do not necessarily want to stop the application
func LogError(action func() error) {
	err := action()
	if err != nil {
		log.Printf("[ERROR] unexpected err: %v", err)
	}
}
