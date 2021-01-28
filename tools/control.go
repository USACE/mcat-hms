package tools

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

// HmsControlData ...
type HmsControlData struct {
	Description  string
	StartDate    string `json:"Start Date"`
	StartTime    string `json:"Start Time"`
	EndDate      string `json:"End Date"`
	EndTime      string `json:"End Time"`
	TimeInterval string `json:"Time Interval"`
	Notes        string
}

//Extract simulation variables from the control file...
func getControlData(hm *HmsModel, file string, wg *sync.WaitGroup) {

	defer wg.Done()

	controlData := HmsControlData{}

	filePath := BuildFilePath(hm.ModelDirectory, file)

	f, err := hm.FileStore.GetObject(filePath)
	if err != nil {
		controlData.Notes += fmt.Sprintf("%s failed to process. ", file)
		hm.Metadata.ControlMetadata[file] = controlData
		return
	}

	defer f.Close()

	sc := bufio.NewScanner(f)

	var line string

	for sc.Scan() {

		line = sc.Text()

		data := strings.Split(line, ": ")

		switch strings.TrimSpace(data[0]) {

		case "Description":
			controlData.Description = data[1]

		case "Start Date":
			controlData.StartDate = data[1]

		case "Start Time":
			controlData.StartTime = data[1]

		case "End Date":
			controlData.EndDate = data[1]

		case "End Time":
			controlData.EndTime = data[1]

		case "Time Interval":
			controlData.TimeInterval = data[1]
		}
	}
	hm.Metadata.ControlMetadata[file] = controlData
}
