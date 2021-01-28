package tools

import (
	"path/filepath"
	"sync"

	"github.com/USACE/filestore"
)

// Model is a general type should contain all necessary data for a model of any type.
type Model struct {
	Type           string
	Version        string
	DefinitionFile string
	Files          ModelFiles
}

// ModelFiles ...
type ModelFiles struct {
	InputFiles        InputFiles
	OutputFiles       OutputFiles
	SupplementalFiles SupplementalFiles
}

// InputFiles is a general type that should contain all data pulled from the models input files
type InputFiles struct {
	ControlFiles        ControlFiles
	ForcingFiles        ForcingFiles
	GeometryFiles       GeometryFiles
	SimulationVariables interface{} // placeholder
	LocalVariables      interface{} // placeholder
}

// ControlFiles ...
type ControlFiles struct {
	Paths []string
	Data  map[string]interface{} // placeholder
}

// ForcingFiles ...
type ForcingFiles struct {
	Paths []string
	Data  map[string]interface{} // placeholder
}

// GeometryFiles is a general type that should contain all data pulled from the models spatial files
type GeometryFiles struct {
	Paths              []string
	FeaturesProperties map[string]interface{} // placeholder
	Georeference       interface{}            // placeholder
}

// OutputFiles is a general type that should contain all data pulled from the models output files
type OutputFiles struct {
	Paths           []string
	ModelPrediction interface{} // placeholder
	RunFiles        []string
	RunLogs         []string
}

// SupplementalFiles is a general type that should contain all data pulled from the models supplemental files
type SupplementalFiles struct {
	Paths             []string
	Visulizations     interface{} // placeholder
	ObservationalData interface{} // placeholder
}

// HmsModel ...
type HmsModel struct {
	Type           string
	Version        string
	Description    string
	FileStore      filestore.FileStore
	ModelDirectory string
	DefinitionFile string
	Files          HmsModelFiles
	Metadata       HmsModelMetadata
}

// HmsModelMetadata ...
type HmsModelMetadata struct {
	ControlMetadata  map[string]HmsControlData
	ForcingMetadata  map[string]HmsForcingData
	GeometryMetadata map[string]HmsGeometryData
}

// FileExt ...
type FileExt struct {
	Definition string
	Control    string
	Forcing    string
	Geometry   string
	Grid       string
}

var hmsFileExt FileExt = FileExt{
	Definition: ".hms",
	Control:    ".control",
	Forcing:    ".met",
	Geometry:   ".basin",
	Grid:       ".grid",
}

// holder of multiple wait groups to help process files concurrency
type hmsWaitGroup struct {
	Control  sync.WaitGroup
	Forcing  sync.WaitGroup
	Geometry sync.WaitGroup
}

// IsAModel ...
func (hm *HmsModel) IsAModel() bool {
	if len(hm.Metadata.GeometryMetadata) == 0 {
		return false
	}
	return true
}

// ModelVersion ...
func (hm *HmsModel) ModelVersion() string {
	return hm.Version
}

// ModelType ...
func (hm *HmsModel) ModelType() string {
	return hm.Type
}

// IsGeospatial ...
func (hm *HmsModel) IsGeospatial() bool {
	for _, geometryData := range hm.Metadata.GeometryMetadata {

		for _, geoRefFile := range geometryData.GeoRefFiles {
			filePath := BuildFilePath(hm.ModelDirectory, geoRefFile)
			_, err := hm.FileStore.GetObject(filePath)
			if err == nil {
				return true
			}
		}
	}
	return false
}

// GeospatialData  ...
func (hm *HmsModel) GeospatialData() interface{} {
	return ""
}

// Index ...
func (hm *HmsModel) Index() Model {
	mod := Model{
		Type:           hm.Type,
		Version:        hm.Version,
		DefinitionFile: hm.DefinitionFile,
		Files: ModelFiles{
			InputFiles: InputFiles{
				ControlFiles: ControlFiles{
					Paths: hm.Files.InputFiles.ControlFiles,
					Data:  make(map[string]interface{}),
				},
				ForcingFiles: ForcingFiles{
					Paths: hm.Files.InputFiles.ForcingFiles,
					Data:  make(map[string]interface{}),
				},
				GeometryFiles: GeometryFiles{
					Paths:              hm.Files.InputFiles.GeometryFiles,
					FeaturesProperties: make(map[string]interface{}),
					Georeference:       nil,
				},
				SimulationVariables: nil,
				LocalVariables:      nil,
			},
			OutputFiles: OutputFiles{
				Paths:           hm.Files.OutputFiles.Paths(),
				ModelPrediction: nil,
				RunFiles:        make([]string, 0),
				RunLogs:         make([]string, 0),
			},
			SupplementalFiles: SupplementalFiles{
				Paths:             make([]string, 0),
				Visulizations:     nil,
				ObservationalData: nil,
			},
		},
	}
	for file, metaData := range hm.Metadata.ControlMetadata {
		mod.Files.InputFiles.ControlFiles.Data[file] = metaData
	}

	for file, metaData := range hm.Metadata.ForcingMetadata {
		mod.Files.InputFiles.ForcingFiles.Data[file] = metaData
	}

	for file, metaData := range hm.Metadata.GeometryMetadata {
		mod.Files.InputFiles.GeometryFiles.FeaturesProperties[file] = metaData
	}
	return mod
}

// NewHmsModel ...
func NewHmsModel(key string, fs filestore.FileStore) (*HmsModel, error) {
	modelMetadata := HmsModelMetadata{ControlMetadata: make(map[string]HmsControlData),
		ForcingMetadata:  make(map[string]HmsForcingData),
		GeometryMetadata: make(map[string]HmsGeometryData)}

	hm := HmsModel{Type: "HMS", FileStore: fs, ModelDirectory: filepath.Dir(key), Metadata: modelMetadata}

	err := verifyDefinitionPath(key, &hm)
	if err != nil {
		return &hm, err
	}

	err = getDefinitionData(&hm)
	if err != nil {
		return &hm, err
	}

	getGridPath(&hm)

	var hmsWG hmsWaitGroup

	for _, file := range hm.Files.Paths() {

		fileExt := filepath.Ext(file)

		switch fileExt {

		case hmsFileExt.Control:
			hmsWG.Control.Add(1)
			go getControlData(&hm, file, &hmsWG.Control)

		case hmsFileExt.Forcing:
			hmsWG.Forcing.Add(1)
			go getForcingData(&hm, file, &hmsWG.Forcing)

		case hmsFileExt.Geometry:
			hmsWG.Geometry.Add(1)
			go getGeometryData(&hm, file, &hmsWG.Geometry)
		}
	}

	hmsWG.Control.Wait()
	hmsWG.Forcing.Wait()
	hmsWG.Geometry.Wait()

	exportGeometryData(&hm)

	return &hm, nil
}
