package main

import (
	"fmt"
	//"strconv"
	"github.com/mcronce/sevrest"
)

func main() {
	// Create Client and Login
	var c = sevrest.Client("http://localhost:8080/api/v1/")
	var err = c.Auth("admin", "yourpassword")
	if err != nil {
		fmt.Printf(err.Error())
	}

	var respMap []map[string]interface{}

	// Get all devicegroups
	fmt.Printf("attempting devicegroups\n")
	resp, err := c.Get("devicegroups")
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
	fmt.Println("StatusCode: ", resp.StatusCode)
	err = resp.Decode(&respMap)

	sevrest.PrettyPrint(respMap)
}
