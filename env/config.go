package env

import (
	"github.com/noah-blockchain/noah-explorer-tools/helpers"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Init(configPath string)
}

type viperConfig struct {
}

func NewViperConfig(configPath string) Config {
	v := &viperConfig{}
	v.Init(configPath)
	return v
}

func (v *viperConfig) Init(configPath string) {
	viper.AutomaticEnv()

	fullPath := strings.Split(configPath, "/")
	configFile := fullPath[len(fullPath)-1]
	config := strings.Split(configFile, ".")

	var path string

	if len(fullPath) > 1 {
		path = strings.Join(fullPath[:len(fullPath)-1], "/")
	} else {
		path = "./"
	}

	viper.AddConfigPath(path) // path to look for the config file in
	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(config[1])
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	helpers.HandleError(err)
}

func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}


func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
