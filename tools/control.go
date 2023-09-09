package tools

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"sync"
)

// HmsControlData ...
type HmsControlData struct {
	Title        string
	Hash   		 string
	Description  string
	StartDate    string `json:"Start Date"`
	StartTime    string `json:"Start Time"`
	EndDate      string `json:"End Date"`
	EndTime      string `json:"End Time"`
	TimeInterval string `json:"Time Interval"`
	Notes        string
}

//Extract simulation variables from the control file...
func getControlData(hm *HmsModel, file string, wg *sync.WaitGroup, mu *sync.Mutex) {

	defer wg.Done()

	controlData := HmsControlData{}

	filePath := BuildFilePath(hm.ModelDirectory, file)

	f, err := hm.FileStore.GetObject(filePath)
	if err != nil {
		controlData.Notes += fmt.Sprintf("%s failed to process. ", file)
		mu.Lock()
		hm.Metadata.ControlMetadata[file] = controlData
		mu.Unlock()
		return
	}

	defer f.Close()

	hasher := sha256.New()

	fs := io.TeeReader(f, hasher) // fs is still a stream
	sc := bufio.NewScanner(fs)

	var line string

	for sc.Scan() {

		line = sc.Text()

		data := strings.Split(line, ": ")

		switch strings.TrimSpace(data[0]) {
		case "Control":
			controlData.Title = data[1]

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
	controlData.Hash = fmt.Sprintf("%x", hasher.Sum(nil))

	mu.Lock()
	hm.Metadata.ControlMetadata[file] = controlData
	mu.Unlock()
}
