package cors

import "os"

// Cors configuration used during middleware setup
type Cors struct {
	AllowedOrigins []string
	AllowedHeaders []string
}

var currentEnv = func() string {
	return os.Getenv("ENVIRONMENT")
}

// Config returns the CORS configuration for the environment
// TODO, read the configuration from the overrides config file argument
func Config(configFile string) (Cors, error) {
	env := currentEnv()
	if env == "development" {
		return Cors{
			AllowedOrigins: []string{"*"},
		}, nil
	}
	// Defaults
	return Cors{
		AllowedOrigins: []string{"my-api-server"},
		AllowedHeaders: []string{"access-control-allow-headers", "x-xsrf-token"},
	}, nil
}
