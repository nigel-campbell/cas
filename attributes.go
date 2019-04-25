package cas

import (
	"encoding/xml"
	"fmt"
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
	XMLName   xml.Name `xml:Response"`
	Status    SoapStatus
	Assertion SoapAssertion
}

type SoapStatus struct {
	XMLName    xml.Name `xml:Status"`
	StatusCode SoapStatusCode
}

type SoapStatusCode struct {
	XMLName xml.Name `xml:StatusCode"`
	Value   string   `xml:"Value,attr"`
}

type SoapAssertion struct {
	XMLName                 xml.Name `xml:Assertion`
	Conditions              SoapConditions
	AuthenticationStatement SoapAuthenticationStatement
	Attributes              []SAMLAttribute `xml:"AttributeStatement>Attribute"`
}

type SoapConditions struct {
	XMLName xml.Name `xml:Conditions`
}

type SoapAuthenticationStatement struct {
	XMLName xml.Name `xml:AuthenticationStatement`
}

type SoapAttributeStatement struct {
	XMLName    xml.Name        `xml:AttributeStatement`
	Attributes []SAMLAttribute `xml:">Attribute"`
}

type SAMLSubject struct {
	XMLName        xml.Name `xml.Subject`
	NameIdentifier string   `xml:NameIdentifier`
}

type SAMLAttribute struct {
	XMLName        xml.Name `xml:Attribute`
	AttributeName  string   `xml:"AttributeName,attr"`
	AttributeValue []string `xml:AttributeValue`
}

func (envelope Envelope) printAttributes() {
	for _, attribute := range envelope.Body.Response.Assertion.Attributes {
		for _, value := range attribute.AttributeValue {
			fmt.Println(value)
		}
	}
}

func (envelope Envelope) attributeExists(attributeValue string) bool {
	for _, attribute := range envelope.Body.Response.Assertion.Attributes {
		for _, value := range attribute.AttributeValue {
			if value == attributeValue {
				return true
			}
		}
	}
	return false
}
