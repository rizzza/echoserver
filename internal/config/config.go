package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel     string        `envconfig:"LOG_LEVEL" default:"debug"`
	Port         string        `envconfig:"PORT" default:"8081"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
}

func Get() Config {
	var cfg Config
	if err := envconfig.Process("echoserver", &cfg); err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}

	return cfg
}

// String returns a json formatted version of the configuration with
// sensitive field redacted
func (c *Config) String() string {
	if c == nil {
		return "null" // json for nil
	}

	cfgMap := make(map[string]string)
	// We know c is of type struct, which allows us to safely use the reflect package to access
	// its members (along with the it's tags i.e. the metadata between the back-ticks)
	cfg := reflect.ValueOf(c).Elem()
	for i := 0; i < cfg.NumField(); i++ {
		name := cfg.Type().Field(i).Name
		val := fmt.Sprintf("%v", cfg.Field(i).Interface())
		tag := cfg.Type().Field(i).Tag.Get("redact")
		if redact, err := strconv.ParseBool(tag); err == nil && redact {
			val = "****"
		}
		cfgMap[name] = val
	}
	return fmt.Sprintf("Configuration: %+v", cfgMap)
}
