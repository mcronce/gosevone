package main

import (
    "fmt"
    "sevrest"
    //"strconv"
)

func main() {

    // Create Client and Login
    var c = sevrest.Client("http://bhale-56-10.sevone.com/api/v1/")
    var err = c.Auth("admin", ***REMOVED***)
    if(err != nil) {
        fmt.Printf(err.Error())
    }

    var respMap []map[string]interface{}

    // Get all devicegroups
    fmt.Printf("attempting devicegroups\n")
    resp, err := c.Get("devicegroups")
    if(err != nil) {
        fmt.Printf("ERROR: %s", err.Error())
    }
    fmt.Println("StatusCode: ", resp.StatusCode)
    err = resp.Decode(&respMap)

    sevrest.PrettyPrint(respMap)

}
