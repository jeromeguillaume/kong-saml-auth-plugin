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
	Message string
}

func New() interface{} {
	return &Config{}
}

//-----------------------------------------------------------------------------
// Access => Executed for every request from a client and before it is being 
// proxied to the upstream service
//-----------------------------------------------------------------------------
func (conf Config) Access(kong *pdk.PDK) {
	
	kong.Log.Notice("*** saml-auth - Begin Access() ***")
	
	// Get Header from the request Consumer
	host, err := kong.Request.GetHeader("host")
	if err != nil {
		kong.Log.Err("Error reading 'host' header: %s", err.Error())	
	}

	// Add a header in the request before to call the Backend
	kong.ServiceRequest.AddHeader("x-saml-auth-req","test")

	// Get message value from Plugin property
	message := conf.Message
	if message == "" {
		message = "hello"
	}
	
	// Prepare the response sent to the Consumer: add a header in the response 
	kong.Response.SetHeader("x-saml-auth-res", fmt.Sprintf("Go says %s to %s", message, host))

	kong.Log.Notice("*** saml-auth - End Access() ***")
}

//-----------------------------------------------------------------------------
// Response => Executed after the whole response has been received from the 
// upstream service, but before sending any part of it to the client
//-----------------------------------------------------------------------------
func (conf Config) Response(kong *pdk.PDK) {
	kong.Log.Notice("*** saml-auth - Begin Response() ***")
	
	// Add a 2nd header in the response 
	kong.Response.SetHeader("x-saml-auth-res2", fmt.Sprintf("2nd Header"))
	
	/*
	mapHeaders, err := kong.Response.GetHeaders(-1)
	if err != nil {
		kong.Log.Err("Error reading 'GetHeaders': %s", err.Error())	
	}
	kong.Response.Exit(599, "PDK custom error", mapHeaders)
	*/

	kong.Log.Notice("*** saml-auth - End Response() ***")
}