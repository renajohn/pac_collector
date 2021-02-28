package swcsource

import (
	"encoding/xml"
	"log"
	"strconv"
	"strings"
)

// TODO: using index is horrible, use ID
const heatingOutboundTemperature = 0
const heatingInboundTemperature = 1
const outsideTemperature = 4
const tankTemperature = 6
const targetTankTemperature = 7
const drillInboundTemperature = 8
const drillOutboundTemperature = 9
const ambiantIndoorTemperature = 20
const ambiantIndoorTargetTemperature = 21

// SWCMeasurement represents all monitored temperatures out of the SWC heating system
type SWCMeasurement struct {
	HeatingOutboundTemperature     float64
	HeatingInboundTemperature      float64
	OutsideTemperature             float64
	TankTemperature                float64
	TargetTankTemperature          float64
	DrillInboundTemperature        float64
	DrillOutboundTemperature       float64
	AmbiantIndoorTemperature       float64
	AmbiantIndoorTargetTemperature float64
}

type _Values struct {
	XMLName xml.Name `xml:"values"`
	Items   []_Item  `xml:"item"`
}

type _Item struct {
	XMLName xml.Name `xml:"item"`
	Ref     string   `xml:"id,attr"`
	Value   string   `xml:"value"`
}

func parseXMLMeasurement(byteXML []byte) (SWCMeasurement, error) {
	var values _Values

	err := xml.Unmarshal(byteXML, &values)

	if err != nil {
		return SWCMeasurement{}, err
	}

	valuesMap := make(map[string]string)
	for _, item := range values.Items {
		valuesMap[item.Ref] = item.Value
	}

	swcMeasurement := SWCMeasurement{
		HeatingOutboundTemperature:     convertToFloat64(values.Items[heatingOutboundTemperature].Value),
		HeatingInboundTemperature:      convertToFloat64(values.Items[heatingInboundTemperature].Value),
		OutsideTemperature:             convertToFloat64(values.Items[outsideTemperature].Value),
		TankTemperature:                convertToFloat64(values.Items[tankTemperature].Value),
		TargetTankTemperature:          convertToFloat64(values.Items[targetTankTemperature].Value),
		DrillInboundTemperature:        convertToFloat64(values.Items[drillInboundTemperature].Value),
		DrillOutboundTemperature:       convertToFloat64(values.Items[drillOutboundTemperature].Value),
		AmbiantIndoorTemperature:       convertToFloat64(values.Items[ambiantIndoorTemperature].Value),
		AmbiantIndoorTargetTemperature: convertToFloat64(values.Items[ambiantIndoorTargetTemperature].Value),
	}

	return swcMeasurement, nil
}

func convertToFloat64(value string) float64 {
	measureAsString := strings.Replace(value, "°C", "", 1)
	measure, err := strconv.ParseFloat(measureAsString, 64)

	if err != nil {
		log.Printf("Failed to convert %s to a float, setting it to 0.0", value)
		return 0.0
	}

	return measure
}
