package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Rgb       bool   `yaml:"rgb"`
	All       bool   `yaml:"all"`
	Palette   string `yaml:"palette"`
	Dither    bool   `yaml:"dither"`
	Greyscale bool   `yaml:"greyscale"`
	Quantize  `yaml:"quantize"`
	ColorGrid bool   `yaml:"color_grid"`
	CsvFile   string `yaml:"csv_file"`
	Width     int    `yaml:"width"`
	LogLevel  int    `yaml:"log_level"`
}

type Quantize struct {
	Enabled bool `yaml:"enabled"`
	N       int  `yaml:"n"`
}

func NewConfig() (*Config, error) {
	c := Config{}

	data, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		return &c, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	return &c, err
}
