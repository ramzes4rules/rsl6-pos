package soap

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

// ---------- BuildEnvelope tests ----------

func TestBuildEnvelope(t *testing.T) {
	type testReq struct {
		XMLName struct{} `xml:"http://tempuri.org/ TestOp"`
		Value   string   `xml:"http://tempuri.org/ value"`
	}
	data, err := BuildEnvelope("http://tempuri.org/IRSLoyaltyService/TestOp", "https://example.com/service", testReq{Value: "hello"})
	if err != nil {
		t.Fatalf("BuildEnvelope: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, "soap-envelope") {
		t.Error("envelope should contain soap-envelope namespace")
	}
	if !strings.Contains(s, "addressing") {
		t.Error("envelope should contain addressing namespace")
	}
	if !strings.Contains(s, "TestOp") {
		t.Error("envelope should contain operation name")
	}
	if !strings.Contains(s, "hello") {
		t.Error("envelope should contain request value")
	}
	if !strings.Contains(s, "mustUnderstand") {
		t.Error("envelope should contain mustUnderstand attribute")
	}
}

// ---------- ParseEnvelope tests ----------

func soapResp(body string) string {
	return `<?xml version="1.0" encoding="utf-8"?>` +
		`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">` +
		`<s:Body>` + body + `</s:Body></s:Envelope>`
}

func soapFaultResp(code, reason, detail string) string {
	d := ""
	if detail != "" {
		d = `<s:Detail>` + detail + `</s:Detail>`
	}
	return soapResp(`<s:Fault><s:Code><s:Value>` + code + `</s:Value></s:Code>` +
		`<s:Reason><s:Text>` + reason + `</s:Text></s:Reason>` + d + `</s:Fault>`)
}

func TestParseEnvelope_Success(t *testing.T) {
	data := []byte(soapResp(`<TestResponse xmlns="http://tempuri.org/"><TestResult>world</TestResult></TestResponse>`))
	result, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	if result.Fault != nil {
		t.Fatal("expected no fault")
	}
	if len(result.BodyContent) == 0 {
		t.Fatal("expected non-empty body content")
	}

	type testResp struct {
		XMLName struct{} `xml:"TestResponse"`
		Value   string   `xml:"TestResult"`
	}
	var r testResp
	if err := UnmarshalBody(result.BodyContent, &r); err != nil {
		t.Fatalf("UnmarshalBody: %v", err)
	}
	if r.Value != "world" {
		t.Errorf("expected world, got %s", r.Value)
	}
}

func TestParseEnvelope_Fault(t *testing.T) {
	data := []byte(soapFaultResp("s:Receiver", "Test fault", "detail info"))
	result, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	if result.Fault == nil {
		t.Fatal("expected fault")
	}
	if result.Fault.Code != "s:Receiver" {
		t.Errorf("expected code s:Receiver, got %s", result.Fault.Code)
	}
	if result.Fault.Reason != "Test fault" {
		t.Errorf("expected reason 'Test fault', got %s", result.Fault.Reason)
	}
	if result.Fault.Detail != "detail info" {
		t.Errorf("expected detail, got %s", result.Fault.Detail)
	}
}

func TestParseEnvelope_FaultWithSubcode(t *testing.T) {
	data := []byte(soapResp(`<s:Fault><s:Code><s:Value>s:Receiver</s:Value><s:Subcode><s:Value>custom:Code</s:Value></s:Subcode></s:Code><s:Reason><s:Text>err</s:Text></s:Reason></s:Fault>`))
	result, err := ParseEnvelope(data)
	if err != nil {
		t.Fatalf("ParseEnvelope: %v", err)
	}
	if result.Fault == nil {
		t.Fatal("expected fault")
	}
	if result.Fault.Code != "s:Receiver/custom:Code" {
		t.Errorf("expected combined code, got %s", result.Fault.Code)
	}
}

func TestParseEnvelope_InvalidXML(t *testing.T) {
	_, err := ParseEnvelope([]byte("not xml"))
	if err == nil {
		t.Fatal("expected error on invalid XML")
	}
}

func TestUnmarshalBody_Nil(t *testing.T) {
	err := UnmarshalBody([]byte("<X/>"), nil)
	if err != nil {
		t.Fatalf("expected no error for nil result: %v", err)
	}
}

// ---------- XmlDecimal tests ----------

func TestXmlDecimal_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		value    XmlDecimal
		expected string
	}{
		{"zero", 0, "0"},
		{"integer", 100, "100"},
		{"fractional", 125.75, "125.75"},
		{"negative", -50.5, "-50.5"},
		{"large", 1000000.123, "1000000.123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type wrapper struct {
				Val XmlDecimal `xml:"val"`
			}
			data, err := xml.Marshal(wrapper{Val: tt.value})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if !strings.Contains(string(data), tt.expected) {
				t.Errorf("marshaled XML %s does not contain %s", data, tt.expected)
			}
			var w2 wrapper
			if err := xml.Unmarshal(data, &w2); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if float64(w2.Val) != float64(tt.value) {
				t.Errorf("expected %f, got %f", float64(tt.value), float64(w2.Val))
			}
		})
	}
}

func TestXmlDecimal_UnmarshalEmpty(t *testing.T) {
	type wrapper struct {
		Val XmlDecimal `xml:"val"`
	}
	var w wrapper
	if err := xml.Unmarshal([]byte(`<wrapper><val></val></wrapper>`), &w); err != nil {
		t.Fatalf("unmarshal empty: %v", err)
	}
	if float64(w.Val) != 0 {
		t.Errorf("expected 0, got %f", float64(w.Val))
	}
}

// ---------- XmlDateTime tests ----------

func TestXmlDateTime_MarshalUnmarshal(t *testing.T) {
	type wrapper struct {
		Val XmlDateTime `xml:"val"`
	}
	tm := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	data, err := xml.Marshal(wrapper{Val: XmlDateTime(tm)})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), "2024-03-15T10:30:45") {
		t.Errorf("expected datetime string, got %s", data)
	}
	var w2 wrapper
	if err := xml.Unmarshal(data, &w2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	result := time.Time(w2.Val)
	if !result.Equal(tm) {
		t.Errorf("expected %v, got %v", tm, result)
	}
}

func TestXmlDateTime_UnmarshalRFC3339(t *testing.T) {
	type wrapper struct {
		Val XmlDateTime `xml:"val"`
	}
	var w wrapper
	if err := xml.Unmarshal([]byte(`<wrapper><val>2024-03-15T10:30:45Z</val></wrapper>`), &w); err != nil {
		t.Fatalf("unmarshal RFC3339: %v", err)
	}
	result := time.Time(w.Val)
	expected := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// ---------- XmlLongArray tests ----------

func TestXmlLongArray_Marshal(t *testing.T) {
	type wrapper struct {
		Arr XmlLongArray `xml:"arr"`
	}
	data, err := xml.Marshal(wrapper{Arr: XmlLongArray{Items: []int64{1, 2, 3}}})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, "long") {
		t.Error("should contain 'long' elements")
	}
	if !strings.Contains(s, NsArrays) {
		t.Error("should contain arrays namespace")
	}
}

func TestXmlLongArray_MarshalEmpty(t *testing.T) {
	type wrapper struct {
		Arr XmlLongArray `xml:"arr"`
	}
	_, err := xml.Marshal(wrapper{Arr: XmlLongArray{}})
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}
}

