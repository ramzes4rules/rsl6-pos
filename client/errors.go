package client

import "fmt"

// FaultError represents a SOAP fault returned by the service.
type FaultError struct {
	Code   string
	Reason string
	Detail string
}

func (e *FaultError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("soap fault: [%s] %s (detail: %s)", e.Code, e.Reason, e.Detail)
	}
	return fmt.Sprintf("soap fault: [%s] %s", e.Code, e.Reason)
}

// IsFaultError checks if err is a SOAP FaultError and returns it.
func IsFaultError(err error) (*FaultError, bool) {
	if fe, ok := err.(*FaultError); ok {
		return fe, true
	}
	return nil, false
}
