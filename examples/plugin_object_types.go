package main

import (
	"fmt"
	"github.com/sevone/sevrest"
)

func main() {
	// Create Client and Login
	var c = sevrest.New("http://10.129.12.2/api/v1/")
	var err = c.Auth("admin", "SevOne")
	if(err != nil) {
		fmt.Printf(err.Error())
	}

	// TODO:  Create type
	response, err := c.GetPluginObjectTypes(false, nil)
	sevrest.PrettyPrint(response)
}

