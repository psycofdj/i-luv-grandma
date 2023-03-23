// Copyright 2023 Xavier MARCELET. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"gihub.com/psycofdj/i-luv-grandma/pbm"
)

var (
	Version   = "unknown"
	Revision  = "unknown"
	Branch    = "unknown"
	BuildUser = "unknown"
	BuildDate = "unknown"
)

type App struct {
	help           bool
	version        bool
	profilePath    string
	inputFilePath  string
	outputFilePath string
	rotationAngle  float64
}

func NewApp() *App {
	return &App{}
}

func (a *App) printUsage() {
	stream := flag.CommandLine.Output()
	fmt.Fprintf(stream, "usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(stream)
	fmt.Fprintf(stream, "Rotate pbm image by given angle. Result is written to output file.\n")
	fmt.Fprintln(stream)
	flag.PrintDefaults()
}

func (a *App) printVersion() {
	fmt.Printf("%s, %s (branch: %s, revision: %s)\n", os.Args[0], Version, Branch, Revision)
	fmt.Printf("build user: %s\n", BuildUser)
	fmt.Printf("build date: %s\n", BuildDate)
	fmt.Printf("go version: %s\n", runtime.Version())
	fmt.Printf("platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func (a *App) parseArgs() {
	flag.BoolVar(&a.help, "help", false, "print usage")
	flag.BoolVar(&a.version, "version", false, "outputs version and revision informations")
	flag.StringVar(&a.profilePath, "profile", "", "generate pprof profile output")
	flag.StringVar(&a.inputFilePath, "input", "input.pbm", "process given input file path, '-' for stdin")
	flag.StringVar(&a.outputFilePath, "output", "output.pbm", "write to given output file path, '-' for stdout")
	flag.Float64Var(&a.rotationAngle, "angle", 90, "rotation of given decimal angle (positive or negative)")
	flag.Parse()
	if a.help {
		a.printUsage()
		os.Exit(0)
	}
	if a.version {
		a.printVersion()
		os.Exit(0)
	}
}

func (a *App) run() error {
	a.parseArgs()
	if len(a.profilePath) != 0 {
		f, err := os.Create(a.profilePath)
		if err != nil {
			return fmt.Errorf("unable to create profile file '%s': %s", a.profilePath, err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not enable CPU profiling: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	image, err := pbm.NewImageFromFile(a.inputFilePath)
	if err != nil {
		return err
	}

	image.Rotate(a.rotationAngle)
	if err := image.EncodeASCIIToFile(a.outputFilePath); err != nil {
		return fmt.Errorf("could not write output file '%s': %s", a.outputFilePath, err)
	}

	return nil
}

func main() {
	app := NewApp()
	if err := app.run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
