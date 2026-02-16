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
	ColorGrid bool     `yaml:"color_grid"`
	CsvFile   string   `yaml:"csv_file"`
	Width     int      `yaml:"width"`
	LogLevel  int      `yaml:"log_level"`
	Metric    bool     `yaml:"metric"`
	Excludes  []string `yaml:"excludes"`
	DMC       Type     `yaml:"dmc"`
	Anchor    Type     `yaml:"anchor"`
	Lego      Type     `yaml:"lego"`
}

type Quantize struct {
	Enabled bool `yaml:"enabled"`
	N       int  `yaml:"n"`
}

type Type struct {
	PixelSizeMM float32 `yaml:"pixel_size_mm"`
	Fabric      Fabric  `yaml:"fabric"`
}

type Fabric struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Color   string `yaml:"color"`
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
