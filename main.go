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
	statusChan chan string // channel to send status to dashboard
	h          *host.Host
)

func init() {
	logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	infoLogger.Println("***** host.settings.init()")

	// Initialize the startup parameters to be parsed in command line
	flag.StringVar(&assemble, "a", "", "assemble file")
	flag.BoolVar(&gui, "g", false, "Activate GUI")
	flag.CommandLine.Usage = func() {
		fmt.Println("Usage: go6502 [script] ..\nOptions:")
		flag.PrintDefaults()
	}

}

func main() {
	logFile, err = os.OpenFile("6502Emu.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	infoLogger.Println("***** Entered go6502.main()")
	// Create the host
	infoLogger.Println("***** Create the host")
	h = host.New()
	defer h.Cleanup()

	flag.Parse()
	infoLogger.Printf("GUI: %v, Assemble: %s", gui, assemble)
	fmt.Println("GUI: ", gui)
	fmt.Println("Assemble: ", assemble)

	// Create dashboard GUI
	infoLogger.Println("***** Open dashboard.")
	os.Setenv("FYNE_THEME", "light")
	// Set up Fyne window before trying to write to Status line!!!
	w, statusChan = dashboard.New(h.GetCPU(), reset, load, step, run, pause, exit)
	// Kick off test ticker to report to dashboard very second
	go ticker() // Remove after testing complete, used to verify program running

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
				exitOnError(err)
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
		h.EnableProcessedMode(logFile, os.Stdout)
		w.ShowAndRun()
	}
	// !!!!!!!! Will never get to next line if ShowAnRun() is executed. Program won't return to the main
	// thread until the fyne window is closed.

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

func ticker() {
	infoLogger.Println("Start Ticker Example 1")

	// define an interval and the ticker for this interval
	interval := time.Duration(2) * time.Second
	// create a new Ticker
	tk := time.NewTicker(interval)
	// start the ticker by constructing a loop
	i := 0
	for range tk.C {
		i++
		//countFuncCall(i)
		statusChan <- fmt.Sprintf("Tick at: %s  --- count = %d", time.Now().UTC(), i)
	}
}

func load() {
	//dashboard.UpdateAll()
	//statusChan <- "'Load' pressed"
	//dashboard.SetStatus("'Load' pressed")
	h.ProcessGUICmd("load monitor.bin $F800")
	h.ProcessGUICmd("load sample.bin")
	h.ProcessGUICmd("set compact true")
	h.ProcessGUICmd("reg PC START")
	h.ProcessGUICmd("d .")
	dashboard.UpdateAll()
}

func run() {
	//dashboard.UpdateAll()
	//statusChan <- "'Run' pressed"
	dashboard.SetStatus("'Run' pressed")
	dashboard.UpdateAll()
}

func step() {
	h.ProcessGUICmd("step in")
	//statusChan <- "'Step' pressed"
	//dashboard.SetStatus("'Step' pressed")
	dashboard.UpdateAll()
}

func reset() {
	//dashboard.UpdateAll()
	//statusChan <- "'Reset' pressed"
	dashboard.SetStatus("'Reset' pressed")
	dashboard.UpdateAll()
}

func pause() {
	//dashboard.UpdateAll()
	//statusChan <- "'Pause' pressed"
	dashboard.SetStatus("'Pause' pressed")
	dashboard.UpdateAll()
}

func exit() {
	dashboard.SetStatus("'Exit' pressed")
	dashboard.UpdateAll()
	os.Exit(0)
}
