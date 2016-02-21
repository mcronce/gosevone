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
    deviceId := "205";
    objectId := "4613";
    indicatorId := "44555";

    // Start and End time (in milliseconds)
    // 
    end := float64(time.Now().Unix())*1000
    start := end-86400000
    startString := sevrest.Float64ToIntString(start)
    endString := sevrest.Float64ToIntString(end)

    var respSlice []interface{}

    resp, err := c.Get("devices/"+deviceId+"/objects/"+objectId+"/indicators/"+indicatorId+"/data?startTime="+startString+"&endTime="+endString)
    err = resp.Decode(&respSlice)
    sevrest.PrettyPrint(respSlice)

}