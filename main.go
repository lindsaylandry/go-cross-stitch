package main

import (
	"errors"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/lindsaylandry/go-cross-stitch/src/config"
	"github.com/lindsaylandry/go-cross-stitch/src/convert"
	"github.com/lindsaylandry/go-cross-stitch/src/writer"
)

var conf *config.Config

func main() {
	rootCmd := &cobra.Command{
		Use:   "cross-stitch",
		Short: "Generate cross-stitch pattern",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				panic(errors.New("No input image provided"))
			}

			return CrossStitch(args[0])
		},
	}

	var err error
	conf, err = config.NewConfig()
	if err != nil {
		panic(err)
	}

	switch conf.LogLevel {
	case 0:
		slog.SetLogLoggerLevel(slog.Level(-8))
	case 1:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case 2:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case 3:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case 4:
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func CrossStitch(filename string) error {
	c, err := convert.NewConverter(filename, conf)
	if err != nil {
		return err
	}

	d, err := c.Convert(conf.Dither)
	if err != nil {
		return err
	}

	w := writer.NewWriter(d)

	return w.WriteFiles()
}
