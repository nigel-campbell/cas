package cas

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"
)

// Testing
// go test -v attributes.go attribute_test.go
func TestAttributeMarshall(t *testing.T) {
	file, _ := ioutil.ReadFile("./_examples/response_pretty.xml")

	var envelope Envelope

	err := xml.Unmarshal(file, &envelope)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, attribute := range envelope.Body.Response.Assertion.Attributes {
		if attribute.AttributeName == "gtAccountEntitlement" {
			for _, value := range attribute.AttributeValue {
				fmt.Println(value)
			}
		}
	}
}
