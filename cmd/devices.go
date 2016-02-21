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
        "name":             "Test Device 1",
        "description":      "Test Description 1",
        "ipAddress":        "127.0.0.1",
        "pollingFrequency": "300",
    }
    respMap, err := c.Post("devices", device)
    sevrest.PrettyPrint(respMap)

    // Get my ID (and numbers in the map are float64)
    deviceId := sevrest.FloatToString(respMap["id"].(float64), 0)

    // Get our device
    respMap, err = c.Get("devices/"+deviceId)
    sevrest.PrettyPrint(respMap)
    
    // Get all devices
    respMap, err = c.Get("devices")
    sevrest.PrettyPrint(respMap)

}