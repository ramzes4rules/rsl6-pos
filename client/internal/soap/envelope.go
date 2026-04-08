// Package soap provides low-level SOAP 1.2 envelope construction and parsing
// for WCF-compatible services. This is an internal package and should not be
// imported directly by consumers of the library.
package soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// XML Namespaces used in SOAP messages.
const (
	NsSoap12   = "http://www.w3.org/2003/05/soap-envelope"
	NsAddr     = "http://www.w3.org/2005/08/addressing"
	NsTempuri  = "http://tempuri.org/"
	NsProtocol = "http://schemas.datacontract.org/2004/07/RS.Loyalty.Domain.Protocol"
	NsArrays   = "http://schemas.microsoft.com/2003/10/Serialization/Arrays"
	NsSystem   = "http://schemas.datacontract.org/2004/07/System"

	ActionPrefix = "http://tempuri.org/IRSLoyaltyService/"
)

// BuildEnvelope constructs a SOAP 1.2 envelope with WS-Addressing headers.
func BuildEnvelope(action, to string, body interface{}) ([]byte, error) {
	bodyBytes, err := xml.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	buf.WriteString(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing">`)
	buf.WriteString(`<s:Header>`)
	buf.WriteString(`<a:Action s:mustUnderstand="1">`)
	xml.EscapeText(&buf, []byte(action))
	buf.WriteString(`</a:Action>`)
	buf.WriteString(`<a:To s:mustUnderstand="1">`)
	xml.EscapeText(&buf, []byte(to))
	buf.WriteString(`</a:To>`)
	buf.WriteString(`</s:Header>`)
	buf.WriteString(`<s:Body>`)
	buf.Write(bodyBytes)
	buf.WriteString(`</s:Body>`)
	buf.WriteString(`</s:Envelope>`)

	return buf.Bytes(), nil
}

// RawFault holds parsed SOAP fault data without coupling to domain error types.
type RawFault struct {
	Code   string
	Reason string
	Detail string
}

// ParseResult holds the outcome of parsing a SOAP response envelope.
type ParseResult struct {
	Fault       *RawFault
	BodyContent []byte
}

// SOAP response envelope types for unmarshaling.
type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    body     `xml:"Body"`
}

type body struct {
	Fault   *fault `xml:"Fault"`
	Content []byte `xml:",innerxml"`
}

type fault struct {
	Code   faultCode   `xml:"Code"`
	Reason faultReason `xml:"Reason"`
	Detail string      `xml:"Detail"`
}

type faultCode struct {
	Value   string     `xml:"Value"`
	Subcode *faultCode `xml:"Subcode"`
}

type faultReason struct {
	Text string `xml:"Text"`
}

// ParseEnvelope parses a SOAP 1.2 response envelope and returns the result.
// Callers should check ParseResult.Fault to handle SOAP faults.
func ParseEnvelope(data []byte) (*ParseResult, error) {
	var env envelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("unmarshal soap envelope: %w", err)
	}
	if env.Body.Fault != nil {
		f := env.Body.Fault
		code := f.Code.Value
		if f.Code.Subcode != nil {
			code += "/" + f.Code.Subcode.Value
		}
		return &ParseResult{
			Fault: &RawFault{
				Code:   code,
				Reason: f.Reason.Text,
				Detail: f.Detail,
			},
		}, nil
	}
	return &ParseResult{
		BodyContent: env.Body.Content,
	}, nil
}

// UnmarshalBody is a convenience to unmarshal the body content into a result struct.
func UnmarshalBody(bodyContent []byte, result interface{}) error {
	if result == nil {
		return nil
	}
	if err := xml.Unmarshal(bodyContent, result); err != nil {
		return fmt.Errorf("unmarshal soap body: %w", err)
	}
	return nil
}
