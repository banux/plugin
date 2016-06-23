package iriscontrol

import "github.com/imdario/mergo"

const (
	// TimeFormat default time format for any kind of datetime parsing
	TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
)

var (
	// DefaultUsername used for default for Ethe ditor's configuration
	DefaultUsername = "iris"
	// DefaultPassword used for default for the Editor's configuration
	DefaultPassword = "admin!123"
)

// Config the options which iris control needs
// contains the port (int) and authenticated users with their passwords (map[string]string)
type Config struct {
	// Port the port
	Port int
	// Users the authenticated users, [username]password
	Users map[string]string
}

// DefaultConfig returns the default configs for IrisControl plugin
func DefaultConfig() Config {
	users := make(map[string]string, 0)
	users[DefaultUsername] = DefaultPassword
	return Config{4000, users}
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
