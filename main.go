// Copyright 2018 Brett Vickers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"fyne.io/fyne/v2"
	"github.com/cjr29/go6502/asm"
	"github.com/cjr29/go6502/dashboard"
	"github.com/cjr29/go6502/host"
)

var (
	assemble   string
	gui        bool
	logFile    *os.File
	err        error
	infoLogger *log.Logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	w          fyne.Window
	h          *host.Host
	outbuffer  *bytes.Buffer
)

func init() {
	/* 	logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	   	if err != nil {
	   		log.Fatal("Failed to open log file:", err)
	   	}
	   	infoLogger.Println("***** host.settings.init()") */

	// Initialize the startup parameters to be parsed in command line
	flag.StringVar(&assemble, "a", "", "assemble file")
	flag.BoolVar(&gui, "g", false, "Activate GUI")
	flag.CommandLine.Usage = func() {
		fmt.Println("Usage: go6502 [script] ..\nOptions:")
		flag.PrintDefaults()
	}
}

func main() {
	/* logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	} */

	//infoLogger.Println("***** Entered go6502.main()")

	// Create the host
	//infoLogger.Println("***** Create the host")
	h = host.New()
	defer h.Cleanup()

	flag.Parse()
	//infoLogger.Printf("GUI: %v, Assemble: %s", gui, assemble)

	// Create dashboard GUI
	//infoLogger.Println("***** Open dashboard.")

	// Set up Fyne window before trying to write to Status line!!!
	w, outbuffer = dashboard.New(h.GetCPU(), h)

	// Initiate assembly from the command line if requested.
	if assemble != "" {
		err := asm.AssembleFile(assemble, 0, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to assemble (%v).\n", err)
		}
		os.Exit(0)
	}

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

	// Activate dashboard process if startup flag set to true (-g)
	if gui {
		// Run every second in background to update dashboard current time display
		go func() {
			for range time.Tick(time.Second) {
				dashboard.UpdateTime()
			}
		}()

		h.EnableProcessedMode(os.Stdin, outbuffer)
		w.ShowAndRun()
	}
	// !!!!!!!! Will never get to next line if ShowAnRun() is executed. Program won't return to the main
	// thread until the gui window is closed.

	// Interactively run commands entered by the user.
	//infoLogger.Println("***** Interactively run commands entered by the user.")
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
