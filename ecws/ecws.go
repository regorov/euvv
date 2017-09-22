// Package ecws implements client to European Comission web services.
package ecws

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const checkVatServiceURL = "http://ec.europa.eu/taxation_customs/vies/services/checkVatService"

var (
	ErrEmptyResponse = errors.New("empty response")
)

// CheckVat describes request attributes.
type CheckVat struct {
	XMLName xml.Name `xml:"urn:ec.europa.eu:taxud:vies:services:checkVat:types checkVat"`

	CountryCode string `xml:"countryCode,omitempty"`

	VatNumber string `xml:"vatNumber,omitempty"`
}

// CheckVatResponse describes response attributes.
type CheckVatResponse struct {
	XMLName xml.Name `xml:"urn:ec.europa.eu:taxud:vies:services:checkVat:types checkVatResponse"`

	CountryCode string `xml:"countryCode,omitempty"`

	VatNumber string `xml:"vatNumber,omitempty"`

	// RequestDate ignored due to special datetime format
	RequestDate time.Time `xml:"-"` // xml:"requestDate,omitempty"

	Valid bool `xml:"valid,omitempty"`

	Name string `xml:"name,omitempty"`

	Address string `xml:"address,omitempty"`
}

// CheckVatService provides interface to VIES VAT number validation service.
type CheckVatService struct {
	client *SOAPClient
}

// NewCheckVatService creates new CheckVatService instance.
func NewCheckVatService(timeout int, verboseMode bool) *CheckVatService {

	client := NewSOAPClient(checkVatServiceURL, time.Duration(timeout)*time.Second, verboseMode)

	return &CheckVatService{
		client: client,
	}
}

// CheckVat sends SOAP request.
func (service *CheckVatService) CheckVat(request *CheckVat) (*CheckVatResponse, error) {
	response := new(CheckVatResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type SOAPClient struct {
	hc          *http.Client
	url         string
	verboseMode bool
	headers     []interface{}
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

// NewSOAPClient returns instantiated SOAPClient.
func NewSOAPClient(url string, timeout time.Duration, verboseMode bool) *SOAPClient {
	return &SOAPClient{
		hc:          &http.Client{Timeout: timeout},
		url:         url,
		verboseMode: verboseMode,
	}
}

// AddHeader addd HTTP header.
func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

// Call sends SOAP envelope over HTTP, receive and process result.
func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &SOAPHeader{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		return errors.New("SOAP message creation failed: " + err.Error())
	}

	if err := encoder.Flush(); err != nil {
		return errors.New("XML encoder flushing failed: " + err.Error())
	}

	if s.verboseMode {
		fmt.Println("Request dump:\n", buffer.String())
	}

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return errors.New("HTTP POST request creation failed: " + err.Error())
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	req.Header.Set("User-Agent", "euvv/0.1")
	req.Close = true

	res, err := s.hc.Do(req)
	if err != nil {
		return errors.New("HTTP request execution failed: " + err.Error())
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("HTTP response body reading failed: " + err.Error())
	}
	if len(rawbody) == 0 {
		return ErrEmptyResponse
	}

	if s.verboseMode {
		fmt.Println("Response dump:\n", string(rawbody))
	}

	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return errors.New("SOAP response unmarshalling failed: " + err.Error())
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return errors.New("Service returned error: " + fault.String)
	}

	return nil
}
