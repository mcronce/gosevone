package main
    
import (
    "fmt"
    "sevrest"
    // "flag"
)

func main() {

    // Create Client and Login    
    var c = sevrest.Client("http://zbrown56c.sevone.com:8080/api/v1/")
    // var err = c.Auth("admin", ***REMOVED***)
    // if(err != nil) {
    //     fmt.Printf(err.Error())
    // }
    // 
    
    var respMap map[string]interface{}
   
    // Get all devices
    resp, err := c.Get("api-docs")
    if(err != nil) {
        fmt.Printf("ERROR: %s", err.Error())
    }
    fmt.Println("StatusCode: ", resp.StatusCode)
    err = resp.Decode(&respMap)

    for uri, methods := range respMap["paths"].(map[string]interface{}) {
        fmt.Printf("%s\n", uri)
        for method, data := range methods.(map[string]interface{}) {
            dataX, _ := data.(map[string]interface{})
            fmt.Printf("    %s - %s\n", method, dataX["description"])
        }
        
    }

    // Dump the giant API JSON
    sevrest.PrettyPrint(respMap)

}
