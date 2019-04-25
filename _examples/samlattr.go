package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/nigel-campbell/cas"
	"net/http"
	"net/url"
	"text/template"
)

type myHandler struct{}

var MyHandler = &myHandler{}
var casURL string

func init() {
	flag.StringVar(&casURL, "url", "", "CAS server URL")
}

const SAML_ASSERTION_TEMPLATE = `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
<SOAP-ENV:Header/>
<SOAP-ENV:Body>
<samlp:Request xmlns:samlp="urn:oasis:names:tc:SAML:1.0:protocol"
MajorVersion="1"
MinorVersion="1"
RequestID="{{.RequestId}}"
IssueInstant="{{.Timestamp}}">
<samlp:AssertionArtifact>{{.Ticket}}</samlp:AssertionArtifact></samlp:Request>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

func fetch_saml_validation(serviceUrl string, ticket string) string {
	return "true"
}

func newRequestId() string {
	const alphabet = "abcdef0123456789"

	// generate 64 character string
	bytes := make([]byte, 64)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = alphabet[v%byte(len(alphabet))]
	}

	return string(bytes)
}

func main() {
	Example()
}

func Example() {
	flag.Parse()

	if casURL == "" {
		flag.Usage()
		return
	}

	glog.Info("Starting up")

	m := http.NewServeMux()
	m.Handle("/", MyHandler)

	url, _ := url.Parse(casURL)
	client := cas.NewClient(&cas.Options{
		URL:        url,
		CasVersion: "CAS_2_SAML_1_0",
	})

	server := &http.Server{
		Addr:    ":7979",
		Handler: client.Handle(m),
	}

	if err := server.ListenAndServe(); err != nil {
		glog.Infof("Error from HTTP Server: %v", err)
	}

	glog.Info("Shutting down")
}

type templateBinding struct {
	Username   string
	Attributes cas.UserAttributes
	Ticket     string
}

type samlBinding struct {
	RequestId string
	Timestamp string
	Ticket    string
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !cas.IsAuthenticated(r) {
		cas.RedirectToLogin(w, r)
		return
	}

	if r.URL.Path == "/logout" {
		cas.RedirectToLogout(w, r)
		return
	}

	envelope := cas.MarshalledResponse(r)

	w.Header().Add("Content-Type", "text/html")

	tmpl, err := template.New("index.html").Parse(index_html)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	fmt.Printf("Envelope %v\n", envelope)
	binding := &templateBinding{
		Username:   cas.Username(r),
		Attributes: cas.Attributes(r),
		// Ticket:     cas.TicketStore.Read(cas.Username(r)),
	}

	html := new(bytes.Buffer)
	if err := tmpl.Execute(html, binding); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, error_500, err)
		return
	}

	//query := r.URL.Query()
	ticket := query.Get("ticket")
	//serviceUrl := query.Get("service")
	fmt.Println("Response*************\n: ", r)
	//fmt.Println("URL: \n", r.URL.RawPath)
	fmt.Printf("Req: %s%s\n", r.Host, r.URL.Path)
	fmt.Printf("Ticket %s: ", ticket)

	html.WriteTo(w)

	//samlBinding := &samlBinding{
	//	RequestId: newRequestId(),
	//	Timestamp: time.Now().Format(time.RFC3339),
	//	Ticket: query.Get("ticket"),
	//}
	//
	//t := template.Must(template.New("saml").Parse(SAML_ASSERTION_TEMPLATE))
	//buf := new(bytes.Buffer)
	//
	//err = t.Execute(os.Stdout, samlBinding)
	//err = t.Execute(buf, samlBinding)
	//if err != nil {
	//	log.Println("executing template:", err)
	//}
	////const target = "https://cop-dev.police.gatech.edu/gocas"
	//const myurl = "https://login.gatech.edu/cas/samlValidate?TARGET=https://cop-dev.police.gatech.edu/gocas"
	//fmt.Println(buf.String())
	//resp, err := http.Post(myurl, "text/xml", strings.NewReader(buf.String()))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//fmt.Println("SAML Response Status: ", resp.StatusCode)
	//if resp.StatusCode == http.StatusOK {
	//	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//	bodyString := string(bodyBytes)
	//	fmt.Println("SAML RESPONSE: ", bodyString)
	//}
}

const index_html = `<!DOCTYPE html>
<html>
  <head>
    <title>Welcome {{.Username}}</title>
  </head>
  <body>
    <h1>Welcome {{.Username}} <a href="/logout">Logout</a></h1>
	<p>Your ticket is: {{.Ticket}}</p>
    <p>Your attributes are:</p>
    <ul>{{range $key, $values := .Attributes}}
      <li>{{$len := len $values}}{{$key}}:{{if gt $len 1}}
        <ul>{{range $values}}
          <li>{{.}}</li>{{end}}
        </ul>
      {{else}} {{index $values 0}}{{end}}</li>{{end}}
    </ul>
  </body>
</html>
`

const error_500 = `<!DOCTYPE html>
<html>
  <head>
    <title>Error 500</title>
  </head>
  <body>
    <h1>Error 500</h1>
    <p>%v</p>
  </body>
</html>
`
