package sevrest

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "bytes"
    "io"
    "strconv"
)

const (
    libraryVersion      = "0.1"
    headerUserAgent     = "sevrest/" + libraryVersion
    headerAccept        = "application/json"
    headerContentType   = "application/json;charset=UTF-8"
)

type ClientStruct struct {
    // HTTP client
    client *http.Client

    // Base URL for API requests.
    BaseURL *url.URL

    // The Authentication token
    X_Auth_Token string
}

type Response http.Response

// Build the initial client 
func Client(apiURL string) *ClientStruct {
    baseURL, _ := url.Parse(apiURL)
    c := &ClientStruct{
        client: http.DefaultClient,
        BaseURL: baseURL,
        X_Auth_Token: "",
    }
    return c
}

func (c *ClientStruct) Auth(username string, password string) (error) {

    // Username Password JSON
    authMap := map[string]string { "name": username, "password": password }
    resp, err := c.Post("authentication/signin", authMap)

    if(err != nil || resp.StatusCode != 200) {
        return fmt.Errorf("Unable to log into SevOne. Status %i", resp.StatusCode)
    }

    // We get back a json with just the token
    type Token struct {
        Token string `json: token`
    }

    var t Token
    err = resp.Decode(&t)
    c.X_Auth_Token = t.Token

    return nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash.
func (c *ClientStruct) Request(method string, urlStr string, body io.Reader) (*http.Response, error) {
    
    // Sign into the API
    rel, err := url.Parse(urlStr)
    if err != nil {
        return nil, err
    }

    // Build the URL
    apiUrl := c.BaseURL.ResolveReference(rel)

    req, err := http.NewRequest(method, apiUrl.String(), body)   
    if err != nil {
        return nil, err
    }

    // Headers
    req.Header.Add("Accept", headerAccept)
    req.Header.Add("Content-Type", headerContentType)
    req.Header.Add("User-Agent", headerUserAgent)
    if(c.X_Auth_Token != "") {
        req.Header.Add("X-Auth-Token", c.X_Auth_Token)
    }

    // Do the request
    return c.client.Do(req)

}

func (c *ClientStruct) Get(urlStr string) (*Response, error) { 
    httpresp, err := c.Request("GET", urlStr, nil)
    resp := Response(*httpresp)
    return &resp, err
}

func (c *ClientStruct) Delete(urlStr string) (*Response, error) { 
    httpresp, err := c.Request("DELETE", urlStr, nil)
    resp := Response(*httpresp)
    return &resp, err
}

func (c *ClientStruct) Post(urlStr string, data interface{}) (*Response, error) { 
    JSONReader, err := NewJSONReader(data)
    if(err != nil) {
        return nil, nil
    }
    httpresp, err := c.Request("POST", urlStr, JSONReader)
    resp := Response(*httpresp)
    return &resp, err
}

func (c *ClientStruct) Put(urlStr string, data interface{}) (*Response, error) { 
    JSONReader, err := NewJSONReader(data)
    if(err != nil) {
        return nil, nil
    }
    httpresp, err := c.Request("PUT", urlStr, JSONReader)
    resp := Response(*httpresp)
    return &resp, err
}

func (resp *Response) Decode(target interface{}) (error) {
    return json.NewDecoder(resp.Body).Decode(target)
}

func NewJSONReader(source interface{}) (io.Reader, error) {
    JSONBytes, err := json.Marshal(source)
    if(err != nil) {
        return nil, nil
    }
    JSONReader := bytes.NewReader(JSONBytes)
    return JSONReader, nil
}

func PrettyPrint(x interface{}) {
    b, err := json.MarshalIndent(x, "", "  ")
    if err != nil {
        fmt.Println("error:", err)
    }
    fmt.Print(string(b))
}

func Float64ToIntString(input float64) string {
    return strconv.FormatFloat(input, 'f', 0, 64)
}
