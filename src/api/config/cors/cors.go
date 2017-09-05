package cors

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type corsYAML struct {
	Cors Cors
}

// Cors configuration used during middleware setup
type Cors struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

func defaultCors() (Cors, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		return Cors{
			AllowedOrigins: []string{"*"},
		}, nil
	}
	// Defaults
	return Cors{
		AllowedOrigins: []string{""},
		AllowedHeaders: []string{"content-type", "x-xsrf-token"},
	}, nil
}

// Config returns the CORS configuration for the environment
func Config(configFile string) (Cors, error) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		log.Info("Loading default CORS config")
		return defaultCors()
	}

	log.Info("Loading CORS from config file")
	cors, err := loadCorsFromFile(configFile)
	if err != nil {
		return Cors{}, err
	}

	if len(cors.AllowedOrigins) == 0 && len(cors.AllowedHeaders) == 0 {
		return defaultCors()
	}

	return cors, nil
}

func loadCorsFromFile(filePath string) (Cors, error) {
	var yamlStruct corsYAML
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Cors{}, err
	}
	if err := yaml.Unmarshal(bytes, &yamlStruct); err != nil {
		return Cors{}, err
	}
	return yamlStruct.Cors, nil
}
