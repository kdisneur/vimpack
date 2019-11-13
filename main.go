package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"
	"vimpack/internal"
)

type Config struct {
	Destination string
	AskVersion  bool
	Vimpackfile string
}

func main() {
	var config Config
	parseFlags(&config)

	if config.AskVersion {
		fmt.Printf("%#+v", internal.GetVersionInfo())
		os.Exit(0)
	}

	plugins, err := internal.NewParser().ParseFile(config.Vimpackfile)
	if err != nil {
		fatalf("%s\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancelOnKill(ctx, cancel)

	ui := internal.NewUI(ctx, os.Stdout)

	fetcher := internal.NewFetcher(ui)

	fetcher.All(ctx, plugins, config.Destination)
}

func cancelOnKill(ctx context.Context, cancelFunc func()) {
	killSignal := make(chan os.Signal, 3)
	signal.Notify(killSignal, os.Interrupt)

	go func(sig <-chan os.Signal, cancelFunc func()) {
		<-sig
		cancelFunc()
	}(killSignal, cancelFunc)
}

func fatalf(format string, variables ...interface{}) {
	fmt.Printf(format, variables...)
	os.Exit(1)
}

func parseFlags(c *Config) {
	currentUser, err := user.Current()
	if err != nil {
		fatalf("can't get current user: %s\n", err)
	}

	defaultVimpackFolder := path.Join(currentUser.HomeDir, ".vim", "pack")
	defaultVimpackfilePath := path.Join(currentUser.HomeDir, ".vim", "Vimpackfile")

	flag.BoolVar(&c.AskVersion, "version", false, "display current version")
	flag.StringVar(&c.Destination, "dest", defaultVimpackFolder, "path where to download plugins")
	flag.StringVar(&c.Vimpackfile, "file", defaultVimpackfilePath, "path to the Vimpackfile")

	flag.Parse()
}
