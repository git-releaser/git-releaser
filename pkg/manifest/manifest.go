package manifest

import (
	"encoding/json"
	"github.com/Masterminds/semver"
	"github.com/thschue/git-releaser/pkg/naming"
	"io"
	"os"
)

func GetCurrentVersion() (*semver.Version, error) {
	jsonFile, err := os.Open(naming.DefaultManifestFileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}

	version, err := semver.NewVersion(result["version"].(string))
	if err != nil {
		return nil, err
	}
	return version, nil
}
