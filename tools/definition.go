package tools

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// HmsModelFiles ...
type HmsModelFiles struct {
	InputFiles        HmsInputFiles
	OutputFiles       HmsOutputFiles
	SupplementalFiles HmsSupplementalFiles
}

// HmsInputFiles ...
type HmsInputFiles struct {
	ControlFiles  []string
	ForcingFiles  []string
	GeometryFiles []string
}

// HmsOutputFiles ...
type HmsOutputFiles struct {
	PredictionFiles []string
	RunFiles        []string
	RunLogs         []string
}

// HmsSupplementalFiles ...
type HmsSupplementalFiles struct {
	GridFiles          []string
	VisualizationFiles []string
	ObservationFiles   []string
}

// Paths will pull all paths from the HmsInputFiles, HmsOutputFiles, and HmsSupplementalFiles members of the model type
func (mf HmsModelFiles) Paths() []string {
	paths := make([]string, 0)
	paths = append(paths, mf.InputFiles.Paths()...)
	paths = append(paths, mf.OutputFiles.Paths()...)
	paths = append(paths, mf.SupplementalFiles.Paths()...)
	return paths
}

// Paths ...
func (i HmsInputFiles) Paths() []string {
	paths := make([]string, 0)
	paths = append(paths, i.ControlFiles...)
	paths = append(paths, i.ForcingFiles...)
	paths = append(paths, i.GeometryFiles...)
	return paths
}

// Paths ...
func (o HmsOutputFiles) Paths() []string {
	paths := make([]string, 0)
	paths = append(paths, o.PredictionFiles...)
	paths = append(paths, o.RunFiles...)
	paths = append(paths, o.RunLogs...)
	return paths
}

// Paths ...
func (s HmsSupplementalFiles) Paths() []string {
	paths := make([]string, 0)
	paths = append(paths, s.GridFiles...)
	paths = append(paths, s.VisualizationFiles...)
	paths = append(paths, s.ObservationFiles...)
	return paths
}

// verifyDefinitionPath ...
func verifyDefinitionPath(key string, hm *HmsModel) error {

	if filepath.Ext(key) != ".hms" {
		return fmt.Errorf("%s is not a .hms file", key)
	}

	firstLine, err := readFirstLine(hm.FileStore, key)
	if err != nil {
		return err
	}
	if !strings.Contains(firstLine, "Project:") {
		return fmt.Errorf("%s is not a HMS Definition file", key)
	}

	hm.DefinitionFile = filepath.Base(key)

	return nil
}

// nextLineData ...
func nextLineData(sc *bufio.Scanner, delimiter string) string {
	sc.Scan()
	nextdata := strings.Split(sc.Text(), delimiter)
	return strings.TrimSpace(nextdata[1])
}

// Extract all model file paths from the definition file...
func getDefinitionData(hm *HmsModel) error {
	defFilePath := BuildFilePath(hm.ModelDirectory, hm.DefinitionFile)
	inputFiles := HmsInputFiles{}
	outputFiles := HmsOutputFiles{}

	f, err := hm.FileStore.GetObject(defFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	projectBlock := true
	for sc.Scan() {
		line := sc.Text()

		match, err := regexp.MatchString(":", line)
		if err != nil {
			return err
		}
		if match {
			data := strings.Split(line, ": ")

			switch strings.TrimSpace(data[0]) {

			case "Description:":
				if projectBlock {
					hm.Description = data[1]
				}

			case "Version":
				hm.Version = data[1]

			case "DSS File Name":
				outputFiles.PredictionFiles = append(outputFiles.PredictionFiles, data[1])

			case "End":
				projectBlock = false

			case "Control":
				nextdata := nextLineData(sc, ": ")
				inputFiles.ControlFiles = append(inputFiles.ControlFiles, nextdata)

			case "Precipitation":
				nextdata := nextLineData(sc, ": ")
				inputFiles.ForcingFiles = append(inputFiles.ForcingFiles, nextdata)

			case "Basin":
				nextdata := nextLineData(sc, ": ")
				inputFiles.GeometryFiles = append(inputFiles.GeometryFiles, nextdata)

			}
		}
	}

	hm.Files = HmsModelFiles{inputFiles, outputFiles, HmsSupplementalFiles{}}

	return nil
}
