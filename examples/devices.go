package main

import (
	"fmt"
	"github.com/sevone/sevrest"
)

func main() {
	// Create Client and Login
	var c = sevrest.New("http://localhost:8080/api/v1")
	var err = c.Auth("admin", "yourpassword")
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Data for creating a new device
	device := map[string]string{
		"name":             "Test Device O",
		"description":      "Test Description 1",
		"ipAddress":        "127.0.0.1",
		"pollingFrequency": "300",
	}

	// We will get a full device JSON back, but we only care to parse the deviceId in the response
	type CreateDeviceResponse struct {
		DeviceId int `json:"id"`
	}
	var respDevice CreateDeviceResponse

	// Create the device
	resp, err := c.Rest.Post("devices", device)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
	fmt.Println("StatusCode: ", resp.StatusCode)
	err = resp.Decode(&respDevice)
	sevrest.PrettyPrint(respDevice) // We will only see the id Field

	// Get our newly created device
	devices, err := c.GetDevices(map[string]interface{}{"ids" : []int{respDevice.DeviceId}})
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
	sevrest.PrettyPrint(devices)
	fmt.Println("---")

	// Get all devices
	devices, err = c.GetDevices(nil)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
	sevrest.PrettyPrint(devices)
}

