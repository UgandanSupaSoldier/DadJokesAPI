package shared

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

var config *toml.Tree
var configPath = "config.toml"

func SeConfigtPath(path string) {
	configPath = path
}

func loadConf() (*toml.Tree, error) {
	if config == nil {
		configFile, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer configFile.Close()

		config, err = toml.LoadReader(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}
	return config, nil
}

func getConf(name string) (interface{}, error) {
	conf, err := loadConf()
	if err != nil {
		return "", err
	}
	value := conf.Get(name)
	if value == nil {
		return "", fmt.Errorf("config key %s not found", name)
	}
	return value, nil
}

func getConfDef(name string, def interface{}) interface{} {
	conf, err := loadConf()
	if err != nil {
		return def
	}

	value := conf.Get(name)
	if value == nil {
		return def
	}
	return value
}

func GetStr(name string) (string, error) {
	value, err := getConf(name)
	return value.(string), err
}

func GetStrDef(name string, def string) string {
	return getConfDef(name, def).(string)
}

func GetInt(name string) (int, error) {
	value, err := getConf(name)
	return int(value.(int64)), err
}

func GetIntDef(name string, def int) int {
	return int(getConfDef(name, int64(def)).(int64))
}

func GetBool(name string) (bool, error) {
	value, err := getConf(name)
	return value.(bool), err
}

func GetBoolDef(name string, def bool) bool {
	return getConfDef(name, def).(bool)
}
