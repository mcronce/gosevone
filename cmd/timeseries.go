package main
    
import (
    "fmt"
    "sevrest"
    "time"
)

func main() {
    
    var c = sevrest.Client("http://zbrown56c.sevone.com/api/v1/")

    var err = c.Auth("admin", ***REMOVED***)
    if(err != nil) {
        fmt.Printf(err.Error())
    }

    // Device, Object, Indicator
    deviceId := "226";
    objectId := "5041";
    indicatorId := "71024";

    // Start and End time (in milliseconds)
    // 
    end := float64(time.Now().Unix())*1000
    start := end-86400000
    startString := sevrest.FloatToString(start, 0)
    endString := sevrest.FloatToString(end, 0)
    respMap, err := c.Get("/devices/"+deviceId+"/objects/"+objectId+"/indicators/"+indicatorId+"/data?startTime="+startString+"&endTime="+endString)
    sevrest.PrettyPrint(respMap)

}