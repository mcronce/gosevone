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
    device := map[string]string {
        "name":             "Test Device Y",
        "description":      "Test Description 1",
        "ipAddress":        "127.0.0.1",
        "pollingFrequency": "300",
    }

    var respMap map[string]interface{}

    // We get generic interface back from these calls
    resp, err := c.Post("devices", device)
    err = resp.Decode(&respMap)
    sevrest.PrettyPrint(respMap)

    // Get my ID (and numbers in the map are float64)
    deviceId := sevrest.Float64ToIntString(respMap["id"].(float64))

    // Get our device
    resp, err = c.Get("devices/"+deviceId)
    err = resp.Decode(&respMap)
    sevrest.PrettyPrint(respMap)
    
    // Get all devices
    resp, err = c.Get("devices")
    err = resp.Decode(&respMap)
    sevrest.PrettyPrint(respMap)

}