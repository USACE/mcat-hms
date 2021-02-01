package tools

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

// HmsForcingData ...
type HmsForcingData struct {
	Title            string
	Description      string
	Units            string `json:"Unit System"`
	MissingToDefault string `json:"Set Missing Data to Default"`
	Precip           string `json:"Precipitation"`
	SWave            string `json:"Short-Wave Radiation"`
	LWave            string `json:"Long-Wave Radiation"`
	Snowmelt         string
	ET               string   `json:"Evapotranspiration"`
	BasinModel       []string `json:"Use Basin Model"`
	Subbasin         []string
	Notes            string
}

//Extract meteorological variables from the forcing files...
func getForcingData(hm *HmsModel, file string, wg *sync.WaitGroup) {

	defer wg.Done()

	forcingData := HmsForcingData{}

	filePath := BuildFilePath(hm.ModelDirectory, file)

	f, err := hm.FileStore.GetObject(filePath)
	if err != nil {
		forcingData.Notes += fmt.Sprintf("%s failed to process. ", file)
		hm.Metadata.ForcingMetadata[file] = forcingData
		return
	}

	defer f.Close()

	sc := bufio.NewScanner(f)

	var line string

	for sc.Scan() {

		line = sc.Text()

		data := strings.Split(line, ": ")

		switch strings.TrimSpace(data[0]) {

		case "Meteorology":
			forcingData.Title = data[1]

		case "Description":
			forcingData.Description = data[1]

		case "Unit System":
			forcingData.Units = data[1]

		case "Set Missing Data to Default":
			forcingData.MissingToDefault = data[1]

		case "Precipitation Method":
			forcingData.Precip = data[1]

		case "Short-Wave Radiation Method":
			forcingData.SWave = data[1]

		case "Long-Wave Radiation Method":
			forcingData.LWave = data[1]

		case "Snowmelt Method":
			forcingData.Snowmelt = data[1]

		case "Evapotranspiration Method":
			forcingData.ET = data[1]

		case "Use Basin Model":
			forcingData.BasinModel = append(forcingData.BasinModel, data[1])

		case "Subbasin":
			forcingData.Subbasin = append(forcingData.Subbasin, data[1])

		}
	}
	hm.Metadata.ForcingMetadata[file] = forcingData
}
