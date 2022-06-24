/*
	A "hello world" plugin in Go,
	which reads a request header and sets a response header.
*/

package main

import (
	"fmt"
	"log"
	
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.1"
var Priority = 1

type Config struct {
	Message string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	host, err := kong.Request.GetHeader("host")
	if err != nil {
		log.Printf("Error reading 'host' header: %s", err.Error())
	}

	// Add a header in the request before the call of Backend
	kong.ServiceRequest.AddHeader("x-saml-auth-req","test")

	message := conf.Message
	if message == "" {
		message = "hello"
	}
	
	// Add a header in the response sent to the Consumer
	kong.Response.SetHeader("x-saml-auth-res", fmt.Sprintf("Go says %s to %s", message, host))
}