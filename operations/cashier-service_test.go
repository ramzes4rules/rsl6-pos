package operations

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// soapOKResponse is a minimal valid SOAP 1.2 envelope with an empty body.
const soapOKResponse = `<?xml version="1.0" encoding="utf-8"?>` +
	`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">` +
	`<s:Body></s:Body>` +
	`</s:Envelope>`

// soapFaultResponse is a SOAP 1.2 fault envelope.
const soapFaultResponse = `<?xml version="1.0" encoding="utf-8"?>` +
	`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">` +
	`<s:Body>` +
	`<s:Fault><s:Code><s:Value>s:Receiver</s:Value></s:Code>` +
	`<s:Reason><s:Text>Service unavailable</s:Text></s:Reason>` +
	`</s:Fault>` +
	`</s:Body>` +
	`</s:Envelope>`

// newTestService creates a CashierService backed by a fake SOAP server
// that responds with the given body for every request.
func newTestService(t *testing.T, responseBody string) (*CashierService, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, responseBody)
	}))
	svc := newServiceFromURL(t, srv.URL)
	return svc, srv.Close
}

// newServiceFromURL creates a CashierService pointing at the given URL.
func newServiceFromURL(t *testing.T, url string) *CashierService {
	t.Helper()
	c := newClientForURL(url)
	return NewCashierService(c, nil)
}

func Test_Ping_Success(t *testing.T) {
	svc, cleanup := newTestService(t, soapOKResponse)
	defer cleanup()

	if err := svc.Ping(context.Background()); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func Test_Ping_Fault(t *testing.T) {
	svc, cleanup := newTestService(t, soapFaultResponse)
	defer cleanup()

	err := svc.Ping(context.Background())
	if err == nil {
		t.Fatal("expected an error for SOAP fault, got nil")
	}
}

func Test_Ping_ConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	srv.Close() // close immediately

	svc := newServiceFromURL(t, srv.URL)
	err := svc.Ping(context.Background())
	if err == nil {
		t.Fatal("expected an error for closed server, got nil")
	}
}
