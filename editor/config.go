package editor

import (
	"github.com/imdario/mergo"
	"github.com/kataras/iris/utils"
)

// Default values for the configuration
const (
	DefaultUsername = "iris"
	DefaultPassword = "admin!123"
)

// Config the configs for the Editor plugin
type Config struct {
	// Host if empty used the iris server's host
	Host string
	// Port if 0 4444
	Port int
	// WorkingDir if empty "./"
	WorkingDir string
	// Username if empty iris
	Username string
	// Password if empty admin!123
	Password string
}

// DefaultConfig returns the default configs for the Editor plugin
func DefaultConfig() Config {
	return Config{"", 4444, "." + utils.PathSeparator, DefaultUsername, DefaultPassword}
}

// Merge merges the default with the given config and returns the result
func (c Config) Merge(cfg []Config) (config Config) {

	if len(cfg) > 0 {
		config = cfg[0]
		mergo.Merge(&config, c)
	} else {
		_default := c
		config = _default
	}

	return
}

// MergeSingle merges the default with the given config and returns the result
func (c Config) MergeSingle(cfg Config) (config Config) {

	config = cfg
	mergo.Merge(&config, c)

	return
}
