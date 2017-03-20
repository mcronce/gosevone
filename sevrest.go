package sevrest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"bytes"
	"io"
)

const (
	libraryVersion      = "0.1"
	headerUserAgent     = "sevrest/" + libraryVersion
	headerAccept        = "application/json"
	headerContentType   = "application/json;charset=UTF-8"
)

// Our client
type ClientStruct struct {
	// HTTP client
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// The Authentication token
	X_Auth_Token string
}

// This type will allow us to create type functions for this library
type Response http.Response

// Build the initial client
func Client(apiURL string) *ClientStruct {
	// Ensure the URL ends with a slash
	if(apiURL[len(apiURL)-1] != '/') {
		apiURL += "/"
	}

	// Setup the client
	baseURL, _ := url.Parse(apiURL)
	client := &ClientStruct{
		client: http.DefaultClient,
		BaseURL: baseURL,
		X_Auth_Token: "",
	}
	return client
}

// Authenticate to the API and store the token for sending in the header
func (c *ClientStruct) Auth(username string, password string) (error) {
	// Username Password JSON
	authMap := map[string]string { "name": username, "password": password }
	resp, err := c.Post("authentication/signin", authMap)
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
	c.X_Auth_Token = t.Token

	return nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash.
func (c *ClientStruct) Request(method string, urlStr string, body io.Reader) (*http.Response, error) {
	// Ensure the Url doesn't start with a slash
	if(urlStr[0] == '/') {
		urlStr = urlStr[1:]
	}

	// Parse the Url
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Build the URL
	apiUrl := c.BaseURL.ResolveReference(rel)

	// Make the request
	req, err := http.NewRequest(method, apiUrl.String(), body)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Add("Accept", headerAccept)
	req.Header.Add("Content-Type", headerContentType)
	req.Header.Add("User-Agent", headerUserAgent)
	// SevOne Auth Token
	if(c.X_Auth_Token != "") {
		req.Header.Add("X-Auth-Token", c.X_Auth_Token)
	}

	// Do the request
	return c.client.Do(req)
}

// GET Request
func (c *ClientStruct) Get(urlStr string) (*Response, error) {
	httpresp, err := c.Request("GET", urlStr, nil)
	if(err != nil) {
		return nil, err
	}
	resp := Response(*httpresp)
	if(err != nil) {
		return nil, err
	}
	return &resp, nil
}

// Get the auth token
func (c *ClientStruct) GetToken() (string) {
	return c.X_Auth_Token
}

// Set the auth token
func (c *ClientStruct) SetToken(token string) {
	c.X_Auth_Token = token
}

// DELETE Request
func (c *ClientStruct) Delete(urlStr string) (*Response, error) {
	httpresp, err := c.Request("DELETE", urlStr, nil)
	if(err != nil) {
		return nil, err
	}
	resp := Response(*httpresp)
	if(err != nil) {
		return nil, err
	}
	return &resp, nil
}

// POST Request
func (c *ClientStruct) Post(urlStr string, data interface{}) (*Response, error) {
	// If it's a reader, we'll assume we're already passing good JSON
	// Otherwise we'll hand off to a function to return JSON from about anything
	var JSONReader io.Reader
	var err error
	switch data := data.(type) {
	case io.Reader:
		JSONReader = data
	default:
		JSONReader, err = NewJSONReader(data)
		if(err != nil) {
			return nil, err
		}
	}
	httpresp, err := c.Request("POST", urlStr, JSONReader)
	if(err != nil) {
		return nil, err
	}
	resp := Response(*httpresp)
	return &resp, nil
}

// PUT Request
func (c *ClientStruct) Put(urlStr string, data interface{}) (*Response, error) {
	// If it's a reader, we'll assume we're already passing good JSON
	// Otherwise we'll hand off to a function to return JSON from about anything
	var JSONReader io.Reader
	var err error
	switch data := data.(type) {
	case io.Reader:
		JSONReader = data
	default:
		JSONReader, err = NewJSONReader(data)
		if(err != nil) {
			return nil, err
		}
	}
	httpresp, err := c.Request("PUT", urlStr, JSONReader)
	if(err != nil) {
		return nil, err
	}
	resp := Response(*httpresp)
	return &resp, nil
}

// This will decode the return JSON into whatever you provide as a container
func (resp *Response) Decode(target interface{}) (error) {
	return json.NewDecoder(resp.Body).Decode(target)
}

// This reads from your source container and provides a Reader for the request
func NewJSONReader(source interface{}) (io.Reader, error) {
	JSONBytes, err := json.Marshal(source)
	if(err != nil) {
		return nil, nil
	}
	JSONReader := bytes.NewReader(JSONBytes)
	return JSONReader, nil
}

// Turns most objects into JSON and prints them pretty
func PrettyPrint(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}

