package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type PlotGraph struct {
	Data []PlotData `json:"data"`
}

type PlotData struct {
	Phase int     `json:"phase"`
	Hours float64 `json:"hours"`
}

const (
	phaseStringTemplate            = "Time for phase"
	distanceBetweenPhaseAndHours   = 5
	phaseStringTemplateV2          = "Phase"
	distanceBetweenPhaseAndHoursV2 = 8
)

func ParseGraphData(data, taskName string) (int, float64, error) {
	phaseTemplate := phaseStringTemplate
	hoursTemplate := distanceBetweenPhaseAndHours

	if taskName == "posv2" {
		phaseTemplate = phaseStringTemplateV2
		hoursTemplate = distanceBetweenPhaseAndHoursV2
	}

	phase, err := parsePhaseNumber(data, phaseTemplate)
	if err != nil {
		return 0, 0, err
	}
	hours, err := parsePhaseHours(data, phaseTemplate, hoursTemplate)
	if err != nil {
		return 0, 0, err
	}
	return phase, hours, nil
}

func parsePhaseNumber(data, template string) (int, error) {
	index := strings.Index(data, template)
	phaseStr := getValueUntilSpace(data, index+len(template)+1)
	phase, err := strconv.Atoi(phaseStr)
	if err != nil {
		return 0, err
	}
	return phase, nil
}

func parsePhaseHours(data, phaseTemplate string, phaseHoursDistance int) (float64, error) {
	index := strings.Index(data, phaseTemplate)

	hoursStr := getValueUntilSpace(data, index+len(phaseTemplate)+phaseHoursDistance)
	hours, err := strconv.ParseFloat(hoursStr, 32)
	if err != nil {
		return 0, err
	}

	// convert seconds to hours
	hours = hours / 3600

	// round to 2 decimal places
	hours = math.Round(hours*100) / 100

	return hours, nil
}

func getValueUntilSpace(info string, index int) string {
	end := strings.Index(info[index:], " ") + index
	return info[index:end]
}

func SavePlotGraphDataToJSON(data PlotGraph, filename string) error {
	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(fmt.Sprintf("%s.json", filename), file, 0777)
	return nil
}
