// Package rslpos provides a Go client library for the RSLoyaltyService WCF SOAP service.
//
// # Architecture
//
// The library follows Clean Architecture principles:
//
//   - Domain layer (models.go, service.go, errors.go): pure domain types, the Service
//     interface (port), and domain error types. No dependencies on infrastructure.
//
//   - Infrastructure layer (internal/soap/): SOAP 1.2 envelope construction, XML
//     serialization types, and WCF-specific protocol details. Hidden behind Go's
//     internal package convention — not accessible to consumers.
//
//   - Client layer (client.go, operations.go): the Client struct that implements
//     the Service interface by bridging domain calls to the SOAP infrastructure.
//
//   - Mock layer (mock/): testify-based mock of the Service interface for use
//     in consumer unit tests.
//
// # Usage
//
// Create a client with [NewClient] and call any of the 36 service operations:
//
//	client := rslpos.NewClient("https://server/RS.Loyalty.Service/RSLoyaltyService.svc")
//	online, err := client.IsOnline(ctx, "2.0")
//
// Use [IsFaultError] to detect SOAP faults:
//
//	if fe, ok := rslpos.IsFaultError(err); ok {
//	    log.Printf("SOAP Fault: %s", fe.Reason)
//	}
//
// For testing, use the mock:
//
//	m := new(rslposmock.MockService)
//	m.On("IsCardValid", mock.Anything, "CARD001").Return(true, nil)
package client
