package main
    
import (
    "fmt"
    "sevrest"
    "flag"
    "os"
    "io"
    "strings"
    "sort"
    "bufio"
)

const (
    // Defaults
    defaultURL = "http://localhost:8080/api/v1"
    defaultUsername = "admin"
    defaultPassword = ***REMOVED***

    defaultMethod = "get"

    // Help
    helpHelp = "Command line help"

    helpURL = "The server API URL"
    helpUsername = "The server Username"
    helpPassword = "The server Password"
    helpToken = "A valid auth token for this server"
    helpGetToken = "Will authenticate and return you a valid token"
    
    helpList = "List all possible commands from the API"
    helpMethod = "The HTTP method to use (get, post, put, delete)"
    helpJSON = "JSON File for input (for put, post requeusts)"
    helpJSONstdin = "JSON will be provided via stdin"

    // Other Const
    apiPath = "/api/v1"
)

func main() {

    var apiUrl, username, password, token string
    var help, jsonStdin bool
    var method, jsonFile string
    var err error

    // Flags
    flag.BoolVar(&help, "help", false, helpHelp)
    flag.BoolVar(&help, "h", false, helpHelp + " (shorthand)")

    flag.StringVar(&apiUrl, "api", defaultURL, helpURL)
    flag.StringVar(&apiUrl, "a", defaultURL, helpURL + " (shorthand)")
    flag.StringVar(&username, "username", defaultUsername, helpUsername)
    flag.StringVar(&username, "u", defaultUsername, helpUsername + " (shorthand)")
    flag.StringVar(&password, "password", defaultPassword, helpPassword)
    flag.StringVar(&password, "p", defaultPassword, helpPassword + " (shorthand)")
    flag.StringVar(&token, "token", "", helpToken)
    flag.StringVar(&token, "t", "", helpToken + " (shorthand)")

    flag.StringVar(&method, "method", defaultMethod, helpMethod)
    flag.StringVar(&method, "m", defaultMethod, helpMethod + " (shorthand)")

    flag.StringVar(&jsonFile, "-json", "", helpJSON)
    flag.StringVar(&jsonFile, "j", "", helpJSON + " (shorthand)")
    flag.BoolVar(&jsonStdin, "json-stdin", false, helpJSONstdin)
    flag.BoolVar(&jsonStdin, "s", false, helpJSONstdin + " (shorthand)")

    // Override built in help
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "%s [options] [commands]: \n", "sevrestcli")
        flag.PrintDefaults()
    }

    // Parse the args
    flag.Parse()
    var args = flag.Args()

    // Help Messsages
    if(len(args) == 0) {
        fmt.Println("\nUsage: sevrestcli [options] <url>\n")
        flag.Usage();
        fmt.Println("\nAny option without a corresponding flag is considered part of the url")
        fmt.Println("\nIf the unflagged option includes an = sign, it will be interpreted as part of a url option")
        fmt.Println("Examples:")
        fmt.Println("  FLAGS: /devices/1/2/3?test=1 URL: /devices/1/2/3?test=1")
        fmt.Println("  FLAGS: devices 1 2 3 test=1 4 5 test=2 6 URL: /devices/1/2/3/4/5/6?test=1&test=2")
        fmt.Println("Help:")
        fmt.Println("  Run with help option with a url of endpoints to list all possible endpoints")
        fmt.Println("    -h endpoints")
        fmt.Println("  Run with help and a specific endpoint to get detailed help on that endpoint")
        fmt.Println("    -method post -help /authentication/signin")
        fmt.Println("Authentication Token:")
        fmt.Println("  Run endpoint 'token' and it will return an authentication token you can use for subsequent calls")
        fmt.Println("Environment Variables: (command line flags take precedence over env variables)")
        fmt.Println("  SEVONE_API - Set to the API endpoint")
        fmt.Println("  SEVONE_USERNAME - The username to use to connect")
        fmt.Println("  SEVONE_PASSWORD - The password to use to connect")
        fmt.Println("  SEVONE_TOKEN - An authentication token received from the 'token' endpoint")
        fmt.Println("    If you provide a valid token, you don't need a username and password")
        os.Exit(0);
    }

    // Allow setting values using environment variables
    if(apiUrl == defaultURL && os.Getenv("SEVONE_API") != "") {
        apiUrl = os.Getenv("SEVONE_API")
    }
    if(username == defaultUsername && os.Getenv("SEVONE_USERNAME") != "") {
        username = os.Getenv("SEVONE_USERNAME")
    }
    if(password == defaultURL && os.Getenv("SEVONE_PASSWORD") != "") {
        password = os.Getenv("SEVONE_PASSWORD")
    }
    if(token == "" && os.Getenv("SEVONE_TOKEN") != "") {
        token = os.Getenv("SEVONE_TOKEN")
    }

    // Build the URL and options string
    var urlSlice, optionsSlice []string
    var optionsIndex int
    
    // The first option should be the URL if they decided to put it in /blah/blah/blah?option=1&option=2 format
    optionsIndex = strings.Index(args[0],"?")
    if(optionsIndex > -1) {
        urlSlice = strings.Split(args[0][:optionsIndex],"/")
        optionsSlice = strings.Split(args[0][optionsIndex+1:],"&")
    } else {
        urlSlice = strings.Split(args[0],"/")
    }

    // We've handled the first arg regardless of format, pop it off the front
    args = args[1:]
    for _, v := range args {
        if(strings.Index(v, "=") == -1) {
            urlSlice = append(urlSlice, v)
        } else {
            optionsSlice = append(optionsSlice, v)
        }
    }

    // This is the URL and options string that can be passed to the API
    url := strings.Join(urlSlice, "/")
    options := strings.Join(optionsSlice, "&")
    fullUrl := url
    if(len(options) > 0) {
        fullUrl += "?"+options
    }

    // Create Client
    var c = sevrest.Client(apiUrl)
    // If we specified the token
    if(token != "") {
        c.SetToken(token)
    // We need to login    
    } else {
        var err = c.Auth(username, password)
        if(err != nil) {
            fmt.Printf("Error authenticating to SevOne. Error: %s\n", err.Error())
            os.Exit(1)
        }
        // Spit out the token assuming we authed okay
        if(url == "token") {
            fmt.Println(c.GetToken())
            os.Exit(0)
        }
    }

    // API HELP
    if(help) {
        var apiDocs sevrest.SevRestApiDocs

        // Get the helpdocs from the api
        resp, err := c.Get("api-docs")
        if(err != nil) {
            fmt.Printf("ERROR: %s", err.Error())
        }
        err = resp.Decode(&apiDocs)
        
        // Use this to debug help structure
        // sevrest.PrettyPrint(apiDocs)
        
        // Sort the endpoints for pretty printing
        var endpoints []string
        for e := range apiDocs.Paths {
            endpoints = append(endpoints, e)
        }
        sort.Strings(endpoints)
        
        // Check to see if they specified a valid endpoint
        _, apiCallPresent := apiDocs.Paths[apiPath+url][method]

        if(url == "endpoints" || !apiCallPresent) {
            if(url != "endpoints") {
                fmt.Printf("Could not find call: (%s):%s searching endpoints.\n", method, url)
            }
            for _, endpoint := range endpoints {

                // Remove the standard /api/v1 from the beginning as it's assumed for the short endpoint
                shortEndpoint := endpoint;
                if strings.Index(shortEndpoint, apiPath) == 0 {
                    shortEndpoint = shortEndpoint[len(apiPath):]
                }

                // If url isn't help we tried to find something that wasn't there
                // We will search wildcard for that call
                if strings.Index(endpoint, url) == -1 && url != "endpoints" {
                    continue;
                }
                // Print high level help for this call
                fmt.Printf("%s\n", shortEndpoint)
                for method, data := range apiDocs.Paths[endpoint] {
                    fmt.Printf("    %s - %s\n", method, data.Description)
                }
            }
        } else {
            // We specified a specific URL, let's search for it and print it's info
            urlHelp := apiDocs.Paths[apiPath+url][method]
            fmt.Printf("URL: %s\n", url)
            fmt.Printf("METHOD: %s\n", method)
            fmt.Printf("DESCRIPTION: %s\n", urlHelp.Description)
            fmt.Printf("PARAMETERS:\n")
            for _, parameter := range urlHelp.Parameters {
                if len(parameter.Schema) != 0 {
                    apiDocs.PrintSchema(parameter.Schema, "        ")
                } else {
                    fmt.Printf("    %s(%s) - %s\n", parameter.Name, parameter.Type, parameter.Description)        
                }
                if parameter.Enum != nil {
                    fmt.Printf("      Valid: %s\n", strings.Join(parameter.Enum, ","))
                }
            }
            fmt.Printf("RESPONSES:\n")
            for statusCode, response := range urlHelp.Responses {
                fmt.Printf("STATUS: %s - %s\n", statusCode, response.Description)
                //sevrest.PrettyPrint(response)
                
                // Check to see if they specified a valid endpoint
                apiDocs.PrintSchema(response.Schema, "    ")
            }
        }
        os.Exit(0)
    }

    // JSON Input
    var reader io.Reader

    // We specified a file
    if(len(jsonFile) != 0) {
        reader, err = os.Open(jsonFile)
        if(err != nil) {
            fmt.Printf("ERROR: %s\n", err.Error())
            os.Exit(1)
        }
    // Stdin
    } else if(jsonStdin) {
        reader = bufio.NewReader(os.Stdin)
    // Nothing
    } else {
        reader, err = os.Open("/dev/null")
        if(err != nil) {
            fmt.Printf("ERROR: %s\n", err.Error())
            os.Exit(1)
        }
    }

    var resp *sevrest.Response
    var respJSON interface{}

    // Make the call
    switch method {
    case "get":
        resp, err = c.Get(fullUrl)
    case "post":
        resp, err = c.Post(fullUrl, reader)
    case "delete":
        resp, err = c.Delete(fullUrl)
    case "put":
        resp, err = c.Put(fullUrl, reader)
    }
    if(err != nil) {
        fmt.Printf("ERROR: %s", err.Error())
        os.Exit(1)
    }
    
    // Decode and print the output
    err = resp.Decode(&respJSON)        
    sevrest.PrettyPrint(respJSON)

    // // If the status code wasn't in the 200's we'll assume something was wrong
    if(resp.StatusCode < 200 || resp.StatusCode >= 300) {
        os.Exit(255)
    }

}
