// Copyright 2024 Christopher J. Riddick, covered by BSD 3-clause license included in this repo
//
// This package creates a GUI for monitoring and controlling CPU simulators
// It uses the fyne.io package which is cross-compatible with Linux, macOS, and Windows
// The gui window is started from the main go routine by calling dashboard.New with
// arguments including pointers to the CPU structure and host as well as functions to
// implement the GUI button responses. The call to dashboard.New returns a pointer to
// a fyne.Window and a pointer to a console buffer (bytes.Buffer) that implements the
// buffered io insterface. This buffer is used by the main program to redirect
// stdout to the GUI console.

package dashboard

import (
	"bytes"
	"image/color"
	"os"

	//"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/cjr29/go6502/cpu"
	"github.com/cjr29/go6502/host"
)

var (
	c                     *cpu.CPU
	h                     *host.Host
	CPUStatus             string
	status                string = "CPU status is displayed here."
	stackDisplay          string
	flagDisplay           string
	stackLabelWidget      *widget.Label
	stackHeader           *widget.Label
	registerHeader        *widget.Label
	registerDisplay       string
	registerDisplayWidget *widget.Label
	consoleBuffer         bytes.Buffer
	consoleDispString     string
	consoleGridLabel      *widget.Label
	consoleLabel          *widget.Label
	commandLine           *widget.Entry
	commandContainer      *fyne.Container
	loadButton            *widget.Button
	runButton             *widget.Button
	stepButton            *widget.Button
	resetButton           *widget.Button
	pauseButton           *widget.Button
	exitButton            *widget.Button
	helpButton            *widget.Button
	submitButton          *widget.Button
	assembleButton        *widget.Button
	fileButton            *widget.Button
	themeButton           *widget.Button
	currentTime           *widget.Label
	mainContainer         *fyne.Container
	buttonsContainer      *fyne.Container
	settingsContainer     *fyne.Container
	statusContainer       *fyne.Container
	registerContainer     *fyne.Container
	consoleContainer      *container.Scroll
	stackContainer        *fyne.Container
	centerContainer       *fyne.Container
	middleContainer       *fyne.Container
	fd                    *dialog.FileDialog
	selectedDirectory     *widget.Label
	selectedFile          *widget.Label
	fileURI               fyne.URI
	a                     fyne.App
	w                     fyne.Window
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpu.CPU, host *host.Host) (w fyne.Window, o *bytes.Buffer) {

	c = cpu  // All data comes from the CPU structure object
	h = host // Structure that manages the specific CPU implementation

	a = app.NewWithID("6502")
	w = a.NewWindow("6502 Simulator")

	// Color backgrounds to be used in container stacks
	registerBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	stackBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	//consoleBackground := canvas.NewRectangle(color.RGBA{R: 223, G: 159, B: 173, A: 200})

	// Control buttons
	loadButton = widget.NewButton("Load", loadSelectedFile)
	runButton = widget.NewButton("Run", run)
	stepButton = widget.NewButton("Step", step)
	resetButton = widget.NewButton("Reset", reset)
	pauseButton = widget.NewButton("Pause", pause)
	exitButton = widget.NewButton("Exit", exit)
	submitButton = widget.NewButton("Submit manual commands in field below", submit)
	helpButton = widget.NewButton("Help", help)
	selectedFile = widget.NewLabel("")
	fileButton = widget.NewButton("File", func() { showFilePicker(w) })
	assembleButton = widget.NewButton("Assemble", assembleSelectedFile)

	// Display time
	currentTime = widget.NewLabel("")

	// Command entry line
	commandLine = widget.NewEntry()
	commandLine.SetPlaceHolder("Enter command, then press Submit button. For help, type 'help'")

	commandContainer = container.NewVBox(commandLine)

	// Stack
	stackHeader = widget.NewLabel("Stack (bytes)\n(grows downward)\n")
	stackHeader.TextStyle.Monospace = true
	stackHeader.TextStyle.Bold = true
	stackDisplay = c.GetStack()
	stackLabelWidget = widget.NewLabel(stackDisplay)
	stackLabelWidget.TextStyle.Monospace = true
	stackLabelWidget.TextStyle.Bold = true
	stackContainer = container.NewStack(
		stackBackground,
		container.NewVBox(
			stackHeader,
			stackLabelWidget,
		))

	// Registers
	registerHeader = widget.NewLabel("Registers\n")
	registerHeader.TextStyle.Monospace = true
	registerHeader.TextStyle.Bold = true
	registerDisplay = c.GetRegisters()
	registerDisplayWidget = widget.NewLabel(registerDisplay)
	registerDisplayWidget.TextStyle.Monospace = true
	registerDisplayWidget.TextStyle.Bold = true
	registerContainer = container.NewStack(
		registerBackground,
		container.NewVBox(
			registerHeader,
			registerDisplayWidget,
		))

	// Console Display
	consoleDispString = consoleBuffer.String()
	consoleLabel = widget.NewLabel("Console Display\n")
	consoleLabel.TextStyle.Monospace = true
	consoleLabel.TextStyle.Bold = true
	consoleGridLabel = widget.NewLabel(consoleDispString)
	consoleGridLabel.TextStyle.Monospace = true
	consoleContainer = container.NewVScroll(
		consoleGridLabel,
	)

	buttonsContainer = container.NewHBox(
		fileButton,
		selectedFile,
		assembleButton,
		loadButton,
		//submitButton,
		resetButton,
		runButton,
		stepButton,
		pauseButton,
		exitButton,
		helpButton,
		currentTime,
	)

	settingsContainer = container.NewVBox(
		buttonsContainer,
		submitButton,
		commandContainer,
	)

	middleContainer = container.NewHBox(
		registerContainer,
		consoleContainer,
		stackContainer,
	)

	statusContainer = container.NewVBox(ConsoleScroller)
	centerContainer = container.NewHBox(consoleContainer, stackContainer)

	mainContainer = container.NewVBox(
		settingsContainer,
		middleContainer,
		statusContainer,
	)

	w.SetContent(mainContainer)
	consoleBuffer.Write([]byte("********************************************************************************\n" +
		"******************************** 6502 Simulator ********************************\n" +
		"********************************************************************************\n"))
	consoleContainer.Refresh()
	UpdateAll()

	return w, &consoleBuffer
}

func UpdateAll() {

	// Reload
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	consoleDispString = consoleBuffer.String() // Get whatever is in the memory buffer from the host
	consoleGridLabel.SetText(consoleDispString)
	registerDisplay = c.GetRegisters()
	registerDisplayWidget.Text = registerDisplay

	// Refresh
	buttonsContainer.Refresh()
	settingsContainer.Refresh()
	stackLabelWidget.Refresh()
	stackContainer.Refresh()
	consoleGridLabel.Refresh()
	consoleContainer.Refresh()
	consoleContainer.ScrollToBottom()
	registerContainer.Refresh()
	middleContainer.Refresh()
	statusContainer.Refresh()
	centerContainer.Refresh()
	mainContainer.Refresh()
}

// Return the command line entered by user when Submit is pressed
func Command() string {
	return commandLine.Text
}

// Return the command line entered by user when Submit is pressed
func ClearCmdLine() {
	commandLine.SetText("")

}

// ===================================== Button callbacks
//
// Show file picker and return selected file
func showFilePicker(w fyne.Window) {
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		saveFile := "NoFileYet"
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if f == nil {
			return
		}
		saveFile = f.URI().Path()
		fileURI = f.URI()
		selectedFile.SetText(saveFile)
	}, w)
}

// Load selected binary file
func loadSelectedFile() {
	if fileURI == nil || fileURI.Extension() != ".bin" {
		SetStatus("Binary file not selected. File must have .bin extension")
		return
	}
	h.ProcessGUICmd("load " + fileURI.Name())
	h.ProcessGUICmd("reg PC START")
	SetStatus("Program loaded")
	UpdateAll()
}

// Assemble selected source file
func assembleSelectedFile() {
	if fileURI == nil || fileURI.Extension() != ".asm" {
		SetStatus("Source file not selected. File must have .asm extension")
		return
	}
	h.ProcessGUICmd("assemble file " + fileURI.Name())
	SetStatus("Binary .bin and source map .map files generated")
	UpdateAll()
}

func run() {
	SetStatus("Running program ...")
	h.ProcessGUICmd("run")
	UpdateAll()
}

func step() {
	SetStatus("Step in ...")
	h.ProcessGUICmd("step in")
	UpdateAll()
}

func reset() {
	SetStatus("'Reset CPU.")
	h.Reset()
	UpdateAll()
}

func pause() {
	SetStatus("Pause running program.")
	h.Break()
	UpdateAll()
}

func exit() {
	SetStatus("Exit simulator")
	UpdateAll()
	os.Exit(0)
}

func help() {
	h.ProcessGUICmd("help")
	UpdateAll()
}

func submit() {
	cmd := Command()
	SetStatus(cmd)
	h.ProcessGUICmd(cmd)
	ClearCmdLine()
	UpdateAll()
}

//=========================================== End Button Callbacks ========================

// Utilities
func UpdateTime() {
	formatted := time.Now().Format("Time: 15:05:01")
	currentTime.SetText(formatted)
}

func SetStatus(s string) {
	status = s
	ConsoleWrite(status)
}

func ConsoleWrite(text string) {
	Console.Add(&canvas.Text{
		Text: text,
		//Color:     color.Black,
		TextSize:  12,
		TextStyle: fyne.TextStyle{Monospace: true},
	})

	if len(Console.Objects) > 100 {
		Console.Remove(Console.Objects[0])
	}
	delta := (Console.Size().Height - ConsoleScroller.Size().Height) - ConsoleScroller.Offset.Y

	if delta < 50 {
		ConsoleScroller.ScrollToBottom()
	}
	Console.Refresh()
}
