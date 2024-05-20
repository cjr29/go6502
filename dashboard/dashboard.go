package dashboard

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/cjr29/go6502/cpu"
)

var (
	c                     *cpu.CPU
	CPUStatus             string
	sps, pcs, flag        *widget.Label
	w                     fyne.Window
	status                string = "CPU status is displayed here."
	stackDisplay          string
	flagDisplay           string
	stackLabelWidget      *widget.Label
	stackHeader           *widget.Label
	registerHeader        *widget.Label
	registerDisplay       string
	registerDisplayWidget *widget.Label
	memoryDisplay         string
	memoryGridLabel       *widget.Label
	memoryLabel           *widget.Label
	inputCPUClock         *widget.Entry
	loadButton            *widget.Button
	runButton             *widget.Button
	stepButton            *widget.Button
	resetButton           *widget.Button
	pauseButton           *widget.Button
	exitButton            *widget.Button
	currentTime           *widget.Label
	mainContainer         *fyne.Container
	buttonsContainer      *fyne.Container
	settingsContainer     *fyne.Container
	statusContainer       *fyne.Container
	registerContainer     *fyne.Container
	memoryContainer       *fyne.Container
	stackContainer        *fyne.Container
	cpuInternalsContainer *fyne.Container
	speedContainer        *fyne.Container
	centerContainer       *fyne.Container
	middleContainer       *fyne.Container
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpu.CPU) fyne.Window {

	fmt.Println("Entering dashboard.New")
	c = cpu // All data comes from the CPU structure object
	a := app.NewWithID("6502")
	w = a.NewWindow("6502 Simulator")

	// Color backgrounds to be used in container stacks
	registerBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	stackBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	memoryBackground := canvas.NewRectangle(color.RGBA{R: 223, G: 159, B: 173, A: 200})

	// Control buttons
	loadButton = widget.NewButton("Load", load)
	runButton = widget.NewButton("Run", run)
	stepButton = widget.NewButton("Step", step)
	resetButton = widget.NewButton("Reset", reset)
	pauseButton = widget.NewButton("Pause", pause)
	exitButton = widget.NewButton("Exit", exit)

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
	// CPU Internals: PC, SP
	pcs = widget.NewLabel(fmt.Sprintf("PC: x%04x", cpu.Reg.PC))
	pcs.TextStyle.Monospace = true
	sps = widget.NewLabel(fmt.Sprintf("SP: x%04x", cpu.Reg.SP))
	sps.TextStyle.Monospace = true
	flag = widget.NewLabel("Flag: false")
	flag.TextStyle.Monospace = true
	cpuInternalsContainer = container.NewHBox(
		pcs,
		sps,
		flag,
	)

	// Stack
	stackHeader = widget.NewLabel("Top of Stack\n16-bit words\n(grows hi to lo)\n")
	stackHeader.TextStyle.Monospace = true
	stackHeader.TextStyle.Bold = true
	stackDisplay = cpu.GetStack()
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
	registerHeader = widget.NewLabel("Registers\n16-bit words\n")
	registerHeader.TextStyle.Monospace = true
	registerHeader.TextStyle.Bold = true
	registerDisplay = cpu.Reg.GetRegisters()
	registerDisplayWidget = widget.NewLabel(registerDisplay)
	registerDisplayWidget.TextStyle.Monospace = true
	registerDisplayWidget.TextStyle.Bold = true
	registerContainer = container.NewStack(
		registerBackground,
		container.NewVBox(
			registerHeader,
			registerDisplayWidget,
		))

	// Memory
	memoryDisplay = cpu.Mem.GetAllMemory()
	memoryLabel = widget.NewLabel("Memory\nbytes\n")
	memoryLabel.TextStyle.Monospace = true
	memoryLabel.TextStyle.Bold = true
	memoryGridLabel = widget.NewLabel(memoryDisplay)
	memoryGridLabel.TextStyle.Monospace = true
	memoryContainer = container.NewStack(
		memoryBackground,
		container.NewVBox(
			memoryLabel,
			memoryGridLabel,
		))

	buttonsContainer = container.NewHBox(
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
		speedContainer,
		cpuInternalsContainer,
	)

	middleContainer = container.NewHBox(
		registerContainer,
		memoryContainer,
		stackContainer,
	)

	statusContainer = container.NewVBox(ConsoleScroller)
	registerContainer = container.NewHBox(registerContainer)
	centerContainer = container.NewHBox(memoryContainer, stackContainer)

	mainContainer = container.NewVBox(
		settingsContainer,
		middleContainer,
		statusContainer,
	)

	w.SetContent(mainContainer)

	return w
}

func UpdateAll() {

	// Reload
	pcs.SetText(fmt.Sprintf("PC: x%04x", c.Reg.PC))
	sps.SetText(fmt.Sprintf("SP: x%04x", c.Reg.SP))
	if c.Reg.Carry {
		flagDisplay = "Flag: true"
	} else {
		flagDisplay = "Flag: false"
	}
	flag.SetText(flagDisplay)
	//	inputCPUClock.SetText(fmt.Sprintf("%3f", c.Clock))
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	memoryDisplay = GetAllMemory()
	memoryGridLabel.SetText(memoryDisplay)
	registerDisplay = c.GetRegisters()
	registerDisplayWidget.Text = registerDisplay

	// Refresh
	buttonsContainer.Refresh()
	speedContainer.Refresh()
	cpuInternalsContainer.Refresh()
	settingsContainer.Refresh()
	stackLabelWidget.Refresh()
	stackContainer.Refresh()
	memoryGridLabel.Refresh()
	memoryContainer.Refresh()
	registerContainer.Refresh()
	middleContainer.Refresh()
	statusContainer.Refresh()
	centerContainer.Refresh()
	mainContainer.Refresh()
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

//=================== Added Chris Riddick 2024 ======================

func load() {
	// Loads code in []program into CPU memory at index 0
	// DO NOTHING
}

func run() {
	// DO NOTHING
}

func step() {
	// DO NOTHING
}

func reset() {
	// DO NOTHING
}

func pause() {
	// DO NOTHING
}

func exit() {
	os.Exit(0)
}

//=================== Added, Chris Riddick, 2024 ====================

// GetRegisters returns a formatted string of register values
func GetRegisters() string {
	//var s string
	// for i := 0; i < len(c.Reg.A); i++ {
	// 	s = s + fmt.Sprintf("R%02d: x%04x\n", i, c.Registers[i])
	// }
	s := "Register\nlisting\ngoes\nhere\n"
	return s
}

// GetStack returns a formatted string of bytes beginning at SP down to to of stack
func GetStack() string {
	//var s string
	// for i := c.Reg.SP; i <= c.StackHead; i = i + 2 {
	// 	s = s + fmt.Sprintf("%04x: x%04x\n", i, binary.BigEndian.Uint16(c.Memory[i:]))
	// }
	s := "Stack\nlisting\ngoes\nhere\n"
	return s
}

// GetAllMemory returns a 16 byte formatted string starting at 0000
func GetAllMemory() string {
	//var line string
	/*	blocks := len(c.Memory) / 16
		remainder := len(c.Memory) % 16
		// Send header line with memory locations
		line = "       00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f\n"
		k := 0
		for j := 0; j < blocks; j++ {
			line = line + fmt.Sprintf("%04x:  ", k)
			for i := k; i < k+16; i++ {
				line = line + fmt.Sprintf("%02x ", c.Memory[i])
			}
			line = line + "\n"
			k = k + 16
		}
		if k >= len(c.Memory) {
			return line
		}
		endBlock := blocks * 16
		line = line + fmt.Sprintf("%04x:  ", k)
		for i := endBlock; i < endBlock+remainder; i++ {
			line = line + fmt.Sprintf("%02x ", c.Memory[i])
		}
		line = line + "\n" */
	line := "Memoryk\nlisting\ngoes\nhere\n"
	return line
}
