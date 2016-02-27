package main
    
import (
    "fmt"
    "sevrest"
    "flag"
    "os"
    "strings"
    "sort"
)

const (
    // Defaults
    defaultURL = "http://zbrown56c.sevone.com:8080/api/v1/"
    defaultUsername = "admin"
    defaultPassword = ***REMOVED***

    defaultMethod = "get"

    // Help
    helpHelp = "Command line help"

    helpURL = "The server API URL"
    helpUsername = "The server Username"
    helpPassword = "The server Password"
    
    helpList = "List all possible commands from the API"
    helpMethod = "The HTTP method to use (get, post, put, delete)"

    // Other Const
    apiPath = "/api/v1"
)

func main() {

    var url, username, password string
    var help, listCommands bool
    var method string

    flag.BoolVar(&help, "-help", false, helpHelp)
    flag.BoolVar(&help, "h", false, helpHelp + " (shorthand)")

    flag.StringVar(&url, "-api", defaultURL, helpURL)
    flag.StringVar(&url, "a", defaultURL, helpURL + " (shorthand)")
    flag.StringVar(&username, "-username", defaultUsername, helpUsername)
    flag.StringVar(&username, "u", defaultUsername, helpUsername + " (shorthand)")
    flag.StringVar(&password, "-password", defaultPassword, helpPassword)
    flag.StringVar(&password, "p", defaultPassword, helpPassword + " (shorthand)")

    flag.BoolVar(&listCommands, "-list", false, helpList)
    flag.BoolVar(&listCommands, "l", false, helpList + " (shorthand)")

    flag.StringVar(&method, "-method", defaultMethod, helpMethod)
    flag.StringVar(&method, "m", defaultMethod, helpMethod + " (shorthand)")

    // Override build in help
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "%s [options] [commands]: \n", "sevrestcli")
        flag.PrintDefaults()
    }

    // Parse the args
    flag.Parse()

    // Debug
    // fmt.Printf("ARGS: %v %d\n", flag.Args(), len(flag.Args()))

    if(len(flag.Args()) == 0) {
        fmt.Println("Usage: sevrestcli [options] <url>")
        flag.Usage();
        os.Exit(0);
    }

    // Create Client and Login    
    var c = sevrest.Client(url)
    var err = c.Auth(username, password)
    if(err != nil) {
        fmt.Printf("Error authenticating to SevOne. Error: %s\n", err.Error())
        os.Exit(1)
    }

    url = flag.Args()[0]

    // Get context sensitive help
    if(help) {
        //var apiDocs map[string]map[string]map[string]interface{}
        var apiDocs sevrest.SevRestApiDocs

        // Get all de
        resp, err := c.Get("api-docs")
        if(err != nil) {
            fmt.Printf("ERROR: %s", err.Error())
        }
        err = resp.Decode(&apiDocs)

        // sevrest.PrettyPrint(apiDocs)
        // os.Exit(0)
        
        // Sort the endpoints for pretty printing
        var endpoints []string
        for e := range apiDocs.Paths {
            endpoints = append(endpoints, e)
        }
        sort.Strings(endpoints)
        
        // Check to see if they specified a valid endpoint
        _, apiCallPresent := apiDocs.Paths[apiPath+url]

        if(url == "help" || !apiCallPresent) {
            if(url != "help") {
                fmt.Printf("Could not find call: %s searching endpoints.\n", url)
            }
            for _, endpoint := range endpoints {

                // Remove the standard /api/v1 from the beginning as it's assumed for the short endpoint
                shortEndpoint := endpoint;
                if strings.Index(shortEndpoint, apiPath) == 0 {
                    shortEndpoint = shortEndpoint[len(apiPath):]
                }

                // If url isn't help we tried to find something that wasn't there
                // We will search wildcard for that call
                if strings.Index(endpoint, url) == -1 {
                    continue;
                }
                // Print high level help for this call
                fmt.Printf("%s\n", shortEndpoint)
                for method, data := range apiDocs.Paths[endpoint] {
                    fmt.Printf("    %s - %s\n", method, data.Description)
                }
            }
        } else {


            urlHelp := apiDocs.Paths[apiPath+url][method]
            fmt.Printf("URL: %s\n", url)
            fmt.Printf("METHOD: %s\n", method)
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

            // sevrest.PrettyPrint(urlHelp)

            // We have a valid call, print the detailed help on it
            //sevrest.PrettyPrint(urldata)
            // parameters := urldata.(map[string]interface{})["parameters"].([]map[string]interface{})
            // for _, parameter := range parameters {
            //       fmt.Printf("   %s(%s) - %s", parameter["name"], parameter["type"], parameter["description"])
            //  }



            //fmt.Printf("%v", respMap["paths"][url])
        }
        os.Exit(0)
    }


    // Dump the giant API JSON
    //sevrest.PrettyPrint(respMap)

}
