package main
    
import (
    "fmt"
    "sevrest"
)

func main() {
    
    var c = sevrest.Client("http://zbrown56c.sevone.com/api/v1/")

    var err = c.Auth("admin", ***REMOVED***)
    if(err != nil) {
        fmt.Printf(err.Error())
    }

    // Create a new Device
    alert := map[string]string {
        "message"       : "Test Alert",
        "deviceId"      : "226",
        "origin"        : "system",
        "closed"        : "0",
        "severity"      : "3",
        "startTime"     : "1429654512",
        "endTime"       : "1429654512",
        "pluginName"    : "KANKEI_NAI",
        "lastProcessed" : "1455374718",
        "objectId"      : "5029",
        "ignoreUntil"   : "0",
    }
    respMap, err := c.Post("/policies/61/alerts", alert)
    sevrest.PrettyPrint(respMap)

}