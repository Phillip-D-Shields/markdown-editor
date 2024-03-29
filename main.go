package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {
	// create a fyne app
	app := app.New()

	// apply theme
	app.Settings().SetTheme(&myTheme{})

	// create a window for the app
	window := app.NewWindow("Markdown Editor")

	// get UI
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(window)

	// set content of the window
	window.SetContent(container.NewHSplit(edit, preview))

	// show the window and run the app
	window.Resize(fyne.Size{Width: 800, Height: 500})
	window.CenterOnScreen()
	window.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app *config) createMenuItems(window fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open", app.openFunc(window))

	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(window))

	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save As", app.saveAsFunc(window))

	// ? order of the menu items is important
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	window.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFunc(window fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *config) openFunc(window fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if read == nil {
				// user cancelled
				return
			}

			defer read.Close()

			// read the file
			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = read.URI()

			// update the window title
			window.SetTitle(window.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false

		}, window)
		// set the filter to only show markdown files
		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) saveAsFunc(window fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if write == nil {
				// user cancelled
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("error", "please use the .md extension", window)
				return
			}

			// save the file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			// !!! prevents resource leak
			defer write.Close()

			// update the window title
			window.SetTitle(window.Title() + " - " + write.URI().Name())

			// enable the save menu item
			app.SaveMenuItem.Disabled = false
		}, window)
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
