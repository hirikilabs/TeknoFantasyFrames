package main

import (
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func extractFrames(video string, folder string, step float64) error {
	err := ffmpeg.Input(video, ffmpeg.KwArgs{}).
		Output(folder+string(os.PathSeparator)+
			"frame_%03d.jpg", ffmpeg.KwArgs{"vf": "select=not(mod(n\\," +
			strconv.FormatFloat(step, 'f', 0, 64) +
			"))", "vsync": "vfr"}).
		OverWriteOutput().ErrorToStdOut().Run()

	return err
}

func main() {
	videoFile := ""
	outputFolder := ""
	inputSelected := false
	outputSelected := false

	frameStep := 5.0

	// create GUI
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	// window and size
	w := a.NewWindow("Frames - TeknoFantasy")
	w.Resize(fyne.NewSize(300, 400))

	// image background
	bgImg := canvas.NewImageFromFile("logo.jpg")
	//bgImg.FillMode = canvas.ImageFillOriginal

	// widgets
	// input
	fileLabel := widget.NewLabel("")
	openLabel := widget.NewLabel("Open Video")
	openButton := widget.NewButton("...", func() {
		dialog.ShowFileOpen(func(u fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if u != nil {
				videoFile = u.URI().String()
				fileLabel.SetText(videoFile)
				inputSelected = true
			}
		}, w)
	})
	openContainer := container.New(layout.NewHBoxLayout(), openLabel, openButton)

	// output
	folderLabel := widget.NewLabel("")
	outputLabel := widget.NewLabel("Output Folder")
	outputButton := widget.NewButton("...", func() {
		// open file dialog and save output
		dialog.ShowFolderOpen(func(u fyne.ListableURI, err error) {
			if err != nil {
				return
			}
			if u != nil {
				outputFolder = u.String()
				folderLabel.SetText(outputFolder)
				outputSelected = true
			}
		}, w)

	})
	outputContainer := container.New(layout.NewHBoxLayout(), outputLabel, outputButton)

	// slider
	sliderLabel := widget.NewLabel("Frames step: " + strconv.FormatFloat(frameStep, 'f', 0, 64))
	slider := widget.NewSlider(5.0, 30.0)
	slider.OnChanged = func(f float64) {
		sliderLabel.SetText("Frames step: " + strconv.FormatFloat(f, 'f', 0, 64))
		frameStep = f
	}
	sliderContainer := container.New(layout.NewVBoxLayout(), sliderLabel, slider)

	// run
	runButton := widget.NewButton("RUN", nil)
	runButtonFunc := func() {
		if inputSelected && outputSelected {
			runButton.SetText("Running...")
			runButton.Disable()

			// launch ffmpeg command
			err := extractFrames(videoFile, outputFolder, frameStep)

			if err != nil {
				dialog.ShowError(err, w)
			}

			runButton.SetText("RUN")
			runButton.Enable()
		}

	}
	runButton.OnTapped = runButtonFunc

	content := container.New(layout.NewVBoxLayout(), openContainer, fileLabel,
		widget.NewSeparator(), outputContainer, folderLabel,
		widget.NewSeparator(), sliderContainer,
		widget.NewSeparator(), layout.NewSpacer(), runButton)

	w.SetContent(container.NewStack(bgImg, content))
	w.ShowAndRun()
}
