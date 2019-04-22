package cas

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

type Envelope struct {
	XMLName xml.Name // `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    *SoapBody
}

type SoapBody struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Response SoapResponse
}

type SoapResponse struct {
	XMLName xml.Name `xml:Response"`
	Status  SoapStatus
}

type SoapStatus struct {
	XMLName    xml.Name `xml:Status"`
	StatusCode SoapStatusCode
}

type SoapStatusCode struct {
	XMLName xml.Name `xml:StatusCode"`
	Value   string   `xml:"Value,attr"`
}

func TestAttributeMarshall(t *testing.T) {
	file, _ := ioutil.ReadFile("./_examples/response_pretty.xml")

	var envelope Envelope

	err := xml.Unmarshal(file, &envelope)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n", envelope.Body)
}
