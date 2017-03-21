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

	// TODO:  Create object type
	object_types, err := c.GetObjectTypes(false, nil)
	sevrest.PrettyPrint(object_types)

	// TODO:  Create indicator type
	indicator_types, err := c.GetIndicatorTypes(false, nil)
	sevrest.PrettyPrint(indicator_types)
}

