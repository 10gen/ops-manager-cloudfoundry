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
