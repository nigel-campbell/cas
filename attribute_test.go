package cas

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Soap    *SoapBody
}

type SoapBody struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Response SoapResponse
}

type SoapResponse struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:1.0:protocol/ Response"`
	Status  SoapStatus
}

type SoapStatus struct {
	XMLName    xml.Name `xml:"urn:oasis:names:tc:SAML:1.0:protocol/ Status"`
	StatusCode SoapStatusCode
}

type SoapStatusCode struct {
	XMLName         xml.Name `xml:"urn:oasis:names:tc:SAML:1.0:protocol/ StatusCode"`
	StatusCodeValue string
}

func TestAttributeMarshall(t *testing.T) {
	file, _ := ioutil.ReadFile("./_examples/response_pretty.xml")

	var envelope Envelope

	err := xml.Unmarshal(file, &envelope)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n", envelope.Soap)
}
