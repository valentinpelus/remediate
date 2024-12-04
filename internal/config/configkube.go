package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AlertCheckInterval []string `yaml:"AlertCheckInterval"`
	QueryURL           string   `yaml:"queryUrl"`
	Servers            []Server `yaml:"servers"`
}

type Server struct {
	LabelValues []string `yaml:"labelValues"`
}

var Conf Config

func LoadConfKube(confPath string) {
	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't read kuberemediate config file")
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Unmarshal error in kube configuration parsing")
		os.Exit(1)
	}
}
