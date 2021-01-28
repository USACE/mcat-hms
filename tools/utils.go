package tools

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/USACE/filestore"
)

// readFirstLine ...
func readFirstLine(fs filestore.FileStore, fn string) (string, error) {
	file, err := fs.GetObject(fn)
	if err != nil {
		fmt.Println("Couldn't open the file", fn)
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	return rmNewLineChar(line), err
}

// rmNewLineChar ...
func rmNewLineChar(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
}

//BuildFilePath ... build the file path given its name and the file path of the definition file...
func BuildFilePath(modelDirectory, fileName string) string {
	return filepath.Join(modelDirectory, strings.Replace(fileName, "\\", "/", -1))

}
