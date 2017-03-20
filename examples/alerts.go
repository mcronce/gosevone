package main

import (
	"fmt"
	"github.com/mcronce/sevrest"
)

func main() {
	// Create Client and Login
	var c = sevrest.Client("http://localhost:8080/api/v1")
	var err = c.Auth("admin", "yourpassword")
	if err != nil {
		fmt.Printf(err.Error())
	}

	// Create a new Alert
	alert := map[string]string{
		"message":       "Test Alert",
		"deviceId":      "226",
		"origin":        "system",
		"closed":        "0",
		"severity":      "3",
		"startTime":     "1429654512",
		"endTime":       "1429654512",
		"pluginName":    "KANKEI_NAI",
		"lastProcessed": "1455374718",
		"objectId":      "5029",
		"ignoreUntil":   "0",
	}

	// The response will just be a map of string
	var respMap map[string]interface{}

	resp, err := c.Post("policies/61/alerts", alert)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
	fmt.Println("StatusCode: ", resp.StatusCode)
	err = resp.Decode(&respMap)
	sevrest.PrettyPrint(respMap)
}
