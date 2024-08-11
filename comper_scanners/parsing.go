package comper_scanners

import (
	"encoding/json"
	"os"
)

func ReadAndUnmarshal[T any](filename string) (T, error) {
	var data T
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(fileContent, &data)
	return data, err
}
