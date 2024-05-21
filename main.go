// Copyright 2018 Brett Vickers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"fyne.io/fyne/v2"
	"github.com/cjr29/go6502/asm"
	"github.com/cjr29/go6502/dashboard"
	"github.com/cjr29/go6502/host"
)

var (
	assemble   string
	logFile    *os.File
	err        error
	infoLogger *log.Logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func init() {
	logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	infoLogger.Println("***** host.settings.init()")
	flag.StringVar(&assemble, "a", "", "assemble file")
	flag.CommandLine.Usage = func() {
		fmt.Println("Usage: go6502 [script] ..\nOptions:")
		flag.PrintDefaults()
	}
}

func main() {
	logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	infoLogger.Println("***** Entered go6502.main()")
	fmt.Println("***** Entered go6502.main()")

	flag.Parse()

	// Initiate assembly from the command line if requested.
	if assemble != "" {
		err := asm.AssembleFile(assemble, 0, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to assemble (%v).\n", err)
		}
		os.Exit(0)
	}

	// Create the host
	infoLogger.Println("***** Create the host")
	h := host.New()
	defer h.Cleanup()

	// Run commands contained in command-line files.
	args := flag.Args()
	if len(args) > 0 {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			ioState := h.EnableProcessedMode(file, os.Stdout)
			h.RunCommands(false)
			h.RestoreIoState(ioState)
			file.Close()
		}
	}

	// Break on Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleInterrupt(h, c)

	// Open dashboard
	// a := app.New()
	// w := a.NewWindow("Hello World")
	// w.SetContent(widget.NewLabel("Hello World!"))
	// w.ShowAndRun()
	infoLogger.Println("***** Open dashboard.")
	os.Setenv("FYNE_THEME", "light")

	// Set up Fyne window before trying to write to Status line!!!
	var w fyne.Window = dashboard.New(h.GetCPU())
	// Activate dashboard process
	w.ShowAndRun()

	// Interactively run commands entered by the user.
	infoLogger.Println("***** Interactively run commands entered by the user.")
	h.EnableRawMode()
	h.RunCommands(true)

}

func handleInterrupt(h *host.Host, c chan os.Signal) {
	for {
		<-c
		h.Break()
	}
}

func exitOnError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	os.Exit(1)
}
