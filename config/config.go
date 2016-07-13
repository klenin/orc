package config

import (
	"strings"
	"os"
	"os/user"
	"path/filepath"
	"bufio"
	"log"
)

const PROJECT_NAME = "orc"

var configEnvPrefix = strings.ToUpper(PROJECT_NAME) + "_";

var config = map[string]string{
	"PORT": "5000",
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(configEnvPrefix + key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}

	configFileName := getEnv("CONFIG_PATH", filepath.Join(usr.HomeDir, "." + PROJECT_NAME + "rc"))
	file, err := os.Open(configFileName)
	if err != nil {
		log.Println(err.Error())
	}
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		words := strings.SplitN(scanner.Text(), "=", 2)
		if len(words) < 2 {
			log.Fatal("Error : config line %d is not correct", i)
		}
		key, value := words[0], words[1]
		config[key] = value
	}
}

func GetValue(key string) string {
	value, _ := config[key]
	return getEnv(key, value)
}
