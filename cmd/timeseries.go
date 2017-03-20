package main

import (
	"fmt"
	"time"
	"strconv"
	"github.com/mcronce/sevrest"
)

// We can built a custom struct to receive data and type it accordingly
type TimeSeriesValue struct {
	Focus string  `json:"focus,string,omitempty"`
	Time int64    `json:time`
	Value float64 `json:value`
}

func main() {
	// Create Client and Login
	var c = sevrest.Client("http://zbrown56c.sevone.com/api/v1")
	var err = c.Auth("admin", ***REMOVED***)
	if(err != nil) {
		fmt.Printf(err.Error())
	}

	// Device, Object, Indicator
	deviceId := "205";
	objectId := "4613";
	indicatorId := "44555";

	// Start and End time (in milliseconds)
	end := float64(time.Now().Unix())*1000
	start := end-86400000 // 24 Hours Previous
	startString := strconv.FormatFloat(start, 'f', 0, 64)
	endString := strconv.FormatFloat(end, 'f', 0, 64)

	// The return will be a slice of our struct
	var respSlice []TimeSeriesValue

	// Do the request
	resp, err := c.Get("/devices/"+deviceId+"/objects/"+objectId+"/indicators/"+indicatorId+"/data?startTime="+startString+"&endTime="+endString)
	if(err != nil) {
		fmt.Printf("ERROR: %s", err.Error())
	}
	fmt.Println("StatusCode: ", resp.StatusCode)
	err = resp.Decode(&respSlice)
	sevrest.PrettyPrint(respSlice)
}

