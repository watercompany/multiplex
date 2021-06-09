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
	phaseTemplate                = "Time for phase"
	distanceBetweenPhaseAndHours = 5
)

func ParseGraphData(data string) (int, float64, error) {
	phase, err := parsePhaseNumber(data)
	if err != nil {
		return 0, 0, err
	}
	hours, err := parsePhaseHours(data)
	if err != nil {
		return 0, 0, err
	}
	return phase, hours, nil
}

func parsePhaseNumber(data string) (int, error) {
	index := strings.Index(data, phaseTemplate)
	phaseStr := getValueUntilSpace(data, index+len(phaseTemplate)+1)
	phase, err := strconv.Atoi(phaseStr)
	if err != nil {
		return 0, err
	}
	return phase, nil
}

func parsePhaseHours(data string) (float64, error) {
	index := strings.Index(data, phaseTemplate)

	hoursStr := getValueUntilSpace(data, index+len(phaseTemplate)+distanceBetweenPhaseAndHours)
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
