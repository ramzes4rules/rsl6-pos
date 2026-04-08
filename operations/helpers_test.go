package operations

import "github.com/ramzes4rules/rsl6-pos/client"

// newClientForURL creates a client.Service backed by a real SOAP client for the given URL.
// Used in tests to wire up httptest servers.
func newClientForURL(url string) client.Service {
	return client.NewClient(url)
}
