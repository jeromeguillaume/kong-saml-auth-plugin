/*
	A custom plugin in Go,
	which reads a request header and sets a response header.

	# Copmile code with this command
	go build -o ~/Documents/tmp ./plugins/saml-auth.go
*/

package main

import (
	"fmt"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.1"
var Priority = 1

// Property of the Plugin
type Config struct {
	Message  string
	Message2 string
}

func New() interface{} {
	return &Config{}
}

//-----------------------------------------------------------------------------
// Access => Executed for every request from a Client and before it is being
// proxied to the upstream service (i.e. producer)
//-----------------------------------------------------------------------------
func (conf Config) Access(kong *pdk.PDK) {

	kong.Log.Notice("*** saml-auth - Begin Access() ***")

	// Get Header from the request Consumer
	host, err := kong.Request.GetHeader("host")
	if err != nil {
		kong.Log.Err("Error reading 'host' header: %s", err.Error())
	}

	// Add a header in the request before to call the Backend
	kong.ServiceRequest.AddHeader("x-saml-auth-req", "test")

	// In early stage we add a header in the response sent to the Consumer
	kong.Response.SetHeader("x-saml-auth-res", fmt.Sprintf("Go says %s to %s", conf.Message, host))

	/* In case of issue in this function, we can setup an error (503 - PDK custom error)
	   by calling Exit(503, ...) and the producer is NOT called

	respHeaders := make(map[string][]string)
	respHeaders["Content-Type"] = append(respHeaders["Content-Type"], "text/plain")
	kong.Log.Notice("*** saml-auth - End Access() ***")
	kong.Response.Exit(503, "PDK custom error", respHeaders)
	*/
	kong.Log.Notice("*** saml-auth - End Access() ***")
}

//-----------------------------------------------------------------------------
// Response => Executed after the whole response has been received from the
// upstream service (i.e. Producer), but before sending any part of it to
// the Client
//-----------------------------------------------------------------------------
func (conf Config) Response(kong *pdk.PDK) {
	kong.Log.Notice("*** saml-auth - Begin Response() ***")

	// Add a 2nd header in the response
	kong.Response.SetHeader("x-saml-auth-res2", conf.Message2)

	kong.Log.Notice("*** saml-auth - End Response() ***")
}
