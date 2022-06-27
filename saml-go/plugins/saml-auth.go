/*
	A custom plugin in Go,
	which reads a request header and sets a response header.

	# Copmile code with this command
	go build -o ~/Documents/tmp ./plugins/saml-auth.go
*/

package main

import (
	"fmt"
	"bytes"
	
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.1"
var Priority = 1

// Property of the Plugin
type Config struct {
	Message string
	Fv_oauth_server string
	Message2 string
}

func New() interface{} {
	return &Config{}
}

//-----------------------------------------------------------------------------
// sapGetAccessToken => Get Access Token
//-----------------------------------------------------------------------------
func sapGetAccessToken(	kong *pdk.PDK, 
						conf Config, 
						fv_b64_authorization string, 
						fv_oauth_client string, 
						fv_odata_scope string, 
						fv_b64_assertion string) {
	
	kong.Log.Notice("*** saml-auth - Begin sapGetAccessToken() ***")

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
	
	// Get fv_oauth_server value from Plugin property
	fv_oauth_server := conf.Fv_oauth_server

    req.SetRequestURI("https://" + fv_oauth_server + "/sap/bc/sec/oauth2/token")
    
	// Send a POST
	req.Header.SetMethod("POST")
	
	// Add Headers for SAP
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Charset", "UTF-8")
	req.Header.Set("User-Agent", "KongProxy")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Authorization", "Basic " + fv_b64_authorization)
	
	// Add Form Parameters for SAP
	req.SetBodyString(	"client_id=" 	+ fv_oauth_client 	+ "&" +
						"scope=" 		+ fv_odata_scope	+ "&" +
						"grant_type=urn:ietf:params:oauth:grant-type:saml2-bearer" + "&" +
						"assertion=" + fv_b64_assertion)
	
	Ok := true
	resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)
	
    // Perform the request
    err := fasthttp.Do(req, resp)
    if Ok && err != nil {
        kong.Log.Err(err.Error())
		Ok = false;
    }
    if Ok && resp.StatusCode() != fasthttp.StatusOK {
        kong.Log.Err(fmt.Sprintf("Expected status code %d but got %d", fasthttp.StatusOK, resp.StatusCode()))
        Ok = false;
    }

    // Verify the content type
    contentType := resp.Header.Peek("Content-Type")
    if Ok && bytes.Index(contentType, []byte("application/json")) != 0 {
        kong.Log.Err(fmt.Sprintf("Expected content type application/json but got %s", contentType))
        Ok = false;
    }

    // Do we need to decompress the response?
	if Ok{
		contentEncoding := resp.Header.Peek("Content-Encoding")
		var body []byte
		if bytes.EqualFold(contentEncoding, []byte("gzip")) {
			body, _ = resp.BodyGunzip()
		} else {
			body = resp.Body()
		}
		kong.Log.Notice(fmt.Sprintf("Response body is: %s", body))
	}
	kong.Log.Notice("*** saml-auth - End sapGetAccessToken() ***")
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
	kong.ServiceRequest.AddHeader("x-saml-auth-req","test")

	// token handling
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
	type MyCustomClaims struct {
		Foo string `json:"foo"`
		Iss string `json:"iss"`
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	kong.Log.Notice(fmt.Sprintf("*** saml-auth token=%s ***", token))
	claims, ok := token.Claims.(*MyCustomClaims)
	kong.Log.Notice(fmt.Sprintf("*** saml-auth foo='%v' iss='%v' exp='%v'", claims.Foo, claims.Iss, claims.StandardClaims.ExpiresAt))
	if ok {}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		kong.Log.Notice(fmt.Sprintf("%v %v", claims.Foo, claims.Iss, claims.StandardClaims.ExpiresAt))
	} else {
		kong.Log.Err(err)
	}

	// sap-get-access-token
	sapGetAccessToken(	kong, 
						conf,
						"fv_b64_authorization",
						"fv_oauth_client",
						"fv_odata_scope",
						"fv_b64_assertion")

	// Get message value from Plugin property
	message := conf.Message
	if message == "" {
		message = "hello"
	}

	// In early stage we add a header in the response sent to the Consumer
	kong.Response.SetHeader("x-saml-auth-res", fmt.Sprintf("Go says %s to %s", message, host))

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
	kong.Response.SetHeader("x-saml-auth-res2", fmt.Sprintf("2nd Header"))
	
	kong.Log.Notice("*** saml-auth - End Response() ***")
}