package sevrest

type DeviceData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	OldestTimestamp uint `json:"oldTs"`
	LatestTimestamp uint `json:"newTs"`
	IP string `json:"ip"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	SourceID uint `json:"sourceId"`
	Objects []DeviceDataObject `json:"objects"`
	// Map of object names to indices
	ObjectMap map[string]int
}

type DeviceDataObject struct {
	Name string `json:"name"`
	Type string `json:"type"`
	PluginID uint `json:"pluginId,omitempty"`
	PluginName string `json:"pluginName,omitempty"`
	Description string `json:"description.omitempty"`
	CreateAutomatically bool `json:"automaticCreation,omitempty"`
	Timestamps []DeviceDataTimestamp `json:"timestamps"`
	// Map of times to indices
	TimestampMap map[uint]int
}

type DeviceDataTimestamp struct {
	Time uint `json:"timestamp"`
	Indicators []DeviceDataIndicator `json:"indicators"`
	// Map of indicator names to indices
	IndicatorMap map[string]int
}

type DeviceDataIndicator struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
	// TODO:  This can only be COUNTER32, COUNTER64, or GAUGE; it should
	//    probably be an enum
	Format string `json:"format"`
	Units string `json:"units,omitempty"`
	MaxValue float64 `json:"maxValue,omitempty"`
}

// TODO:  More args?
func NewDeviceData(name string, initial_timestamp uint, source_id uint) DeviceData {
	return DeviceData{
		Name : name,
		Type : "Generic",
		OldestTimestamp : initial_timestamp,
		LatestTimestamp : initial_timestamp,
		SourceID : source_id,
		IP : "0.0.0.0",
		Objects : make([]DeviceDataObject, 0),
		ObjectMap : make(map[string]int),
	}
}

// TODO:  More args?
func (this *DeviceData) NewObject(name string, type_name string) *DeviceDataObject {
	object := DeviceDataObject{
		Name : name,
		Type : type_name,
		PluginID : 17,
		PluginName : "BULKDATA",
		Timestamps : make([]DeviceDataTimestamp, 0),
	}
	this.ObjectMap[name] = len(this.Objects)
	this.Objects = append(this.Objects, object)
	return &object
}

// TODO:  More args?
func (this *DeviceDataObject) NewTimestamp(time uint) *DeviceDataTimestamp {
	timestamp := DeviceDataTimestamp{
		Time : time,
		Indicators : make([]DeviceDataIndicator, 0),
		IndicatorMap : make(map[string]int),
	}
	this.TimestampMap(time) = len(this.Timestamps)
	this.Timestamps = append(this.Timestamps, timestamp)
	return &timestamp
}

// TODO:  More args?
func (this *DeviceDataTimestamp) NewIndicator(name string, value float64) *DeviceDataIndicator {
	indicator := DeviceDataIndicator{
		Name : name,
		Value : value,
		Format : "GAUGE",
	}
	this.IndicatorMap[name] = len(this.Indicators)
	this.Indicators = append(this.Indicators, indicator)
	return &indicator
}

