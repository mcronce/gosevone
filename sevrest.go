package sevrest

import (
	"encoding/json"
	"fmt"

	"github.com/mcronce/gorest"
)

const (
	libraryVersion      = "0.1"
	headerUserAgent     = "sevrest/" + libraryVersion
	headerAccept        = "application/json"
	headerContentType   = "application/json;charset=UTF-8"
)

// Our client
type SevRest struct {
	// The low-level REST client that we're wrapping
	Rest *gorest.Client
}

// Build the initial client
func New(base_url string) *SevRest {
	return &SevRest{
		Rest : gorest.New(base_url, map[string]string{
			"User-Agent" : headerUserAgent,
			"Content-Type" : headerContentType,
			"Accept" : headerAccept,
		}),
	}
}

// Authenticate to the API and store the token for sending in the header
func (this *SevRest) Auth(username string, password string) (error) {
	// Username Password JSON
	auth_map := map[string]string {"name" : username, "password" : password}
	resp, err := this.Rest.Post("authentication/signin", auth_map)
	if(err != nil) {
		return err
	}
	if(resp.StatusCode != 200) {
		return fmt.Errorf("Unable to log into SevOne. Status: %d", resp.StatusCode)
	}

	// We get back a json with just the token
	type Token struct {
		Token string `json: token`
	}

	// Decode and store the auth token to use for future requests
	var t Token
	err = resp.Decode(&t)
	this.Rest.Headers["X-Auth-Token"] = t.Token

	return nil
}

// Set the auth token
func (this *SevRest) SetToken(token string) {
	this.Rest.Headers["X-Auth-Token"] = token
}

// Turns most objects into JSON and prints them pretty
func PrettyPrint(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

