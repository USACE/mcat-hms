package tools

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
)

// HmsGeometryData ...
type HmsGeometryData struct {
	Title          string
	Hash           string
	Description    string
	Units          string `json:"Unit System"`
	MissingtoZero  string `json:"Missing Flow To Zero"`
	FlowRatio      string `json:"Enable Flow Ratio"`
	LocalFlow      string `json:"Local Flow At Junctions"`
	SedRouting     string `json:"Sediment Routing"`
	QualityRouting string `json:"Quality Routing"`
	Features       map[string][]string
	GeoRefFiles    []string `json:"Geospatial Reference Files"`
	CRS            string   `json:"Coordinate System"`
	LossRate       string   `json:"LossRate"`
	Transform      string   `json:"Transform"`
	Notes          string
}

// GeometryFeatureTypes ...
var GeometryFeatureTypes []string = []string{"Subbasin", "Reach", "Junction", "Source", "Sink", "Reservoir", "Diversion"}

// getGridPath ...
func getGridPath(hm *HmsModel) {
	matchingFiles := make([]string, 0)

	prefix := hm.ModelDirectory + "/"

	files, err := hm.FileStore.GetDir(prefix, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range *files {
		if file.Type == hmsFileExt.Grid {
			firstLine, err := readFirstLine(hm.FileStore, filepath.Join(file.Path, file.Name))
			if err != nil {
				fmt.Println(err)
				continue
			}
			if strings.Contains(firstLine, "Grid Manager:") {
				matchingFiles = append(matchingFiles, file.Name)
			}
		}
	}

	hm.Files.SupplementalFiles.GridFiles = matchingFiles
	return
}

//Extract features and their properties from the geometry files...
func getGeometryData(hm *HmsModel, file string, wg *sync.WaitGroup, mu *sync.Mutex) {

	defer wg.Done()

	geometryData := HmsGeometryData{Features: make(map[string][]string)}
	var basinProperties bool = false
	var endBasinProperties bool = true

	filePath := BuildFilePath(hm.ModelDirectory, file)

	f, err := hm.FileStore.GetObject(filePath)
	if err != nil {
		geometryData.Notes += fmt.Sprintf("%s failed to process. ", file)
		mu.Lock()
		hm.Metadata.GeometryMetadata[file] = geometryData
		mu.Unlock()
		return
	}

	defer f.Close()

	hasher := sha256.New()

	fs := io.TeeReader(f, hasher) // fs is still a stream
	sc := bufio.NewScanner(fs)

	var line string

out:
	for sc.Scan() {

		line = sc.Text()

		data := strings.Split(line, ":")

		key := strings.TrimSpace(data[0])

		switch key {

		case "Basin":
			geometryData.Title = strings.TrimSpace(data[1])

		case "Description":
			geometryData.Description = strings.TrimSpace(data[1])

		case "Unit System":
			geometryData.Units = strings.TrimSpace(data[1])

		case "Missing Flow To Zero":
			geometryData.MissingtoZero = strings.TrimSpace(data[1])

		case "Enable Flow Ratio":
			geometryData.FlowRatio = strings.TrimSpace(data[1])

		case "Compute Local Flow At Junctions":
			geometryData.LocalFlow = strings.TrimSpace(data[1])

		case "Enable Sediment Routing":
			geometryData.SedRouting = strings.TrimSpace(data[1])

		case "Enable Quality Routing":
			geometryData.QualityRouting = strings.TrimSpace(data[1])

		case "Transform":
			geometryData.Transform = strings.TrimSpace(data[1])

		case "LossRate":
			geometryData.LossRate = strings.TrimSpace(data[1])

		case "File":
			filename := strings.TrimSpace(data[1])
			fileparts := strings.Split(filename, ".")
			if fileparts[1] == "sqlite" {
				for _, existingFile := range geometryData.GeoRefFiles {
					if filename == existingFile {
						continue out
					}
				}
				geometryData.GeoRefFiles = append(geometryData.GeoRefFiles, filename)
			}

		case "Basin Layer Properties":
			basinProperties = true
			endBasinProperties = false

		case "End":
			if basinProperties {
				endBasinProperties = true
			}

		case "Filename":
			if basinProperties && !endBasinProperties {
				geometryData.GeoRefFiles = append(geometryData.GeoRefFiles, strings.TrimSpace(data[1]))
			}

		case "Coordinate System":
			geometryData.CRS = strings.TrimSpace(data[1])

		}

		for _, featureType := range GeometryFeatureTypes {
			if key == featureType {
				geometryData.Features[featureType] = append(geometryData.Features[featureType], strings.TrimSpace(data[1]))
			}
		}

	}
	geometryData.Hash = fmt.Sprintf("%x", hasher.Sum(nil))

	mu.Lock()
	hm.Metadata.GeometryMetadata[file] = geometryData
	mu.Unlock()
}

//Check that the geometry reference files exists, read them into memory, and serialize ...
func exportGeometryData(hm *HmsModel) {
	geometryMetadata := hm.Metadata.GeometryMetadata

	for file, geometryData := range geometryMetadata {

		for _, geoRefFile := range geometryData.GeoRefFiles {

			filePath := BuildFilePath(hm.ModelDirectory, geoRefFile)

			_, err := hm.FileStore.GetObject(filePath)
			if err != nil {
				geometryData.Notes += fmt.Sprintf("%s does not exist. ", geoRefFile)
			}
		}
		hm.Metadata.GeometryMetadata[file] = geometryData
	}
}
