package dashboard

import (
	"bytes"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cjr29/go6502/cpu"
	"github.com/cjr29/go6502/host"
)

var (
	c         *cpu.CPU
	h         *host.Host
	CPUStatus string

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
	//inputCPUClock         *widget.Entry
	loadButton        *widget.Button
	runButton         *widget.Button
	stepButton        *widget.Button
	resetButton       *widget.Button
	pauseButton       *widget.Button
	exitButton        *widget.Button
	submitButton      *widget.Button
	currentTime       *widget.Label
	mainContainer     *fyne.Container
	buttonsContainer  *fyne.Container
	settingsContainer *fyne.Container
	statusContainer   *fyne.Container
	registerContainer *fyne.Container
	consoleContainer  *container.Scroll
	stackContainer    *fyne.Container
	speedContainer    *fyne.Container
	centerContainer   *fyne.Container
	middleContainer   *fyne.Container
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpu.CPU, host *host.Host, submit func(), reset func(), load func(), step func(), run func(), pause func(), exit func()) (w fyne.Window, o *bytes.Buffer) {

	c = cpu  // All data comes from the CPU structure object
	h = host // Structure that manages the specific CPU implementation

	a := app.NewWithID("6502")
	w = a.NewWindow("6502 Simulator")

	// Color backgrounds to be used in container stacks
	registerBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	stackBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	//consoleBackground := canvas.NewRectangle(color.RGBA{R: 223, G: 159, B: 173, A: 200})

	// Control buttons
	loadButton = widget.NewButton("Load", load)
	runButton = widget.NewButton("Run", run)
	stepButton = widget.NewButton("Step", step)
	resetButton = widget.NewButton("Reset", reset)
	pauseButton = widget.NewButton("Pause", pause)
	exitButton = widget.NewButton("Exit", exit)
	submitButton = widget.NewButton("Submit", submit)

	// Display time
	currentTime = widget.NewLabel("")

	// Clock settings line
	/* 	inputCPUClock = widget.NewEntry()
	   	inputCPUClock.SetText("1")
	   	speedContainer = container.NewHBox(
	   		canvas.NewText("Clock Speed = ", color.Black),
	   		inputCPUClock,
	   		canvas.NewText("ms  ", color.Black),
	   		widget.NewButton("Save", func() {
	   			if s, err := strconv.ParseFloat(inputCPUClock.Text, 64); err == nil {
	   				if s <= 1.0 {
	   					cpu.Clock = 1 // ticker requires positive value >= 1
	   				} else {
	   					cpu.Clock = s
	   				}
	   			}
	   			SetStatus(fmt.Sprintf("Clock set to %f milliseconds", cpu.Clock))
	   		}),
	   		canvas.NewText("Set clock speed in milliseconds. 1.0 sets clock to full speed.  ", color.Black),
	   		layout.NewSpacer(),
	   	)
	*/

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
		submitButton,
		resetButton,
		loadButton,
		runButton,
		stepButton,
		pauseButton,
		exitButton,
		layout.NewSpacer(),
		currentTime,
	)

	settingsContainer = container.NewVBox(
		buttonsContainer,
		//speedContainer,
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
	consoleBuffer.Write([]byte("******************************************************************************************\n" +
		"************************************* 6502 Simulator *************************************\n" +
		"******************************************************************************************\n"))
	consoleContainer.Refresh()
	UpdateAll()

	return w, &consoleBuffer
}

func UpdateAll() {

	// Reload
	//inputCPUClock.SetText(fmt.Sprintf("%3f", c.Clock))
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	consoleDispString = consoleBuffer.String() // Get whatever is in the memory buffer from the host
	consoleGridLabel.SetText(consoleDispString)
	registerDisplay = c.GetRegisters()
	registerDisplayWidget.Text = registerDisplay

	// Refresh
	buttonsContainer.Refresh()
	//speedContainer.Refresh()
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

func UpdateTime() {
	formatted := time.Now().Format("Time: 03:04:05")
	currentTime.SetText(formatted)
}

func SetStatus(s string) {
	status = s
	ConsoleWrite(status)
}

func ConsoleWrite(text string) {
	Console.Add(&canvas.Text{
		Text:      text,
		Color:     color.Black,
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
