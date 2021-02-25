package swcsource

import (
	"encoding/xml"
	"log"
	"strconv"
	"strings"
)

const heatingOutboundTemperature = "0x46a0ac"
const heatingInboundTemperature = "0x49902c"
const outsideTemperature = "0x498414"
const tankTemperature = "0x4976e4"
const targetTankTemperature = "0x497b4c"
const drillInboundTemperature = "0x497ef4"
const drillOutboundTemperature = "0x497494"
const ambiantIndoorTemperature = "0x4678d4"
const ambiantIndoorTargetTemperature = "0x46790c"

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
		HeatingOutboundTemperature:     convertToFloat64(valuesMap[heatingOutboundTemperature]),
		HeatingInboundTemperature:      convertToFloat64(valuesMap[heatingInboundTemperature]),
		OutsideTemperature:             convertToFloat64(valuesMap[outsideTemperature]),
		TankTemperature:                convertToFloat64(valuesMap[tankTemperature]),
		TargetTankTemperature:          convertToFloat64(valuesMap[targetTankTemperature]),
		DrillInboundTemperature:        convertToFloat64(valuesMap[drillInboundTemperature]),
		DrillOutboundTemperature:       convertToFloat64(valuesMap[drillOutboundTemperature]),
		AmbiantIndoorTemperature:       convertToFloat64(valuesMap[ambiantIndoorTemperature]),
		AmbiantIndoorTargetTemperature: convertToFloat64(valuesMap[ambiantIndoorTargetTemperature]),
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
